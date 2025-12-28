#!/usr/bin/env bash
set -euo pipefail

# dotsec installer for Linux and macOS
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash
#
# Requirements:
#   - curl or wget
#   - tar or unzip
#   - Internet connection to download from GitHub
#
# Copyright (c) 2024 Chad Smith
# Licensed under MIT License

GITHUB_REPO="chadsmith12/dotsec"
BINARY_NAME="dotsec"

# Color codes for output
if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
    readonly RESET='\033[0m'
    readonly RED='\033[31m'
    readonly GREEN='\033[32m'
    readonly YELLOW='\033[33m'
    readonly BLUE='\033[34m'
    readonly CYAN='\033[36m'
    readonly BOLD='\033[1m'
else
    readonly RESET=''
    readonly RED=''
    readonly GREEN=''
    readonly YELLOW=''
    readonly BLUE=''
    readonly CYAN=''
    readonly BOLD=''
fi

log_info() {
    echo -e "${BLUE}==>${RESET} $1"
}

log_success() {
    echo -e "${GREEN}✓${RESET} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${RESET} ${YELLOW}$1${RESET}"
}

log_error() {
    echo -e "${RED}✗${RESET} ${RED}$1${RESET}"
}

show_header() {
    echo -e "${CYAN}${BOLD}"
    cat << 'EOF'
    ____        __  _____          
   / __ \____  / /_/ ___/___  _____
  / / / / __ \/ __/\__ \/ _ \/ ___/
 / /_/ / /_/ / /_ ___/ /  __/ /__  
/_____/\____/\__//____/\___/\___/  
                                    
EOF
    echo -e "${RESET}${BOLD}dotsec installer${RESET}"
    echo -e "${CYAN}Secure development secrets management${RESET}"
    echo
}

detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "x86_64" ;;
        arm64|aarch64)  echo "arm64" ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

get_latest_version() {
    log_info "Fetching latest version..."
    
    local auth_header=""
    if [[ -n "${GITHUB_TOKEN:-}" ]]; then
        auth_header="Authorization: token $GITHUB_TOKEN"
    fi
    
    local version_url="https://api.github.com/repos/$GITHUB_REPO/releases/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL ${auth_header:+-H "$auth_header"} "$version_url" 2>/dev/null | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | head -1
    elif command -v wget >/dev/null 2>&1; then
        local wget_args=()
        [[ -n "$auth_header" ]] && wget_args+=(--header="$auth_header")
        wget -qO- "${wget_args[@]}" "$version_url" 2>/dev/null | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | head -1
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
}

get_install_dir() {
    local install_dir=""
    
    # Try system-wide installation first
    if [[ -w "/usr/local/bin" ]] || [[ "$EUID" == "0" ]]; then
        install_dir="/usr/local/bin"
    # Fall back to user directory
    elif [[ -d "$HOME/.local/bin" ]] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        install_dir="$HOME/.local/bin"
    else
        install_dir="$HOME/bin"
        mkdir -p "$HOME/bin"
    fi
    
    echo "$install_dir"
}

is_in_path() {
    case ":$PATH:" in
        *":$1:"*) return 0 ;;
        *) return 1 ;;
    esac
}

install_binary() {
    local version="$1"
    local os="$2"
    local arch="$3"
    
    # Determine OS title for archive name
    local os_title
    case "$os" in
        linux)   os_title="Linux" ;;
        darwin)  os_title="Darwin" ;;
    esac
    
    # Construct archive name and download URL
    local archive_name="${BINARY_NAME}_${os_title}_${arch}"
    if [[ "$os" == "linux" ]] || [[ "$os" == "darwin" ]]; then
        archive_name="${archive_name}.tar.gz"
    fi
    
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$version/$archive_name"
    
    log_info "Installing dotsec $version..."
    
    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT
    
    # Download archive
    log_info "Downloading dotsec..."
    if ! download_file "$download_url" "$temp_dir/$archive_name"; then
        log_error "Failed to download dotsec"
        exit 1
    fi
    
    # Extract archive
    log_info "Extracting archive..."
    if [[ "$archive_name" == *.tar.gz ]]; then
        tar -xzf "$temp_dir/$archive_name" -C "$temp_dir"
    else
        log_error "Unsupported archive format"
        exit 1
    fi
    
    # Get installation directory
    local install_dir
    install_dir=$(get_install_dir)
    
    log_info "Installing to $install_dir..."
    
    # Backup existing installation
    if [[ -f "$install_dir/$BINARY_NAME" ]]; then
        local backup_file="$install_dir/$BINARY_NAME.backup.$(date +%s)"
        mv "$install_dir/$BINARY_NAME" "$backup_file"
        log_info "Backed up existing installation to $(basename "$backup_file")"
    fi
    
    # Install binary
    if [[ ! -w "$install_dir" ]] && [[ "$EUID" != "0" ]]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo cp "$temp_dir/$BINARY_NAME" "$install_dir/$BINARY_NAME"
            sudo chmod +x "$install_dir/$BINARY_NAME"
        else
            log_error "No permission to write to $install_dir and sudo is not available"
            exit 1
        fi
    else
        cp "$temp_dir/$BINARY_NAME" "$install_dir/$BINARY_NAME"
        chmod +x "$install_dir/$BINARY_NAME"
    fi
    
    log_success "dotsec installed to $install_dir/$BINARY_NAME"
    
    # Verify installation
    verify_installation "$install_dir" "$version"
    
    # Show success message
    show_success_message "$install_dir" "$version"
}

download_file() {
    local url="$1"
    local output="$2"
    
    if command -v curl >/dev/null 2>&1; then
        curl -fL --progress-bar "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget --progress=bar:force -O "$output" "$url"
    else
        log_error "Neither curl nor wget is available"
        return 1
    fi
}

verify_installation() {
    local install_dir="$1"
    local expected_version="$2"
    
    log_info "Verifying installation..."
    
    local binary_path="$install_dir/$BINARY_NAME"
    
    if [[ ! -f "$binary_path" ]]; then
        log_error "Binary not found at $binary_path"
        exit 1
    fi
    
    if [[ ! -x "$binary_path" ]]; then
        log_error "Binary is not executable: $binary_path"
        exit 1
    fi
    
    # Check version
    local installed_version
    if ! installed_version=$("$binary_path" --version 2>/dev/null | grep -o "v[0-9].*" | head -1); then
        log_warn "Could not determine installed version, but binary appears functional"
        return 0
    fi
    
    if [[ -n "$expected_version" ]] && [[ "$installed_version" != "$expected_version" ]]; then
        log_warn "Version mismatch: expected $expected_version, got $installed_version"
    fi
    
    log_success "Installation verified: $installed_version"
}

show_success_message() {
    local install_dir="$1"
    local version="$2"
    
    echo
    log_success "dotsec $version installed successfully!"
    echo
    
    # Check if install dir is in PATH
    if ! is_in_path "$install_dir"; then
        log_warn "Installation directory $install_dir is not in your PATH"
        echo -e "${YELLOW}Add it to your PATH by adding this line to your shell profile:${RESET}"
        echo -e "${CYAN}export PATH=\"$install_dir:\$PATH\"${RESET}"
        echo
    fi
    
    echo -e "${BOLD}Getting Started:${RESET}"
    echo -e "${CYAN}  dotsec init${RESET}       # Initialize project configuration"
    echo -e "${CYAN}  dotsec configure${RESET}  # Set up Passbolt authentication"
    echo -e "${CYAN}  dotsec --help${RESET}     # Show all available commands"
    echo
    
    echo -e "${BOLD}Documentation:${RESET}"
    echo -e "${CYAN}  https://github.com/$GITHUB_REPO${RESET}"
    echo
}

# Preflight checks
preflight_checks() {
    log_info "Performing preflight checks..."
    
    # Check for required commands
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        log_error "Either curl or wget is required for installation"
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        log_error "tar is required for extraction"
        exit 1
    fi
    
    # Check internet connectivity
    if ! curl -fsSL --connect-timeout 5 https://api.github.com >/dev/null 2>&1; then
        if ! wget -q --timeout=5 -O /dev/null https://api.github.com >/dev/null 2>&1; then
            log_error "No internet connectivity or GitHub is not accessible"
            exit 1
        fi
    fi
    
    log_success "Preflight checks passed"
}

# Main function
main() {
    show_header
    
    # Perform preflight checks
    if ! preflight_checks; then
        exit 1
    fi
    
    # Detect platform
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    log_info "Detected platform: $os/$arch"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    if [[ -z "$version" ]]; then
        log_error "Could not determine latest version"
        exit 1
    fi
    
    log_info "Installing dotsec $version"
    
    # Install binary
    install_binary "$version" "$os" "$arch"
}

main
