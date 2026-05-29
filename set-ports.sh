#!/bin/bash
echo "========================================"
echo "🔓 Configuring Permanent Unprivileged Ports"
echo "========================================"

# 1. Check for root/sudo privileges (Required to modify sysctl)
if [[ $EUID -ne 0 ]]; then
   echo "⚠️  Sudo privileges are required to modify system network settings."
   echo "Please enter your password when prompted."
   exec sudo "$0" "$@"
fi

# 2. Define the configuration file path
# Using the /etc/sysctl.d/ directory is the modern, persistent way to configure kernel parameters
CONF_FILE="/etc/sysctl.d/80-unprivileged-ports.conf"

echo "[*] Setting unprivileged port start to 80..."

# 3. Write the rule to the persistent boot file
cat <<EOF > "$CONF_FILE"
# Project R.E.D. - Allow rootless Podman to bind to standard web ports
net.ipv4.ip_unprivileged_port_start=80
EOF

# 4. Reload all system configuration files so it takes effect immediately without a reboot
echo "[*] Applying changes immediately..."
sysctl --system > /dev/null

echo "========================================"
echo "✅ Success! Ports 80 and above can now be bound by standard users."
echo "✅ This setting is now permanently baked into your boot configuration."
echo "========================================"