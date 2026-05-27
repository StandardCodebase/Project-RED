#!/bin/bash
CONFIG_FILE="config.json"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: $CONFIG_FILE not found in the current directory!"
    exit 1
fi

# Extract the current token safely
CURRENT_TOKEN=$(grep '"adminToken"' "$CONFIG_FILE" | sed -E 's/.*"adminToken"[[:space:]]*:[[:space:]]*"([^"]*)".*/\1/')

echo "----------------------------------------"
if [ -z "$CURRENT_TOKEN" ]; then
    echo "Current Admin Token: [NONE / NOT SET]"
else
    echo "Current Admin Token: $CURRENT_TOKEN"
fi
echo "----------------------------------------"
echo ""

read -p "Would you like to generate and save a new secure token? (y/N): " choice
case "$choice" in
  y|Y )
    # Generate a secure 24-character random string
    NEW_TOKEN=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1)

    # Safely replace the token in the JSON file
    if grep -q '"adminToken"' "$CONFIG_FILE"; then
        sed -i.bak -E 's/("adminToken"[[:space:]]*:[[:space:]]*")[^"]*(")/\1'"$NEW_TOKEN"'\2/' "$CONFIG_FILE"
        rm -f "$CONFIG_FILE.bak"
    else
        echo "Error: 'adminToken' key not found in $CONFIG_FILE. Please add it manually."
        exit 1
    fi

    echo "✅ Token updated successfully!"
    echo "Your new token is: $NEW_TOKEN"
    echo "⚠️  Make sure to restart your node: podman-compose restart red_engine"
    ;;
  * )
    echo "Operation cancelled. Token unchanged."
    ;;
esac
