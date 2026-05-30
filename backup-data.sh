#!/bin/bash

cd "$(dirname "$0")"

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y-%m-%d_%H%M%S)
mkdir -p "$BACKUP_DIR"

echo "[*] Starting backup..."
tar -czvf "$BACKUP_DIR/data-backup-$TIMESTAMP.tar.gz" ./data
find "$BACKUP_DIR" -type f -name "*.tar.gz" -mtime +7 -delete
echo "[*] Backup complete: $BACKUP_DIR/data-backup-$TIMESTAMP.tar.gz"