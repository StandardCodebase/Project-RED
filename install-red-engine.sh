#!/bin/bash
echo "========================================"
echo "🚀 Installing RED Engine..."
echo "========================================"

if [ ! -f "docker-compose.yml" ]; then
    echo "[*] Repository not detected in current directory."

    if ! command -v git &> /dev/null; then
        echo "❌ Error: 'git' is not installed. Please install git to continue."
        exit 1
    fi

    echo "[*] Cloning RED Engine repository..."
    git clone https://github.com/RED-Collective/red-engine.git
    if [ $? -ne 0 ]; then
        echo "❌ Error: Failed to clone repository."
        exit 1
    fi

    echo "[*] Navigating into red-engine directory..."
    cd red-engine || exit 1
else
    echo "[*] Running from inside existing repository."
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
    echo "[*] Generated secure Admin Token: $NEW_TOKEN"
    echo "⚠️  PLEASE SAVE THIS TOKEN! You will need it to log in to the admin panel."
else
    echo "[*] config.json already exists. Skipping default generation."
fi

if [ ! -f "contributors.json" ]; then
    echo "[*] Generating default contributors.json..."
    echo "[]" > contributors.json
else
    echo "[*] contributors.json already exists."
fi

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

if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to start containers."
    exit 1
fi

CONFIG_PORT=$(grep '"addr"' config.json | sed -E 's/.*:([0-9]+).*/\1/')

if [ -z "$CONFIG_PORT" ]; then
    CONFIG_PORT="8080"
fi

HOST_IP="localhost"
if command -v hostname &> /dev/null; then
    HOST_IP=$(hostname -I | awk '{print $1}')
    [ -z "$HOST_IP" ] && HOST_IP="localhost"
fi

echo "========================================"
echo "✅ Installation Complete!"
echo "🌐 Your node is running at: http://${HOST_IP}:${CONFIG_PORT}"
echo "⚙️  Admin Panel: http://${HOST_IP}:${CONFIG_PORT}/-/admin"
echo "========================================"
