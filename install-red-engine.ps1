Write-Host "========================================" -ForegroundColor Cyan
Write-Host "🚀 Installing RED Engine (Production Mode)..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Check for Administrator privileges
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
{
    Write-Host "⚠️  Administrator privileges are required to bind network ports." -ForegroundColor Yellow
    Write-Host "Re-launching as Administrator..." -ForegroundColor Yellow
    Start-Process powershell.exe -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`"" -Verb RunAs
    Exit
}

# 1. Repository Check
if (-Not (Test-Path "docker-compose.yml"))
{
    git clone https://github.com/RED-Collective/red-engine.git
    Set-Location "red-engine"
}

# 2. Setup Directories
if (-Not (Test-Path ".\data"))
{ New-Item -ItemType Directory -Path ".\data" | Out-Null 
}

# 3. Handle config.json
if (-Not (Test-Path "config.json"))
{
    $Bytes = New-Object Byte[] 16
    [Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($Bytes)
    $NewToken = [BitConverter]::ToString($Bytes) -replace '-'
    $DefaultConfig = @{ addr = ":8080"; siteName = "RED Engine"; dataDir = "/app/data"; adminToken = $NewToken; startupSync = @() }
    $DefaultConfig | ConvertTo-Json -Depth 10 | Set-Content "config.json"
}

# 4. Handle contributors.json
if (-Not (Test-Path "contributors.json"))
{ "[]" | Set-Content "contributors.json" 
}

# 5. Build and Deploy
Write-Host "[*] Building local image..." -ForegroundColor Green
podman build --network=host -t red-engine-image .

Write-Host "[*] Starting services..." -ForegroundColor Green
podman-compose up -d

# 6. Final Status
$Config = Get-Content "config.json" | ConvertFrom-Json
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ Installation Complete!" -ForegroundColor Green
Write-Host "🌐 Node running at: http://localhost"
Write-Host "⚙️  Admin Panel: http://localhost/-/admin"
Write-Host "🔑 YOUR ADMIN TOKEN: $($Config.adminToken)" -ForegroundColor Yellow
Write-Host "⚠️  PLEASE SAVE THIS TOKEN!" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan
