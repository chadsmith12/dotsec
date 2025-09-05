#!/usr/bin/env bash
# dotsec installer - Professional CLI installation script
# 
# This script installs the latest version of dotsec, a CLI tool for managing
# development secrets with Passbolt integration.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash -s -- --version v1.2.3
#   curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash -s -- --beta
#   curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | INSTALL_DIR=/opt/dotsec bash
#
# Options:
#   --version VERSION    Install specific version (default: latest stable)
#   --beta              Install latest pre-release/beta version
#   --no-completions     Skip shell completion setup
#   --dry-run           Show what would be installed without making changes
#   --uninstall         Remove dotsec installation
#   --help              Show this help message
#
# Environment Variables:
#   INSTALL_DIR         Custom installation directory (default: /usr/local/bin or ~/.local/bin)
#   GITHUB_TOKEN        GitHub token for API requests (optional)
#   NO_COLOR           Disable colored output
#
# Copyright (c) 2024 Chad Smith
# Licensed under MIT License

set -euo pipefail

# Script configuration
readonly SCRIPT_VERSION="1.0.0"
readonly GITHUB_REPO="chadsmith12/dotsec"
readonly BINARY_NAME="dotsec"
readonly INSTALL_TIMEOUT=300

# Color codes for output (disabled if NO_COLOR is set)
if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
    readonly RED='\033[0;31m'
    readonly GREEN='\033[0;32m'
    readonly YELLOW='\033[1;33m'
    readonly BLUE='\033[0;34m'
    readonly PURPLE='\033[0;35m'
    readonly CYAN='\033[0;36m'
    readonly WHITE='\033[1;37m'
    readonly BOLD='\033[1m'
    readonly NC='\033[0m' # No Color
else
    readonly RED=''
    readonly GREEN=''
    readonly YELLOW=''
    readonly BLUE=''
    readonly PURPLE=''
    readonly CYAN=''
    readonly WHITE=''
    readonly BOLD=''
    readonly NC=''
fi

# Global variables
VERSION="${VERSION:-}"
INSTALL_DIR="${INSTALL_DIR:-}"
DRY_RUN=false
NO_COMPLETIONS=false
UNINSTALL=false
BETA=false
TEMP_DIR=""

# Utility functions
log_info() {
    echo -e "${BLUE}${BOLD}==>${NC} ${WHITE}$1${NC}" >&2
}

log_success() {
    echo -e "${GREEN}${BOLD}✓${NC} ${WHITE}$1${NC}" >&2
}

log_warn() {
    echo -e "${YELLOW}${BOLD}⚠${NC} ${YELLOW}$1${NC}" >&2
}

log_error() {
    echo -e "${RED}${BOLD}✗${NC} ${RED}$1${NC}" >&2
}

log_debug() {
    if [[ "${DEBUG:-}" == "1" ]]; then
        echo -e "${PURPLE}${BOLD}DEBUG:${NC} $1" >&2
    fi
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
    echo -e "${NC}${WHITE}${BOLD}dotsec installer v${SCRIPT_VERSION}${NC}"
    echo -e "${CYAN}Secure development secrets management${NC}"
    echo
}

cleanup() {
    if [[ -n "$TEMP_DIR" ]] && [[ -d "$TEMP_DIR" ]]; then
        log_debug "Cleaning up temporary directory: $TEMP_DIR"
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get the latest version from GitHub API
get_latest_version() {
    log_debug "Fetching latest version from GitHub API"
    
    local auth_header=""
    if [[ -n "${GITHUB_TOKEN:-}" ]]; then
        auth_header="Authorization: token $GITHUB_TOKEN"
    fi
    
    local version_url
    if [[ "$BETA" == "true" ]]; then
        version_url="https://api.github.com/repos/$GITHUB_REPO/releases"
        log_debug "Fetching latest pre-release version"
    else
        version_url="https://api.github.com/repos/$GITHUB_REPO/releases/latest"
        log_debug "Fetching latest stable version"
    fi
    
    if command_exists curl; then
        local response
        response=$(curl -fsSL ${auth_header:+-H "$auth_header"} "$version_url" 2>/dev/null) || {
            log_error "Failed to fetch latest version from GitHub API"
            return 1
        }
        
        if [[ "$BETA" == "true" ]]; then
            # Find the first pre-release (which is the latest)
            echo "$response" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | while read -r tag; do
                if [[ "$tag" =~ -[a-zA-Z]+ ]]; then
                    echo "$tag"
                    break
                fi
            done | head -1
        else
            echo "$response" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | head -1
        fi
    elif command_exists wget; then
        local wget_args=()
        if [[ -n "$auth_header" ]]; then
            wget_args+=(--header="$auth_header")
        fi
        local response
        response=$(wget -qO- "${wget_args[@]}" "$version_url" 2>/dev/null) || {
            log_error "Failed to fetch latest version from GitHub API"
            return 1
        }
        
        if [[ "$BETA" == "true" ]]; then
            # Find the first pre-release (which is the latest)
            echo "$response" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | while read -r tag; do
                if [[ "$tag" =~ -[a-zA-Z]+ ]]; then
                    echo "$tag"
                    break
                fi
            done | head -1
        else
            echo "$response" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | head -1
        fi
    else
        log_error "Neither curl nor wget is available"
        return 1
    fi
}

# Detect operating system
detect_os() {
    local os
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          
            log_error "Unsupported operating system: $(uname -s)"
            return 1
            ;;
    esac
    echo "$os"
}

# Detect architecture
detect_arch() {
    local arch
    case "$(uname -m)" in
        x86_64|amd64)   arch="x86_64" ;;
        i386|i686)      arch="i386" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7l)         arch="armv7" ;;
        *)              
            log_error "Unsupported architecture: $(uname -m)"
            return 1
            ;;
    esac
    echo "$arch"
}

# Detect shell and profile file
detect_shell_profile() {
    local shell_name profile_file
    
    # Detect shell
    if [[ -n "${SHELL:-}" ]]; then
        shell_name=$(basename "$SHELL")
    else
        shell_name="bash"
    fi
    
    # Detect profile file based on shell
    case "$shell_name" in
        bash)
            if [[ -f "$HOME/.bash_profile" ]]; then
                profile_file="$HOME/.bash_profile"
            elif [[ -f "$HOME/.bashrc" ]]; then
                profile_file="$HOME/.bashrc"
            elif [[ -f "$HOME/.profile" ]]; then
                profile_file="$HOME/.profile"
            fi
            ;;
        zsh)
            if [[ -f "$HOME/.zshrc" ]]; then
                profile_file="$HOME/.zshrc"
            elif [[ -f "$HOME/.zprofile" ]]; then
                profile_file="$HOME/.zprofile"
            fi
            ;;
        fish)
            if [[ -d "$HOME/.config/fish" ]]; then
                profile_file="$HOME/.config/fish/config.fish"
                mkdir -p "$HOME/.config/fish"
            fi
            ;;
        *)
            shell_name="bash"
            profile_file="$HOME/.bashrc"
            ;;
    esac
    
    echo "$shell_name:${profile_file:-}"
}

# Determine installation directory
get_install_dir() {
    if [[ -n "$INSTALL_DIR" ]]; then
        echo "$INSTALL_DIR"
        return
    fi
    
    # Try system-wide installation first
    if [[ -w "/usr/local/bin" ]] || [[ "$EUID" == "0" ]]; then
        echo "/usr/local/bin"
    elif [[ -d "$HOME/.local/bin" ]] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        echo "$HOME/.local/bin"
    else
        echo "$HOME/bin"
        mkdir -p "$HOME/bin"
    fi
}

# Check if directory is in PATH
is_in_path() {
    local dir="$1"
    case ":$PATH:" in
        *":$dir:"*) return 0 ;;
        *) return 1 ;;
    esac
}

# Download file with progress bar
download_file() {
    local url="$1"
    local output="$2"
    local description="${3:-file}"
    
    log_info "Downloading $description..."
    log_debug "URL: $url"
    log_debug "Output: $output"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would download from $url"
        return 0
    fi
    
    if command_exists curl; then
        if curl --version 2>/dev/null | grep -q "curl 7\.[0-9]"; then
            # Older curl versions
            curl -fL --progress-bar "$url" -o "$output"
        else
            # Newer curl with better progress display
            curl -fL --progress-bar "$url" -o "$output"
        fi
    elif command_exists wget; then
        wget --progress=bar:force -O "$output" "$url"
    else
        log_error "Neither curl nor wget is available"
        return 1
    fi
}

# Verify checksum if available
verify_checksum() {
    local file="$1"
    local expected_checksum="$2"
    
    if [[ -z "$expected_checksum" ]]; then
        log_warn "No checksum provided, skipping verification"
        return 0
    fi
    
    log_info "Verifying checksum..."
    
    if command_exists sha256sum; then
        local actual_checksum
        actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
    elif command_exists shasum; then
        local actual_checksum
        actual_checksum=$(shasum -a 256 "$file" | cut -d' ' -f1)
    else
        log_warn "No checksum utility available, skipping verification"
        return 0
    fi
    
    if [[ "$actual_checksum" == "$expected_checksum" ]]; then
        log_success "Checksum verified"
        return 0
    else
        log_error "Checksum verification failed"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        return 1
    fi
}

# Install the binary
install_binary() {
    local temp_binary="$1"
    local install_dir="$2"
    
    log_info "Installing dotsec to $install_dir..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would install binary to $install_dir/$BINARY_NAME"
        return 0
    fi
    
    # Create install directory if it doesn't exist
    if ! mkdir -p "$install_dir"; then
        log_error "Failed to create installation directory: $install_dir"
        return 1
    fi
    
    # Check if we need sudo for installation
    local use_sudo=false
    if [[ ! -w "$install_dir" ]] && [[ "$EUID" != "0" ]]; then
        if command_exists sudo; then
            use_sudo=true
            log_info "Using sudo for installation to $install_dir"
        else
            log_error "No permission to write to $install_dir and sudo is not available"
            return 1
        fi
    fi
    
    # Backup existing installation
    local backup_file=""
    if [[ -f "$install_dir/$BINARY_NAME" ]]; then
        backup_file="$install_dir/$BINARY_NAME.backup.$(date +%s)"
        if [[ "$use_sudo" == "true" ]]; then
            sudo mv "$install_dir/$BINARY_NAME" "$backup_file"
        else
            mv "$install_dir/$BINARY_NAME" "$backup_file"
        fi
        log_info "Backed up existing installation to $(basename "$backup_file")"
    fi
    
    # Install the binary
    if [[ "$use_sudo" == "true" ]]; then
        if ! sudo cp "$temp_binary" "$install_dir/$BINARY_NAME"; then
            # Restore backup on failure
            if [[ -n "$backup_file" ]] && [[ -f "$backup_file" ]]; then
                sudo mv "$backup_file" "$install_dir/$BINARY_NAME"
            fi
            log_error "Failed to install binary"
            return 1
        fi
        sudo chmod +x "$install_dir/$BINARY_NAME"
    else
        if ! cp "$temp_binary" "$install_dir/$BINARY_NAME"; then
            # Restore backup on failure
            if [[ -n "$backup_file" ]] && [[ -f "$backup_file" ]]; then
                mv "$backup_file" "$install_dir/$BINARY_NAME"
            fi
            log_error "Failed to install binary"
            return 1
        fi
        chmod +x "$install_dir/$BINARY_NAME"
    fi
    
    # Clean up backup on success
    if [[ -n "$backup_file" ]] && [[ -f "$backup_file" ]]; then
        if [[ "$use_sudo" == "true" ]]; then
            sudo rm -f "$backup_file"
        else
            rm -f "$backup_file"
        fi
    fi
    
    log_success "dotsec installed successfully to $install_dir/$BINARY_NAME"
}

# Setup shell completions
setup_completions() {
    if [[ "$NO_COMPLETIONS" == "true" ]] || [[ "$DRY_RUN" == "true" ]]; then
        return 0
    fi
    
    local install_dir="$1"
    local shell_info
    shell_info=$(detect_shell_profile)
    local shell_name="${shell_info%%:*}"
    local profile_file="${shell_info##*:}"
    
    if [[ -z "$profile_file" ]]; then
        log_warn "Could not detect shell profile file, skipping completion setup"
        return 0
    fi
    
    log_info "Setting up shell completion for $shell_name..."
    
    # Check if dotsec is accessible
    local dotsec_cmd="$install_dir/$BINARY_NAME"
    if ! is_in_path "$install_dir"; then
        # Add to PATH temporarily for completion generation
        export PATH="$install_dir:$PATH"
    fi
    
    # Test if dotsec can generate completions
    if ! "$dotsec_cmd" completion "$shell_name" >/dev/null 2>&1; then
        log_warn "dotsec does not support completion for $shell_name, skipping"
        return 0
    fi
    
    local completion_line
    case "$shell_name" in
        bash)
            completion_line='eval "$(dotsec completion bash)"'
            ;;
        zsh)
            completion_line='eval "$(dotsec completion zsh)"'
            ;;
        fish)
            completion_line='dotsec completion fish | source'
            ;;
        *)
            log_warn "Unsupported shell for completions: $shell_name"
            return 0
            ;;
    esac
    
    # Check if completion is already configured
    if [[ -f "$profile_file" ]] && grep -q "dotsec completion" "$profile_file"; then
        log_success "Shell completion already configured in $profile_file"
        return 0
    fi
    
    # Add completion to profile file
    {
        echo ""
        echo "# dotsec shell completion"
        echo "$completion_line"
    } >> "$profile_file"
    
    log_success "Shell completion configured in $profile_file"
    log_info "Restart your shell or run 'source $profile_file' to enable completion"
}

# Verify installation
verify_installation() {
    local install_dir="$1"
    local expected_version="$2"
    
    log_info "Verifying installation..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_success "DRY RUN: Installation verification skipped"
        return 0
    fi
    
    local dotsec_path="$install_dir/$BINARY_NAME"
    
    # Check if binary exists and is executable
    if [[ ! -f "$dotsec_path" ]]; then
        log_error "Binary not found at $dotsec_path"
        return 1
    fi
    
    if [[ ! -x "$dotsec_path" ]]; then
        log_error "Binary is not executable: $dotsec_path"
        return 1
    fi
    
    # Check version
    local installed_version
    if ! installed_version=$("$dotsec_path" --version 2>/dev/null | grep -o "v[0-9].*" | head -1); then
        # Try alternative version format
        if ! installed_version=$("$dotsec_path" version 2>/dev/null | grep -o "v[0-9].*" | head -1); then
            log_warn "Could not determine installed version, but binary appears functional"
            return 0
        fi
    fi
    
    if [[ -n "$expected_version" ]] && [[ "$installed_version" != "$expected_version" ]]; then
        log_warn "Version mismatch: expected $expected_version, got $installed_version"
    fi
    
    log_success "Installation verified: dotsec $installed_version"
}

# Main installation function
install_dotsec() {
    local version="$1"
    
    # Detect platform
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    log_info "Detected platform: $os/$arch"
    
    # Get version
    if [[ -z "$version" ]]; then
        if [[ "$BETA" == "true" ]]; then
            log_info "Fetching latest beta version..."
        else
            log_info "Fetching latest stable version..."
        fi
        version=$(get_latest_version)
        if [[ -z "$version" ]]; then
            if [[ "$BETA" == "true" ]]; then
                log_error "Could not determine latest beta version. No pre-releases found."
            else
                log_error "Could not determine latest stable version"
            fi
            return 1
        fi
    fi
    
    if [[ "$BETA" == "true" ]] || [[ "$version" =~ -[a-zA-Z]+ ]]; then
        log_info "Installing dotsec $version (beta)"
    else
        log_info "Installing dotsec $version"
    fi
    
    # Construct download URL
    # Based on your .goreleaser.yaml, the format should be:
    # dotsec_Linux_x86_64.tar.gz, dotsec_Darwin_x86_64.tar.gz, etc.
    local os_title
    case "$os" in
        linux)   os_title="Linux" ;;
        darwin)  os_title="Darwin" ;;
        windows) os_title="Windows" ;;
    esac
    
    local archive_name="${BINARY_NAME}_${os_title}_${arch}"
    if [[ "$os" == "windows" ]]; then
        archive_name="${archive_name}.zip"
    else
        archive_name="${archive_name}.tar.gz"
    fi
    
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$version/$archive_name"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    local archive_path="$TEMP_DIR/$archive_name"
    
    # Download archive
    if ! download_file "$download_url" "$archive_path" "dotsec $version"; then
        log_error "Failed to download dotsec"
        return 1
    fi
    
    # Extract archive
    log_info "Extracting archive..."
    if [[ "$DRY_RUN" != "true" ]]; then
        if [[ "$archive_name" == *.zip ]]; then
            if command_exists unzip; then
                unzip -q "$archive_path" -d "$TEMP_DIR"
            else
                log_error "unzip is required to extract .zip files"
                return 1
            fi
        else
            tar -xzf "$archive_path" -C "$TEMP_DIR"
        fi
    fi
    
    local binary_path="$TEMP_DIR/$BINARY_NAME"
    if [[ "$os" == "windows" ]]; then
        binary_path="$TEMP_DIR/${BINARY_NAME}.exe"
    fi
    
    # Get installation directory
    local install_dir
    install_dir=$(get_install_dir)
    
    # Install binary
    if ! install_binary "$binary_path" "$install_dir"; then
        log_error "Failed to install binary"
        return 1
    fi
    
    # Setup shell completions
    setup_completions "$install_dir"
    
    # Verify installation
    if ! verify_installation "$install_dir" "$version"; then
        log_error "Installation verification failed"
        return 1
    fi
    
    # Show success message
    show_success_message "$install_dir" "$version"
}

# Show success message and next steps
show_success_message() {
    local install_dir="$1"
    local version="$2"
    
    echo
    if [[ "$BETA" == "true" ]] || [[ "$version" =~ -[a-zA-Z]+ ]]; then
        log_success "dotsec $version (beta) installed successfully!"
        echo -e "${YELLOW}${BOLD}⚠️  This is a beta version - use with caution in production!${NC}"
    else
        log_success "dotsec $version installed successfully!"
    fi
    echo
    
    # Check if install dir is in PATH
    if ! is_in_path "$install_dir"; then
        log_warn "Installation directory $install_dir is not in your PATH"
        echo -e "${YELLOW}Add it to your PATH by adding this line to your shell profile:${NC}"
        echo -e "${CYAN}export PATH=\"$install_dir:\$PATH\"${NC}"
        echo
    fi
    
    echo -e "${WHITE}${BOLD}Getting Started:${NC}"
    echo -e "${CYAN}  dotsec init${NC}       # Initialize project configuration"
    echo -e "${CYAN}  dotsec configure${NC}  # Set up Passbolt authentication"
    echo -e "${CYAN}  dotsec --help${NC}     # Show all available commands"
    echo
    
    echo -e "${WHITE}${BOLD}Documentation:${NC}"
    echo -e "${CYAN}  https://github.com/$GITHUB_REPO${NC}"
    echo
    
    if [[ "$NO_COMPLETIONS" != "true" ]]; then
        echo -e "${WHITE}${BOLD}Shell completion configured!${NC}"
        echo -e "Restart your shell or run: ${CYAN}source ~/.bashrc${NC} (or your shell's profile file)"
        echo
    fi
}

# Uninstall function
uninstall_dotsec() {
    log_info "Uninstalling dotsec..."
    
    local install_dirs=("/usr/local/bin" "$HOME/.local/bin" "$HOME/bin")
    local removed=false
    
    for dir in "${install_dirs[@]}"; do
        if [[ -f "$dir/$BINARY_NAME" ]]; then
            if [[ "$DRY_RUN" == "true" ]]; then
                log_info "DRY RUN: Would remove $dir/$BINARY_NAME"
                removed=true
            else
                if [[ -w "$dir" ]] || [[ "$EUID" == "0" ]]; then
                    rm -f "$dir/$BINARY_NAME"
                    log_success "Removed $dir/$BINARY_NAME"
                    removed=true
                elif command_exists sudo; then
                    sudo rm -f "$dir/$BINARY_NAME"
                    log_success "Removed $dir/$BINARY_NAME"
                    removed=true
                else
                    log_warn "Cannot remove $dir/$BINARY_NAME (no permission and no sudo)"
                fi
            fi
        fi
    done
    
    if [[ "$removed" != "true" ]]; then
        log_warn "dotsec installation not found"
        return 1
    fi
    
    log_info "Note: Shell completion configuration remains in your profile files"
    log_success "dotsec uninstalled successfully"
}

# Show help message
show_help() {
    show_header
    cat << EOF
${WHITE}${BOLD}USAGE:${NC}
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash -s -- [OPTIONS]

${WHITE}${BOLD}OPTIONS:${NC}
    ${CYAN}--version VERSION${NC}     Install specific version (default: latest stable)
    ${CYAN}--beta${NC}                Install latest pre-release/beta version
    ${CYAN}--no-completions${NC}      Skip shell completion setup
    ${CYAN}--dry-run${NC}             Show what would be installed without making changes
    ${CYAN}--uninstall${NC}           Remove dotsec installation
    ${CYAN}--help${NC}                Show this help message

${WHITE}${BOLD}ENVIRONMENT VARIABLES:${NC}
    ${CYAN}INSTALL_DIR${NC}           Custom installation directory
    ${CYAN}GITHUB_TOKEN${NC}          GitHub token for API requests (optional)
    ${CYAN}NO_COLOR${NC}              Disable colored output
    ${CYAN}DEBUG${NC}                 Enable debug output

${WHITE}${BOLD}EXAMPLES:${NC}
    # Install latest stable version
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash

    # Install latest beta version
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash -s -- --beta

    # Install specific version
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash -s -- --version v1.2.3

    # Install to custom directory
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | INSTALL_DIR=/opt/dotsec bash

    # Dry run to see what would be installed
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash -s -- --dry-run

    # Uninstall
    curl -fsSL https://raw.githubusercontent.com/$GITHUB_REPO/main/install.sh | bash -s -- --uninstall

${WHITE}${BOLD}SUPPORT:${NC}
    Report issues: ${CYAN}https://github.com/$GITHUB_REPO/issues${NC}

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --beta)
                BETA=true
                shift
                ;;
            --no-completions)
                NO_COMPLETIONS=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --uninstall)
                UNINSTALL=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                log_info "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

# Preflight checks
preflight_checks() {
    log_info "Performing preflight checks..."
    
    # Check for required commands
    if ! command_exists curl && ! command_exists wget; then
        log_error "Either curl or wget is required for installation"
        return 1
    fi
    
    if ! command_exists tar && ! command_exists unzip; then
        log_error "Either tar or unzip is required for extraction"
        return 1
    fi
    
    # Check internet connectivity
    if ! curl -fsSL --connect-timeout 5 https://api.github.com >/dev/null 2>&1; then
        if ! wget -q --timeout=5 -O /dev/null https://api.github.com >/dev/null 2>&1; then
            log_error "No internet connectivity or GitHub is not accessible"
            return 1
        fi
    fi
    
    log_success "Preflight checks passed"
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    # Show header
    if [[ "$DRY_RUN" != "true" ]] && [[ "$UNINSTALL" != "true" ]]; then
        show_header
    fi
    
    # Handle uninstall
    if [[ "$UNINSTALL" == "true" ]]; then
        uninstall_dotsec
        exit $?
    fi
    
    # Perform preflight checks
    if ! preflight_checks; then
        exit 1
    fi
    
    # Install dotsec
    if ! install_dotsec "$VERSION"; then
        log_error "Installation failed"
        exit 1
    fi
}

# Ensure the entire script is downloaded before execution
{
    main "$@"
}