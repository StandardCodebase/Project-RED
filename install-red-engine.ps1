Write-Host "========================================" -ForegroundColor Cyan
Write-Host "🚀 Installing RED Engine..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. Check if we are inside the repository; if not, clone it.
if (-Not (Test-Path "docker-compose.yml"))
{
    Write-Host "[*] Repository not detected in current directory."

    if (-Not (Get-Command "git" -ErrorAction SilentlyContinue))
    {
        Write-Host "❌ Error: 'git' is not installed. Please install Git for Windows to continue." -ForegroundColor Red
        Exit
    }

    Write-Host "[*] Cloning RED Engine repository..."
    git clone https://github.com/RED-Collective/red-engine.git

    if ($LASTEXITCODE -ne 0)
    {
        Write-Host "❌ Error: Failed to clone repository." -ForegroundColor Red
        Exit
    }

    Write-Host "[*] Navigating into red-engine directory..."
    Set-Location "red-engine"
} else
{
    Write-Host "[*] Running from inside existing repository."
}

# 2. Create data directory safely as the standard user
if (-Not (Test-Path ".\data"))
{
    Write-Host "[*] Creating .\data directory..."
    New-Item -ItemType Directory -Path ".\data" | Out-Null
} else
{
    Write-Host "[*] .\data directory already exists."
}

# 3. Check for or create config.json with a secure token
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

# 4. Detect the container engine
$ComposeCmd = ""
$ComposeArgs = @("up", "--build", "-d")

if (Get-Command "podman-compose" -ErrorAction SilentlyContinue)
{
    $ComposeCmd = "podman-compose"
} elseif (Get-Command "docker-compose" -ErrorAction SilentlyContinue)
{
    $ComposeCmd = "docker-compose"
} elseif (Get-Command "docker" -ErrorAction SilentlyContinue)
{
    # Check if modern 'docker compose' (V2) is available
    try
    {
        $null = Invoke-Expression "docker compose version 2>&1"
        if ($LASTEXITCODE -eq 0)
        {
            $ComposeCmd = "docker"
            $ComposeArgs = @("compose", "up", "--build", "-d")
        }
    } catch
    {
    }
}

if ($ComposeCmd -eq "")
{
    Write-Host "❌ Error: Neither podman-compose nor docker compose found on this system." -ForegroundColor Red
    Write-Host "Please install Podman or Docker Desktop to continue." -ForegroundColor Red
    Exit
}

Write-Host "[*] Starting RED Engine..."
& $ComposeCmd $ComposeArgs

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ Installation Complete!" -ForegroundColor Green
Write-Host "🌐 Your node is running at: http://localhost"
Write-Host "⚙️  Admin Panel: http://localhost/-/admin"
Write-Host "========================================" -ForegroundColor Cyan
