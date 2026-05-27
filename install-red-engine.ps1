Write-Host "========================================" -ForegroundColor Cyan
Write-Host "🚀 Installing RED Engine..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. Create data directory safely as the standard user
if (-Not (Test-Path ".\data"))
{
    Write-Host "[*] Creating .\data directory..."
    New-Item -ItemType Directory -Path ".\data" | Out-Null
} else
{
    Write-Host "[*] .\data directory already exists."
}

# 2. Check for or create config.json with a secure token
if (-Not (Test-Path "config.json"))
{
    Write-Host "[*] Generating default config.json..."

    # Generate a cryptographically secure 32-character hexadecimal token
    $Bytes = New-Object Byte[] 16
    [Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($Bytes)
    $NewToken = [BitConverter]::ToString($Bytes) -replace '-'

    $DefaultConfig = @{
        addr = ":8080"
        siteName = "RED Engine"
        dataDir = "/app/data"
        adminToken = $NewToken
        startupSync = @()
    }
    $DefaultConfig | ConvertTo-Json -Depth 10 | Set-Content "config.json"

    Write-Host "[*] Generated secure Admin Token: $NewToken" -ForegroundColor Green
    Write-Host "⚠️  PLEASE SAVE THIS TOKEN! You will need it to log in to the admin panel." -ForegroundColor Yellow
} else
{
    Write-Host "[*] config.json already exists. Skipping default generation."
}

# 3. Detect the container engine
$ComposeCmd = ""
if (Get-Command "podman-compose" -ErrorAction SilentlyContinue)
{
    $ComposeCmd = "podman-compose"
} elseif (Get-Command "docker-compose" -ErrorAction SilentlyContinue)
{
    $ComposeCmd = "docker-compose"
} else
{
    Write-Host "❌ Error: Neither podman-compose nor docker-compose found on this system." -ForegroundColor Red
    Write-Host "Please install Podman or Docker Desktop to continue." -ForegroundColor Red
    Exit
}

Write-Host "[*] Starting RED Engine using $ComposeCmd..."
if ($ComposeCmd -eq "podman-compose")
{
    podman-compose up --build -d
} else
{
    docker-compose up --build -d
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ Installation Complete!" -ForegroundColor Green
Write-Host "🌐 Your node is running at: http://localhost"
Write-Host "⚙️  Admin Panel: http://localhost/-/admin"
Write-Host "========================================" -ForegroundColor Cyan
