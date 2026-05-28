#!/bin/bash
echo "========================================"
echo "🚀 Installing RED Engine (Production Mode)..."
echo "========================================"

# Check for root/sudo privileges to bind ports 80/443
if [[ $EUID -ne 0 ]]; then
   echo "⚠️  Sudo privileges are required to bind ports 80 and 443."
   echo "Please enter your password when prompted."
   sudo "$0" "$@"
   exit $?
fi

# 1. Repository Check
if [ ! -f "docker-compose.yml" ]; then
    echo "[*] Cloning RED Engine repository..."
    git clone https://github.com/RED-Collective/red-engine.git
    cd red-engine || exit 1
fi

# 2. Setup Directories
mkdir -p ./data

# 3. Handle config.json
if [ ! -f "config.json" ]; then
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

# 4. Handle contributors.json
if [ ! -f "contributors.json" ]; then
    echo "[]" > contributors.json
fi

# 5. Build and Deploy
echo "[*] Building local image..."
podman build --network=host -t red-engine-image .

echo "[*] Starting services..."
podman-compose up -d

# 6. Final Status
TOKEN=$(grep -oP '"adminToken": "\K[^"]+' config.json)
echo "========================================"
echo "✅ Installation Complete!"
echo "🌐 Node running at: http://localhost"
echo "⚙️  Admin Panel: http://localhost/-/admin"
echo "🔑 YOUR ADMIN TOKEN: $TOKEN"
echo "⚠️  PLEASE SAVE THIS TOKEN!"
echo "========================================"
