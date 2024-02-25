#!/bin/bash

# Function to copy files from files directory to system directories
copy_files() {
    echo "Copying files..."
    cp -r files/* /
}

# Function to copy gwng executable based on architecture
copy_executable() {
    local arch=$(uname -m)
    local executable=""
    echo "Architecture: $arch"

    case $arch in
        "x86_64")
            executable="gwng-linux-amd64"
            ;;
        "i386" | "i486" | "i586" | "i686")
            executable="gwng-linux-386"
            ;;
        "armv7l")
            executable="gwng-linux-arm"
            ;;
        "aarch64")
            executable="gwng-linux-arm64"
            ;;
        "mips")
            executable="gwng-linux-mips"
            ;;
        "mips64")
            executable="gwng-linux-mips64"
            ;;
        "mips64el")
            executable="gwng-linux-mips64le"
            ;;
        "mipsel")
            executable="gwng-linux-mipsle"
            ;;
        "ppc64")
            executable="gwng-linux-ppc64"
            ;;
        "ppc64le")
            executable="gwng-linux-ppc64le"
            ;;
        "riscv64")
            executable="gwng-linux-riscv64"
            ;;
        "s390x")
            executable="gwng-linux-s390x"
            ;;
        *)
            echo "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    cp "exec/$executable" /usr/sbin/gwng
    chmod +x /usr/sbin/gwng
    chmod +x /etc/init.d/gwng_autologin
}

# Function to set gwng_autologin as startup service
set_startup_service() {
    echo "Setting gwng_autologin as startup service..."
    /etc/init.d/gwng_autologin enable
}

# Function to install
install() {
    copy_files
    copy_executable
    set_startup_service
    echo "Installation completed. Please refresh the web page to see the changes."
}

# Function to uninstall
uninstall() {
    echo "Uninstalling..."
    rm -rf /etc/config/gwng_autologin
    rm -rf /etc/init.d/gwng_autologin
    rm -rf /usr/lib/lua/luci/model/cbi/gwng_autologin.lua
    rm -rf /usr/lib/lua/luci/controller/gwng_autologin.lua
    rm -rf /usr/sbin/gwng
    rm -rf /var/log/gwng
    echo "Uninstallation completed. The system will restart in 5 seconds to complete the uninstallation process."
    sleep 5
    reboot
}

# Main script

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 {install|uninstall}"
    exit 1
fi

case "$1" in
    "install")
        install
        ;;
    "uninstall")
        uninstall
        ;;
    *)
        echo "Unknown option: $1"
        echo "Usage: $0 {install|uninstall}"
        exit 1
        ;;
esac

exit 0
