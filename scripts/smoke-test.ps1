param(
    [string]$BaseUrl = "http://localhost:8080"
)

$ErrorActionPreference = "Stop"

$BaseUrl = $BaseUrl.TrimEnd("/")
$Title = "CRM AI Assistant smoke test $(Get-Date -Format 'yyyyMMdd-HHmmss')"
$Content = @"
Lead response policy
Sales representatives must make the first contact attempt within one business day after a new lead is assigned. Record the attempt as an activity in the CRM and create a follow-up task if the customer does not respond.

Opportunity updates
Every open opportunity must have an owner, estimated amount, next step, and expected close date based on the customer's stated timeline.
"@

Write-Host "==> Checking API health at $BaseUrl..."
$health = Invoke-RestMethod -Method Get -Uri "$BaseUrl/health"
if ($health -ne "ok") {
    throw "Unexpected health response: $health"
}

Write-Host "==> Creating a sample document and its embeddings..."
$documentBody = @{
    title = $Title
    content = $Content
} | ConvertTo-Json
$created = Invoke-RestMethod `
    -Method Post `
    -Uri "$BaseUrl/documents" `
    -ContentType "application/json" `
    -Body $documentBody

if (-not $created.id) {
    throw "Document creation response did not include an ID"
}

Write-Host "Created document ID $($created.id). Asking a question..."
$askBody = @{
    question = "How soon should sales make the first contact attempt after a new lead is assigned?"
} | ConvertTo-Json
$answer = Invoke-RestMethod `
    -Method Post `
    -Uri "$BaseUrl/ask" `
    -ContentType "application/json" `
    -Body $askBody

if ([string]::IsNullOrWhiteSpace($answer.answer)) {
    throw "The answer endpoint returned an empty answer"
}

$matchingSource = @($answer.sources | Where-Object { $_.document_id -eq $created.id })
if ($matchingSource.Count -eq 0) {
    throw "The answer did not cite the document created by this smoke check"
}

Write-Host "==> Smoke check passed."
Write-Host "Answer: $($answer.answer)"
Write-Host "Source: $($matchingSource[0].title), chunk $($matchingSource[0].chunk_index)"
