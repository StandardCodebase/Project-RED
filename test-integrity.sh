#!/bin/bash
set -e

echo "=========================================="
echo "RED Engine Integrity Verification Test"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Setup
VAULT_PATH="data/remote/API"
CONFIG_FILE="config-integrity-test.json"
PORT=9999

echo -e "${YELLOW}[1] Checking test environment...${NC}"
if [ ! -f "$VAULT_PATH/manifest.json" ]; then
  echo -e "${RED}✗ manifest.json not found${NC}"
  exit 1
fi
echo -e "${GREEN}✓ manifest.json found${NC}"

if [ ! -f "contributors.json" ]; then
  echo -e "${RED}✗ contributors.json not found${NC}"
  exit 1
fi
echo -e "${GREEN}✓ contributors.json found${NC}"

echo ""
echo -e "${YELLOW}[2] Extracting test data...${NC}"

echo "📄 Files in vault:"
ls -1 "$VAULT_PATH"/*.md 2>/dev/null | while read f; do
  FILE=$(basename "$f")
  HASH=$(sha256sum "$f" | cut -d' ' -f1)
  echo "  - $FILE"
  echo "    Hash: ${HASH:0:32}..."
done

echo ""
echo "📝 Manifest entries:"
cat "$VAULT_PATH/manifest.json" | jq 'keys' 2>/dev/null || echo "  (jq not available, showing raw)"

echo ""
echo "👥 Trusted contributors:"
cat contributors.json | jq '.[].name' 2>/dev/null || cat contributors.json

echo ""
echo -e "${YELLOW}[3] Verifying signatures manually...${NC}"

if command -v jq &> /dev/null; then
  echo "Checking File.md signature:"
  ENTRY=$(cat "$VAULT_PATH/manifest.json" | jq '.["File.md"]' 2>/dev/null)
  STORED_HASH=$(echo "$ENTRY" | jq -r '.file_hash' 2>/dev/null)
  PUBKEY=$(echo "$ENTRY" | jq -r '.public_key' 2>/dev/null)
  
  ACTUAL_HASH=$(sha256sum "$VAULT_PATH/File.md" | cut -d' ' -f1)
  
  if [ "$STORED_HASH" = "$ACTUAL_HASH" ]; then
    echo -e "${GREEN}✓ Hash matches${NC}"
    echo "  Stored:  $STORED_HASH"
    echo "  Actual:  $ACTUAL_HASH"
  else
    echo -e "${RED}✗ Hash mismatch!${NC}"
    echo "  Stored:  $STORED_HASH"
    echo "  Actual:  $ACTUAL_HASH"
  fi
  
  if grep -q "$PUBKEY" contributors.json; then
    echo -e "${GREEN}✓ Public key is trusted${NC}"
  else
    echo -e "${RED}✗ Public key is NOT trusted${NC}"
  fi
else
  echo "jq not available - skipping detailed signature verification"
fi

echo ""
echo -e "${YELLOW}[4] Building test configuration...${NC}"
cat > "$CONFIG_FILE" << 'CONFIG'
{
  "addr": ":9999",
  "siteName": "RED Engine Integrity Test",
  "dataDir": "data/remote/API",
  "sourceURL": "",
  "sourceType": "",
  "adminToken": "test-token-12345",
  "startupSync": null
}
CONFIG
echo -e "${GREEN}✓ Created $CONFIG_FILE${NC}"

echo ""
echo -e "${YELLOW}[5] Build and Start Server${NC}"
if [ ! -f "red" ]; then
  echo "Building red-engine..."
  go build -o red ./cmd/red
  echo -e "${GREEN}✓ Built successfully${NC}"
else
  echo -e "${GREEN}✓ red binary already exists${NC}"
fi

echo ""
echo -e "${YELLOW}[6] Server Ready${NC}"
echo "Run the following to start the server and test:"
echo ""
echo -e "${YELLOW}cd $(pwd)${NC}"
echo -e "${YELLOW}./red -config $CONFIG_FILE${NC}"
echo ""
echo "Then visit: ${GREEN}http://localhost:$PORT/File${NC}"
echo ""
echo "Expected UI elements:"
echo "  ✅ Green badge: \"Verified Contributor: StandardCodebase\""
echo "  🔐 SHA-256 checksum at bottom of page"
echo ""
