# dotsec

[![Release](https://img.shields.io/github/v/release/chadsmith12/dotsec)](https://github.com/chadsmith12/dotsec/releases)
[![License](https://img.shields.io/github/license/chadsmith12/dotsec)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/chadsmith12/dotsec)](https://goreportcard.com/report/github.com/chadsmith12/dotsec)

> **Secure development secrets management with Passbolt integration**

**dotsec** is a command-line interface (CLI) tool written in Go that simplifies the process of synchronizing secrets between your password manager (Passbolt) and local development environments. It streamlines secret sharing within development teams by supporting both `dotnet user-secrets` and `.env` file formats.

## Overview

The tool is designed to be used by teams where it can be difficult to keep and share development secrets across the team. Store them in the Passbolt password manager and quickly download and place them in your `secrets.json` file (for .NET) or `.env` file (for general development). When using dotnet, dotsec will run the `dotnet user-secrets` command, and for .env files, it will parse and save the secrets.

## Key Features

- **Secure**: Direct integration with Passbolt for enterprise-grade secret management
- **Bi-directional Sync**: Pull secrets from Passbolt or push local secrets to Passbolt
- **Multi-format Support**: Works with `dotnet user-secrets` and `.env` files
- **Easy Setup**: Simple configuration and installation process
- **Cross-platform**: Available for Linux, macOS, and Windows
- **Development Focused**: Designed specifically for development team workflows
- **Team-based Access Control**: Secrets can be shared with specific teams/groups in Passbolt

## Prerequisites

### Passbolt Setup

Before using dotsec, you need to have a Passbolt server set up with the following:

1. **Passbolt Server**: A running Passbolt instance (self-hosted or cloud)
2. **User Account**: A Passbolt user with permissions to access secrets
3. **Private Key**: Your Passbolt private key file
4. **Team/Group Structure**: Secrets organized in folders and shared with teams

### Local Development Environment

- **For .NET projects**: .NET SDK installed
- **For .env files**: Any development environment that uses environment variables

## Installation

### Manual Installation

This section provides step-by-step instructions for installing dotsec.

#### Windows

1. **Download the latest release**
   - Visit [GitHub Releases](https://github.com/chadsmith12/dotsec/releases)
   - Download the latest `.zip` file (e.g., `dotsec_Windows_x64.zip`)

2. **Extract the archive**
   - Create a directory for dotsec (e.g., `C:\Tools\dotsec`)
   - Extract the zip file contents into this directory
   - You should now have `dotsec.exe` in the directory

3. **Add to PATH**
   - Press `Windows + R`, type `sysdm.cpl`, press Enter
   - Go to the "Advanced" tab
   - Click "Environment Variables"
   - Under "System variables", find "Path" and click "Edit"
   - Click "New" and add the directory containing `dotsec.exe` (e.g., `C:\Tools\dotsec`)
   - Click "OK" on all dialogs
   - **Note**: You may need to restart your command prompt for changes to take effect

4. **Verify installation**
   ```cmd
   dotsec --version
   ```

#### Linux / macOS

1. **Download the latest release**
   - Visit [GitHub Releases](https://github.com/chadsmith12/dotsec/releases)
   - Download the latest archive file (e.g., `dotsec_Linux_x64.tar.gz` or `dotsec_Darwin_x64.tar.gz`)

2. **Extract the archive**

   ```bash
   # Create installation directory
   sudo mkdir -p /usr/local/dotsec
   sudo chown $USER:$USER /usr/local/dotsec

   # Extract to a temporary location first
   tar -xzf dotsec_Linux_x64.tar.gz -C /tmp

   # Move the binary to the installation directory
   sudo mv /tmp/dotsec /usr/local/dotsec/
   ```

3. **Add to PATH**
   - Edit your shell configuration file (`~/.bashrc`, `~/.zshrc`, or `~/.profile`)
   - Add the following line:
     ```bash
     export PATH="$PATH:/usr/local/dotsec"
     ```

   ````
   - Reload the configuration:
     ```bash
   source ~/.bashrc  # or ~/.zshrc, ~/.profile
   ````

4. **Verify installation**
   ```bash
   dotsec --version
   ```

### Build from Source

If you prefer to build from source, follow these steps:

1. **Clone the repository**

   ```bash
   git clone https://github.com/chadsmith12/dotsec.git
   cd dotsec
   ```

2. **Build the binary**

   ```bash
   go build -o dotsec
   ```

3. **Install the binary**

   ```bash
   # Option 1: Install system-wide (requires sudo)
   sudo mv dotsec /usr/local/bin/

   # Option 2: Install to user directory
   mkdir -p ~/.local/bin
   mv dotsec ~/.local/bin/
   ```

4. **Verify installation**
   ```bash
   dotsec --version
   ```

### Alternative: Install to User Directory (All Platforms)

If you don't have permission to install system-wide, you can install to your user directory:

1. **Download and extract**
   - Follow steps 1-2 from your platform's manual installation
   - Extract to `~/bin` or `~/.local/bin`

2. **Ensure directory exists**

   ```bash
   mkdir -p ~/bin
   ```

3. **Move the binary**

   ```bash
   mv dotsec ~/bin/
   ```

4. **Add to PATH** (if not already in PATH)
   - Edit your shell configuration file and add:
     ```bash
     export PATH="$PATH:~/bin"
     ```

5. **Reload shell configuration**

   ```bash
   source ~/.bashrc  # or ~/.zshrc, ~/.profile
   ```

6. **Verify installation**
   ```bash
   dotsec --version
   ```

### Verification

After installation, you should see the version information:

```bash
dotsec --version
```

If you encounter any issues, check that:

- The `dotsec` executable is in your PATH
- The binary is executable (Linux/macOS: `chmod +x dotsec`)
- You have the required permissions

### Troubleshooting

**Binary not found**

- Run `which dotsec` to check if it's in your PATH
- Run `echo $PATH` to verify your PATH environment variable

**Permission denied**

- On Linux/macOS, ensure the binary is executable: `chmod +x dotsec`
- On Windows, ensure you're running Command Prompt as administrator if installing system-wide

**Wrong version**

- Check that you have the correct binary for your platform (Windows, Linux, or macOS)
- Compare the file properties with other binaries in your PATH

## Quick Start

### Step 1: Configure Passbolt Connection

```bash
dotsec configure
```

This will prompt you for:

- **Passbolt Server URL**: Your Passbolt instance URL (e.g., `https://passbolt.example.com`)
- **Private Key File**: Path to your Passbolt private key file (e.g., `~/.passbolt/etc/key`)
- **Password**: Optional password for the private key (leave blank to be prompted each time)

### Step 2: Initialize Project

```bash
dotsec init
```

This creates a `.dotsecrc` file in your project root that stores:

- **Folder**: The Passbolt folder name containing your secrets
- **Type**: Secret storage format (`dotnet` or `env`)
- **Path**: The target file or project directory
- **Team**: The Passbolt team name for access control

### Step 3: Start Using

```bash
# Pull secrets from Passbolt to your development environment
dotsec pull "my-project-secrets"

# Push local secrets to Passbolt
dotsec push "my-project-secrets"
```

## Basic Usage

### Pull Secrets (Passbolt → Local Environment)

#### For .NET Projects

```bash
# Pull secrets to current .NET project's secrets.json
dotsec pull "my-api-secrets"

# Pull secrets to specific .NET project directory
dotsec pull "my-api-secrets" --project /path/to/my-api

# Pull secrets to custom .env file
dotsec pull "my-app-secrets" --type env --file .env.development
```

#### For Environment Files

```bash
# Pull secrets to default .env file
dotsec pull "my-app-secrets" --type env

# Pull secrets to custom .env file
dotsec pull "my-app-secrets" --file .env.local --type env
```

### Push Secrets (Local Environment → Passbolt)

#### From .NET Projects

```bash
# Push secrets from current .NET project
dotsec push "my-api-secrets"

# Push secrets from specific .NET project
dotsec push "my-api-secrets" --project /path/to/my-api
```

#### From Environment Files

```bash
# Push secrets from .env file
dotsec push "my-app-secrets" --type env

# Push secrets from custom .env file
dotsec push "my-app-secrets" --file .env.local --type env
```

### Additional Commands

```bash
# Configure Passbolt connection (run once)
dotsec configure

# Initialize project configuration (run once per project)
dotsec init

# See the current configuration dotsec will use
dotsec test

# List team members
dotsec team

# Add a user to a team
dotsec team add user@example.com

# View help
dotsec --help
```

## Configuration

### Configuration File

dotsec uses the XDG Base Directory specification for configuration:

- **Global config**: `~/.config/dotsec/.config` (JSON)
- **Project config**: `.dotsecrc` (JSON) in project root

### Configuration Options

#### Global Configuration (`~/.config/dotsec/.config`)

- `server`: Passbolt server URL
- `privateKey`: Path to Passbolt private key file
- `password`: Password for the private key (optional)

#### Project Configuration (`.dotsecrc`)

- `folder`: Passbolt folder name containing secrets
- `type`: Secret storage format (`dotnet` or `env`)
- `path`: Target file or project directory
- `team`: Passbolt team name for access control

### Environment Variables

You can also set configuration using environment variables with the `DOTSEC_` prefix:

```bash
export DOTSEC_SERVER=https://passbolt.example.com
export DOTSEC_PRIVATEKEY=~/.passbolt/etc/key
export DOTSEC_PASSWORD=your_password
```

## Passbolt Setup Requirements

### Folder Structure

Secrets in Passbolt should be organized in folders:

```
/my-secrets/
├── api/
│   ├── database.json
│   ├── api-key.json
│   └── config.json
├── frontend/
│   ├── env-config.json
│   └── api-endpoints.json
└── shared/
    └── common-secrets.json
```

### Team-Based Access Control

1. **Create Teams**: In Passbolt, create teams for different project groups
2. **Share Folders**: Share secret folders with the appropriate teams
3. **Set Permissions**: Ensure users have the correct permissions (read/update)

Example team structure:

- `api-team`: Access to API secrets
- `frontend-team`: Access to frontend secrets
- `shared-team`: Access to common secrets

### Secret Format

Each secret in Passbolt should be stored as a key-value pair:

- **Name**: The secret key (e.g., `DATABASE_URL`)
- **Value**: The secret value (e.g., `postgresql://localhost/db`)

## Advanced Usage

### Command Reference

#### `dotsec pull <folder-name>`

Retrieves secrets from a Passbolt folder and saves them to your local environment.

**Arguments:**

- `folder-name` (required): The name of the Passbolt folder containing your secrets

**Flags:**

- `--project, -p` (optional): Path to the dotnet project directory (default: current directory)
  - Only valid with `--type dotnet`
- `--file, -f` (optional): Target `.env` file path (default: `.env`)
  - Only valid with `--type env`
- `--type` (optional): Secret storage format (default: `dotnet`)
  - Values: `dotnet` | `env`

#### `dotsec push <folder-name>`

Uploads secrets from your local environment to a Passbolt folder.

**Arguments:**

- `folder-name` (required): The name of the Passbolt folder to update

**Flags:**

- `--project, -p` (optional): Path to the dotnet project directory
  - Only valid with `--type dotnet`
- `--file, -f` (optional): Target `.env` file path
  - Only valid with `--type env`
- `--type` (optional): Secret storage format (default: `dotnet`)
  - Values: `dotnet` | `env`
- `--prune` (optional): Delete secrets in Passbolt that don't exist locally
- `--force` (optional): Skip dirty check warning

### Examples

#### Complete Workflow

```bash
# 1. Configure Passbolt connection (once)
dotsec configure

# 2. Initialize project (once per project)
dotsec init

# 3. Pull secrets to current .NET project
dotsec pull "api-secrets" --type dotnet

# 4. Use secrets in your .NET application

# 5. Update secrets locally (e.g., in .env file)
echo "NEW_SECRET=value" >> .env

# 6. Push updated secrets back to Passbolt
dotsec push "api-secrets" --type env
```

#### Team Collaboration

```bash
# Pull secrets for the API team
dotsec pull "production-api" --project ./api --type dotnet

# Pull secrets for the frontend team
dotsec pull "frontend-config" --file ./frontend/.env --type env

# Push secrets from the API team
dotsec push "production-api" --project ./api --type dotnet
```

### Team Management

#### List Team Members

```bash
dotsec team
```

#### Add User to Team

```bash
dotsec team add user@example.com --manager
```

## Troubleshooting

### Common Issues

#### "folder is required"

Ensure you have either:

1. Run `dotsec init` to create a `.dotsecrc` file
2. Provide the folder name as an argument

#### "Failed to login to Passbolt"

Check:

1. Your Passbolt server URL is correct
2. Your private key file exists and is readable
3. Your password is correct

#### "Permission denied"

Ensure:

1. Your user is a member of the required team
2. The team has access to the folder
3. Your user has the correct permissions (read/update)

#### "No secrets configured"

This error occurs when:

1. The folder exists but has no secrets
2. The folder doesn't exist in Passbolt
3. You don't have access to the folder

#### File Permission Errors

If you encounter permission errors:

1. Ensure the target file/directory is writable
2. On Unix systems, check file permissions (e.g., `chmod 600 .env`)
3. Run with appropriate privileges

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
 go test ./env/
 go test ./passbolt/
 go test ./dotnet/
```

### Building

```bash
# Build the binary
go build -o dotsec

# Run the CLI
./dotsec --help
```

### Code Generation

```bash
# Run go generate (used by goreleaser)
go generate ./...

# Tidy go.mod
go mod tidy
```
