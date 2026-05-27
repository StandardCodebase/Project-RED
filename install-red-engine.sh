#!/bin/bash
echo "========================================"
echo "🚀 Installing RED Engine..."
echo "========================================"

# 1. Create data directory safely as the standard user (NO SUDO)
if [ ! -d "./data" ]; then
    echo "[*] Creating ./data directory..."
    mkdir -p ./data
else
    echo "[*] ./data directory already exists."
fi

# 2. Check for or create config.json with a secure token
if [ ! -f "config.json" ]; then
    echo "[*] Generating default config.json..."
    # Generate a secure 24-character token
    NEW_TOKEN=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1)

    cat <<EOF > config.json
{
  "addr": ":8080",
  "siteName": "RED Engine",
  "dataDir": "/app/data",
  "adminToken": "$NEW_TOKEN",
  "startupSync": []
}
EOF
    echo "[*] Generated secure Admin Token: $NEW_TOKEN"
    echo "⚠️  PLEASE SAVE THIS TOKEN! You will need it to log in to the admin panel."
else
    echo "[*] config.json already exists. Skipping default generation."
fi

# 3. Detect the container engine
if command -v podman-compose &> /dev/null; then
    COMPOSE_CMD="podman-compose"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
else
    echo "❌ Error: Neither podman-compose nor docker-compose found on this system."
    echo "Please install Podman or Docker to continue."
    exit 1
fi

echo "[*] Starting RED Engine using $COMPOSE_CMD..."
$COMPOSE_CMD up --build -d

echo "========================================"
echo "✅ Installation Complete!"
echo "🌐 Your node is running at: http://localhost"
echo "⚙️  Admin Panel: http://localhost/-/admin"
echo "========================================"
