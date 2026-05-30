#!/bin/sh

# Ensure the config file exists so the engine doesn't panic
if [ ! -f "/app/config.json" ]; then
    echo "⚠️  No config.json found! Creating a placeholder..."
    echo "{}" > /app/config.json
fi

# Check if the container was started as root
if [ "$(id -u)" = "0" ]; then
    # Running in standard Docker (Root): Fix permissions and drop to reduser
    chown -R reduser:redgroup /app/data
    chown reduser:redgroup /app/config.json
    exec su-exec reduser "$@"
else
    # Running in Podman keep-id (Non-Root): Permissions are already mapped natively
    exec "$@"
fi