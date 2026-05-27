$ConfigFile = "config.json"

if (-Not (Test-Path $ConfigFile)) {
    Write-Host "Error: $ConfigFile not found in the current directory!" -ForegroundColor Red
    Exit
}

# Parse the JSON file
$Config = Get-Content $ConfigFile -Raw | ConvertFrom-Json

$CurrentToken = $Config.adminToken

Write-Host "----------------------------------------"
if ([string]::IsNullOrEmpty($CurrentToken)) {
    Write-Host "Current Admin Token: [NONE / NOT SET]" -ForegroundColor Yellow
} else {
    Write-Host "Current Admin Token: $CurrentToken" -ForegroundColor Cyan
}
Write-Host "----------------------------------------`n"

$Choice = Read-Host "Would you like to generate and save a new secure token? (y/N)"

if ($Choice -match "^[yY]") {
    # Generate a cryptographically secure 32-character hexadecimal token
    $Bytes = New-Object Byte[] 16
    [Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($Bytes)
    $NewToken = [BitConverter]::ToString($Bytes) -replace '-'

    # Update the object and save it back to disk
    $Config.adminToken = $NewToken
    $Config | ConvertTo-Json -Depth 10 | Set-Content $ConfigFile$ConfigFile = "config.json"

    if (-Not (Test-Path $ConfigFile)) {
        Write-Host "Error: $ConfigFile not found in the current directory!" -ForegroundColor Red
        Exit
    }

    # Parse the JSON file
    $Config = Get-Content $ConfigFile -Raw | ConvertFrom-Json

    $CurrentToken = $Config.adminToken

    Write-Host "----------------------------------------"
    if ([string]::IsNullOrEmpty($CurrentToken)) {
        Write-Host "Current Admin Token: [NONE / NOT SET]" -ForegroundColor Yellow
    } else {
        Write-Host "Current Admin Token: $CurrentToken" -ForegroundColor Cyan
    }
    Write-Host "----------------------------------------`n"

    $Choice = Read-Host "Would you like to generate and save a new secure token? (y/N)"

    if ($Choice -match "^[yY]") {
        # Generate a cryptographically secure 32-character hexadecimal token
        $Bytes = New-Object Byte[] 16
        [Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($Bytes)
        $NewToken = [BitConverter]::ToString($Bytes) -replace '-'

        # Update the object and save it back to disk
        $Config.adminToken = $NewToken
        $Config | ConvertTo-Json -Depth 10 | Set-Content $ConfigFile

        Write-Host "✅ Token updated successfully!" -ForegroundColor Green
        Write-Host "Your new token is: $NewToken" -ForegroundColor Cyan
        Write-Host "⚠️  Make sure to restart your node: podman-compose restart red_engine" -ForegroundColor Yellow
    } else {
        Write-Host "Operation cancelled. Token unchanged."
    }

    Write-Host "✅ Token updated successfully!" -ForegroundColor Green
    Write-Host "Your new token is: $NewToken" -ForegroundColor Cyan
    Write-Host "⚠️  Make sure to restart your node: podman-compose restart red_engine" -ForegroundColor Yellow
} else {
    Write-Host "Operation cancelled. Token unchanged."
}
