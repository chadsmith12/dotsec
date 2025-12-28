#!/usr/bin/env pwsh

# Script configuration
$SCRIPT_VERSION = "1.0.0"
$GITHUB_REPO = "chadsmith12/dotsec"
$BINARY_NAME = "dotsec"
$INSTALL_DIR = "$env:USERPROFILE\.dotsec\bin"

# Color codes for output
$COLOR_RESET = "`e[0m"
$COLOR_RED = "`e[31m"
$COLOR_GREEN = "`e[32m"
$COLOR_YELLOW = "`e[33m"
$COLOR_BLUE = "`e[34m"
$COLOR_CYAN = "`e[36m"
$COLOR_WHITE = "`e[37m"
$COLOR_BOLD = "`e[1m"

function Show-Header {
    Write-Host ""
    Write-Host "    ____        __  _____          " -ForegroundColor Cyan
    Write-Host "   / __ \____  / /_/ ___/___  _____" -ForegroundColor Cyan
    Write-Host "  / / / / __ \/ __/\__ \/ _ \/ ___/" -ForegroundColor Cyan
    Write-Host " / /_/ / /_/ / /_ ___/ /  __/ /__  " -ForegroundColor Cyan
    Write-Host "/_____/\____/\__//____/\___/\___/  " -ForegroundColor Cyan
    Write-Host ""
    Write-Host "dotsec installer v$SCRIPT_VERSION" -ForegroundColor White -NoNewline
    Write-Host " - Secure development secrets management" -ForegroundColor Cyan
    Write-Host ""
}

function Write-Info {
    param([string]$Message)
    Write-Host "$COLOR_BLUE==>$COLOR_RESET $COLOR_WHITE$Message$COLOR_RESET"
}

function Write-Success {
    param([string]$Message)
    Write-Host "$COLOR_GREEN✓$COLOR_RESET $COLOR_WHITE$Message$COLOR_RESET"
}

function Write-Error {
    param([string]$Message)
    Write-Host "$COLOR_RED✗$COLOR_RESET $COLOR_RED$Message$COLOR_RESET"
}

function Write-Warn {
    param([string]$Message)
    Write-Host "$COLOR_YELLOW⚠$COLOR_RESET $COLOR_YELLOW$Message$COLOR_RESET"
}

function Get-Architecture {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { "x86_64" }
        "ARM64" { "arm64" }
        default {
            Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
            exit 1
        }
    }
    return $arch
}

function Get-LatestVersion {
    Write-Info "Fetching latest version..."

    $apiUrl = "https://api.github.com/repos/$GITHUB_REPO/releases/latest"

    try {
        $response = Invoke-RestMethod -Uri $apiUrl -ErrorAction Stop
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to fetch latest version from GitHub API"
        Write-Host $_.Exception.Message -ForegroundColor Red
        exit 1
    }
}

function Install-Dotsec {
    param(
        [string]$Version,
        [string]$Arch
    )

    Write-Info "Installing dotsec $Version..."

    # Create temporary directory
    $tempDir = Join-Path $env:TEMP "dotsec-install-$([Guid]::NewGuid())"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    try {
        # Construct download URL
        $archiveName = "${BINARY_NAME}_Windows_${Arch}.zip"
        $downloadUrl = "https://github.com/$GITHUB_REPO/releases/download/${Version}/${archiveName}"
        $archivePath = Join-Path $tempDir $archiveName

        # Download binary
        Write-Info "Downloading dotsec..."
        try {
            Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -ErrorAction Stop
        }
        catch {
            Write-Error "Failed to download dotsec from $downloadUrl"
            Write-Host $_.Exception.Message -ForegroundColor Red
            exit 1
        }

        # Extract archive
        Write-Info "Extracting archive..."
        try {
            Expand-Archive -Path $archivePath -DestinationPath $tempDir -Force
        }
        catch {
            Write-Error "Failed to extract archive"
            Write-Host $_.Exception.Message -ForegroundColor Red
            exit 1
        }

        # Create installation directory
        Write-Info "Installing to $INSTALL_DIR..."
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null

        # Move binary to installation directory
        $binaryPath = Join-Path $tempDir "${BINARY_NAME}.exe"
        if (-not (Test-Path $binaryPath)) {
            Write-Error "Binary not found in archive: ${BINARY_NAME}.exe"
            exit 1
        }

        Copy-Item -Path $binaryPath -Destination "$INSTALL_DIR\" -Force

        # Add to PATH
        Add-ToPath

        # Verify installation
        Verify-Installation $Version

    }
    finally {
        # Clean up temporary directory
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

function Add-ToPath {
    Write-Info "Adding to PATH..."

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($currentPath -like "*$INSTALL_DIR*") {
        Write-Warn "PATH already contains $INSTALL_DIR"
    }
    else {
        $newPath = "$currentPath;$INSTALL_DIR"
        [Environment]::SetEnvironmentVariable("Path", $newPath, [EnvironmentVariableTarget]::User)
        Write-Success "Added $INSTALL_DIR to user PATH"
    }

    # Update current session
    $env:Path += ";$INSTALL_DIR"
}

function Verify-Installation {
    param([string]$ExpectedVersion)

    Write-Info "Verifying installation..."

    $binaryPath = Join-Path $INSTALL_DIR "${BINARY_NAME}.exe"

    if (-not (Test-Path $binaryPath)) {
        Write-Error "Binary not found at $binaryPath"
        exit 1
    }

    try {
        $installedVersion = & $binaryPath --version 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Installation verified: $installedVersion"
        }
        else {
            Write-Warn "Could not determine installed version, but binary appears functional"
            Write-Host "Exit code: $LASTEXITCODE, Output: $installedVersion" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Warn "Could not verify version, but binary was installed successfully"
        Write-Host "Error details: $($_.Exception.Message)" -ForegroundColor Red
    }
}

function Show-SuccessMessage {
    param([string]$Version)

    Write-Host ""
    Write-Success "dotsec $Version installed successfully!"
    Write-Host ""
    Write-Host "$COLOR_WHITE$COLOR_BOLDGetting Started:$COLOR_RESET"
    Write-Host "  $COLOR_CYAN  dotsec init$COLOR_RESET       # Initialize project configuration"
    Write-Host "  $COLOR_CYAN  dotsec configure$COLOR_RESET  # Set up Passbolt authentication"
    Write-Host "  $COLOR_CYAN  dotsec --help$COLOR_RESET     # Show all available commands"
    Write-Host ""
    Write-Host "$COLOR_WHITE$COLOR_BOLD Documentation:$COLOR_RESET"
    Write-Host "  $COLOR_CYAN  https://github.com/$GITHUB_REPO$COLOR_RESET"
    Write-Host ""
    Write-Warn "Restart your PowerShell terminal or run:"
    Write-Host "  $COLOR_CYAN$`$env:Path += `";$INSTALL_DIR`"$COLOR_RESET"
    Write-Host ""
}

# Main execution
Show-Header

# Preflight checks
Write-Info "Performing preflight checks..."

# Check PowerShell version
$psVersion = $PSVersionTable.PSVersion
Write-Info "PowerShell version: $psVersion"

# Check internet connectivity
try {
    $null = Invoke-WebRequest -Uri "https://api.github.com" -Method Head -TimeoutSec 5 -ErrorAction Stop
}
catch {
    Write-Error "No internet connectivity or GitHub is not accessible"
    exit 1
}

Write-Success "Preflight checks passed"

# Detect platform
$arch = Get-Architecture
Write-Info "Detected platform: Windows/$arch"

# Get latest version
$version = Get-LatestVersion

# Install dotsec
Install-Dotsec -Version $version -Arch $arch

# Show success message
Show-SuccessMessage -Version $version
