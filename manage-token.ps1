$ConfigFile = "config.json"

# Helper function to generate a secure token
function New-SecureToken
{
    $Bytes = New-Object Byte[] 16
    [Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($Bytes)
    return [BitConverter]::ToString($Bytes) -replace '-'
}

# 1. Check/Create Config
if (-Not (Test-Path $ConfigFile))
{
    Write-Host "⚠️  $ConfigFile not found. Creating a default configuration..." -ForegroundColor Yellow

    $NewToken = New-SecureToken
    $DefaultConfig = [PSCustomObject]@{
        addr       = ":8080"
        siteName   = "RED Engine"
        dataDir    = "/app/data"
        adminToken = $NewToken
        startupSync = @(
            @{
                url      = "https://github.com/mundimark/awesome-markdown"
                filename = "awesome-markdown.md"
            }
        )
    }
    $DefaultConfig | ConvertTo-Json -Depth 10 | Set-Content $ConfigFile
    Write-Host "✅ Created new config.json with secure token." -ForegroundColor Green
}

# 2. Parse and Display
$Config = Get-Content $ConfigFile -Raw | ConvertFrom-Json
Write-Host "`n----------------------------------------"
Write-Host "Current Admin Token: $($Config.adminToken)" -ForegroundColor Cyan
Write-Host "----------------------------------------`n"

# 3. Interactive Update
$Choice = Read-Host "Would you like to generate and save a new secure token? (y/N)"

if ($Choice -match "^[yY]")
{
    $Config.adminToken = New-SecureToken
    $Config | ConvertTo-Json -Depth 10 | Set-Content $ConfigFile

    # Update the object and save it back to disk
    $Config.adminToken = $NewToken
    $Config | ConvertTo-Json -Depth 10 | Set-Content $ConfigFile

    Write-Host "✅ Token updated successfully!" -ForegroundColor Green
    Write-Host "Your new token is: $NewToken" -ForegroundColor Cyan
    Write-Host "⚠️  Make sure to restart your node: podman-compose restart red_engine" -ForegroundColor Yellow
} else
{
    Write-Host "Operation cancelled. Token unchanged."
}