#!/usr/bin/env bash
set -euxo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
enableRPiHIDPath=/opt/enable-rpi-hid
enableRPiHIDPathDir=$(dirname enableRPiHIDPath)

loadKernelModules() {
    # Enable device tree overlay, dwc2
    if ! grep 'dtoverlay=dwc2' /boot/config; then
        echo "dtoverlay=dwc2" >> /boot/config.txt
    fi

    if ! grep dwc2 /etc/modules; then
        echo "dwc2" >> /etc/modules
    fi

    # LibComposite puts usb config into userland (ConfigFS).
    if ! grep libcomposite /etc/modules; then
        echo "libcomposite" | tee --append /etc/modules
        #modprobe libcomposite       # Load kernel module while system is up
    fi
}

loadKernelModules

cp -v "$SCRIPT_DIR/enable-rpi-hid" "$enableRPiHIDPathDir"

cat <<EOF | tee /lib/systemd/system/usb-keyboard-gadget.service
[Unit]
Description=Create virtual keyboard USB gadget
After=syslog.target

[Service]
Type=oneshot
User=root
ExecStart=${enableRPiHIDPath}

[Install]
WantedBy=local-fs.target
EOF

systemctl daemon-reload
systemctl enable usb-gadget.service
