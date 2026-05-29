#!/bin/bash
echo "========================================"
echo "🚀 Installing RED Engine (Production Mode)..."
echo "========================================"

if [ ! -f "docker-compose.yml" ]; then
    echo "[*] Cloning RED Engine repository..."
    git clone https://github.com/RED-Collective/red-engine.git
    cd red-engine || exit 1
fi

if [ ! -d "./data" ]; then
    echo "[*] Creating ./data directory..."
    mkdir -p ./data
else
    echo "[*] ./data directory already exists."
fi

if [ ! -f "config.json" ]; then
    echo "[*] Generating default config.json..."
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
    echo "[*] Generated Admin Token: $NEW_TOKEN"
fi

if [ ! -f "contributors.json" ]; then
    echo "[]" > contributors.json
else
    echo "[*] contributors.json already exists."
fi

# 5. Detect the container engine
if command -v podman-compose &> /dev/null; then
    COMPOSE_CMD="podman-compose up --build -d"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose up --build -d"
elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose up --build -d"
else
    echo "❌ Error: Neither podman-compose nor docker compose found on this system."
    echo "Please install Podman or Docker to continue."
    exit 1
fi

echo "[*] Starting RED Engine using container engine..."
$COMPOSE_CMD

echo "========================================"
echo "✅ Installation Complete!"
echo "🌐 Your node is running at: http://${HOST_IP}:${CONFIG_PORT}"
echo "⚙️  Admin Panel: http://${HOST_IP}:${CONFIG_PORT}/-/admin"
echo "========================================"
