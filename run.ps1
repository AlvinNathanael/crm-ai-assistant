$ErrorActionPreference = "Stop"

$ComposeFile = "docker/docker-compose.yml"
$PostgresContainer = "crm-ai-postgres"

if (-not (Test-Path ".env")) {
    Write-Error "Error: .env file not found"
    exit 1
}

Write-Host "==> Loading environment variables..."

Get-Content ".env" | ForEach-Object {
    $line = $_.Trim()

    if ($line -eq "" -or $line.StartsWith("#")) {
        return
    }

    $parts = $line -split "=", 2

    if ($parts.Count -ne 2) {
        return
    }

    $name = $parts[0].Trim()
    $value = $parts[1].Trim()

    Set-Item -Path "Env:$name" -Value $value
}

Write-Host "==> Validating environment variables..."

$RequiredEnvVars = @(
    "POSTGRES_USER",
    "POSTGRES_PASSWORD",
    "POSTGRES_DB",
    "POSTGRES_DSN",
    "MIGRATION_DATABASE_URL",
    "OPENAI_API_KEY"
)

foreach ($name in $RequiredEnvVars) {
    $value = [Environment]::GetEnvironmentVariable($name)

    if ([string]::IsNullOrWhiteSpace($value)) {
        throw "Required environment variable '$name' is missing or empty"
    }
}

Write-Host "==> Checking Docker Desktop..."

function Test-DockerEngine {
    $PreviousErrorActionPreference = $ErrorActionPreference
    try {
        # Windows PowerShell can turn native stderr into a terminating error
        # when ErrorActionPreference is Stop, even when checking a nonzero exit.
        $ErrorActionPreference = "Continue"
        docker info 2>$null | Out-Null
        return $LASTEXITCODE -eq 0
    }
    finally {
        $ErrorActionPreference = $PreviousErrorActionPreference
    }
}

if (-not (Test-DockerEngine)) {
    Write-Host "Docker engine is not reachable. Starting Docker Desktop..."

    try {
        $DockerCLI = (Get-Command docker -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source
        $DockerStartProcess = Start-Process `
            -FilePath $DockerCLI `
            -ArgumentList @("desktop", "start") `
            -WindowStyle Hidden `
            -PassThru
    }
    catch {
        throw "Could not start Docker Desktop. Start it manually, wait for the engine, then run this script again. Details: $($_.Exception.Message)"
    }

    $DockerMaxAttempts = 60
    $DockerAttempt = 0

    while ($DockerAttempt -lt $DockerMaxAttempts) {
        if (Test-DockerEngine) {
            break
        }

        if ($DockerStartProcess.HasExited -and $DockerStartProcess.ExitCode -ne 0) {
            throw "Docker Desktop start command exited with code $($DockerStartProcess.ExitCode). Start Docker Desktop manually and check its Linux engine."
        }

        $DockerAttempt++
        Write-Host "Waiting for Docker Desktop engine... ($DockerAttempt/$DockerMaxAttempts)"
        Start-Sleep -Seconds 2
    }

    if ($DockerAttempt -ge $DockerMaxAttempts) {
        throw "Docker Desktop did not become ready within $($DockerMaxAttempts * 2) seconds. Check Docker Desktop and its Linux engine, then run this script again."
    }
}

Write-Host "==> Docker Desktop is ready."

Write-Host "==> Starting PostgreSQL..."

docker compose -f $ComposeFile up -d

if ($LASTEXITCODE -ne 0) {
    throw "Failed to start PostgreSQL"
}

Write-Host "==> Waiting for PostgreSQL..."

$MaxAttempts = 30
$Attempt = 0

while ($Attempt -lt $MaxAttempts) {
    docker exec $PostgresContainer `
        pg_isready `
        -U $env:POSTGRES_USER `
        -d $env:POSTGRES_DB

    if ($LASTEXITCODE -eq 0) {
        break
    }

    $Attempt++

    Write-Host "PostgreSQL is not ready yet... ($Attempt/$MaxAttempts)"
    Start-Sleep -Seconds 1
}

if ($Attempt -ge $MaxAttempts) {
    throw "PostgreSQL did not become ready within $MaxAttempts seconds"
}

Write-Host "==> PostgreSQL is ready."

Write-Host "==> Running database migrations..."

docker run --rm `
    -v "${PWD}/migrations:/migrations" `
    migrate/migrate `
    -path=/migrations `
    -database "$env:MIGRATION_DATABASE_URL" `
    up

if ($LASTEXITCODE -ne 0) {
    throw "Migration failed"
}

Write-Host "==> Starting CRM AI Assistant API..."

go run ./cmd/api

if ($LASTEXITCODE -ne 0) {
    throw "Failed to start CRM AI Assistant API"
}
