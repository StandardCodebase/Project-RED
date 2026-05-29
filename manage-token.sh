#!/bin/bash
CONFIG_FILE="config.json"

# Helper function for token
generate_token() {
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
}

# 1. Check/Create Config
if [ ! -f "$CONFIG_FILE" ]; then
    echo "⚠️  $CONFIG_FILE not found. Creating a default configuration..."
    NEW_TOKEN=$(generate_token)

    cat <<EOF > "$CONFIG_FILE"
{
  "addr": ":8080",
  "siteName": "RED Engine",
  "dataDir": "/app/data",
  "adminToken": "$NEW_TOKEN",
  "startupSync": [
    {
      "url": "https://github.com/mundimark/awesome-markdown",
      "filename": "awesome-markdown.md"
    }
  ]
}
EOF
    echo "✅ Created new config.json."
fi

# 2. Display Current Token
TOKEN=$(grep -oP '"adminToken": "\K[^"]+' "$CONFIG_FILE")
echo "----------------------------------------"
echo "Current Admin Token: $TOKEN"
echo "----------------------------------------"

# 3. Interactive Update
read -p "Would you like to generate and save a new secure token? (y/N): " choice
if [[ "$choice" =~ ^[yY]$ ]]; then
    NEW_TOKEN=$(generate_token)
    # Use sed to replace the token safely
    sed -i "s/\"adminToken\": \".*\"/\"adminToken\": \"$NEW_TOKEN\"/" "$CONFIG_FILE"

    echo "✅ Token updated successfully!"
    echo "Your new token is: $NEW_TOKEN"
    echo "⚠️  Restart your node: podman-compose restart red_engine"
else
    echo "Operation cancelled. Token unchanged."
fi
