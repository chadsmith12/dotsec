
# dotsec

[![Release](https://img.shields.io/github/v/release/chadsmith12/dotsec)](https://github.com/chadsmith12/dotsec/releases)
[![License](https://img.shields.io/github/license/chadsmith12/dotsec)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/chadsmith12/dotsec)](https://goreportcard.com/report/github.com/chadsmith12/dotsec)

> **Secure development secrets management with Passbolt integration**

**dotsec** is a command-line interface (CLI) tool written in Go that simplifies the process of synchronizing secrets between your password manager and development environment. It streamlines secret sharing within development teams by supporting both `dotnet user-secrets` and `.env` file formats.

## ✨ Features

- 🔐 **Secure**: Direct integration with Passbolt for enterprise-grade secret management
- 🔄 **Bi-directional Sync**: Pull secrets from Passbolt or push local secrets to Passbolt
- 🛠️ **Multi-format Support**: Works with `dotnet user-secrets` and `.env` files
- 🚀 **Easy Setup**: Simple configuration and installation process
- 📦 **Cross-platform**: Available for Linux, macOS, and Windows
- 🔧 **Development Focused**: Designed specifically for development team workflows

## 📦 Installation

### Quick Install (Recommended)

Install the latest stable version:

```bash
curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash
```

### Beta Releases

Install the latest beta version to test new features:

```bash
curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash -s -- --beta
```

### Install Specific Version

```bash
curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash -s -- --version v1.2.3
```

### Custom Installation Directory

```bash
curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | INSTALL_DIR=/opt/dotsec bash
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/chadsmith12/dotsec/releases)
2. Extract the archive
3. Move the binary to a directory in your PATH:

```bash
# Linux/macOS
sudo mv dotsec /usr/local/bin/

# Or to user directory
mkdir -p ~/.local/bin
mv dotsec ~/.local/bin/
```

### Build from Source

```bash
git clone https://github.com/chadsmith12/dotsec.git
cd dotsec
go build -o dotsec
sudo mv dotsec /usr/local/bin/
```

## 🚀 Quick Start

### 1. Install dotsec

```bash
curl -fsSL https://raw.githubusercontent.com/chadsmith12/dotsec/main/install.sh | bash
```

### 2. Configure Passbolt Connection

```bash
dotsec configure
```

This will prompt you for:
- **Passbolt Server URL**: Your Passbolt instance URL
- **Private Key File**: Path to your Passbolt private key file
- **Password**: Optional password for the private key (leave blank to be prompted each time)

### 3. Initialize Project

```bash
dotsec init
```

This creates a project configuration file to manage your secret settings.

### 4. Start Using

```bash
# Pull secrets from Passbolt to your development environment
dotsec pull "my-project-secrets"

# Push local secrets to Passbolt
dotsec push "my-project-secrets"
```

## 🔧 Supported Secret Managers

| Manager | Status | Description |
|---------|--------|-------------|
| **Passbolt** | ✅ Supported | Enterprise-grade open source password manager |
| Others | 🔄 Planned | Additional managers may be supported in future releases |

## 📖 Usage

dotsec provides two primary commands for managing secrets between your development environment and Passbolt:

| Command | Description | Direction |
|---------|-------------|-----------|
| `pull` | Retrieve secrets from Passbolt | Passbolt → Local Environment |
| `push` | Upload secrets to Passbolt | Local Environment → Passbolt |

> **📝 Note**: When working with Passbolt, your secrets must be organized within folders.

### Basic Commands

#### Pull Secrets from Passbolt

```bash
# Pull secrets to dotnet user-secrets (default)
dotsec pull "my-project-secrets"

# Pull secrets to .env file
dotsec pull "my-project-secrets" --type env
```

#### Push Secrets to Passbolt

```bash
# Push secrets from dotnet user-secrets (default)
dotsec push "my-project-secrets"

# Push secrets from .env file
dotsec push "my-project-secrets" --type env
```

## 🔧 Advanced Usage

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
- Same as `pull` command

### 💡 Examples

#### .NET Development

```bash
# Pull secrets for current .NET project
dotsec pull "my-api-secrets" --type dotnet

# Pull secrets for specific .NET project
dotsec pull "my-api-secrets" --project /path/to/my-api --type dotnet

# Push local user-secrets to Passbolt
dotsec push "my-api-secrets" --project /path/to/my-api --type dotnet
```

> **📝 Note**: For .NET projects, if no `secrets.json` file exists, dotsec will automatically run `dotnet user-secrets init`.

#### Environment File Development

```bash
# Pull secrets to default .env file
dotsec pull "my-app-secrets" --type env

# Pull secrets to custom .env file
dotsec pull "my-app-secrets" --file .env.development --type env

# Push secrets from .env file to Passbolt
dotsec push "my-app-secrets" --file .env.local --type env
```

### 🛠️ Additional Commands

```bash
# Configure Passbolt connection
dotsec configure

# Initialize project configuration
dotsec init

# Run tests (development)
dotsec test

# View help
dotsec --help
```

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines for details on how to submit pull requests, report issues, and contribute to the project.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📚 [Documentation](https://github.com/chadsmith12/dotsec)
- 🐛 [Report Issues](https://github.com/chadsmith12/dotsec/issues)
- 💬 [Discussions](https://github.com/chadsmith12/dotsec/discussions)

## 📋 Requirements

- **Go**: Version 1.19 or higher (for building from source)
- **.NET SDK**: Required when using `--type dotnet` 
- **Passbolt**: Access to a Passbolt server instance

---

<div align="center">

**⭐ If you find dotsec useful, please consider giving it a star on GitHub! ⭐**

Made with ❤️ by [Chad Smith](https://github.com/chadsmith12)

</div>
