
# Overview

Dotsec is a command-line interface (CLI) tool written in Go to simplify the process of downloading secrets from a password/secretmanager during development for your project. It is designed to streamline the sharing of secrets within your team, supporting either `dotnet user-secrets` or a `.env` file as the storage mechanism.

## Installation

You can install Dotsec by downloading the latest release from the [Releases](https://github.com/chadsmith12/dotsec/releases) or by building it from source. Make sure to add the `dotsec` binary to your system's PATH for easy access.

```shell
# Install Dotsec from source
git clone https://github.com/chadsmith12/dotsec.git
cd dotsec
go build
mv dotsec /usr/local/bin/ # Move to a directory in your PATH
```

## Supported Secret/Password Managers

Currently `dotsec` only supports Passbolt, though it could be expanded in the future to support other secret/password managers.

## Configuration

Before using `dotsec`, you need to configure it by providing the Passbolt server details, your private key file, and an optional password. You can configure `dotsec` using the following command:

```shell
dotsec configure
```

This command will prompt you for the Passbolt server URL, private key file path, and password. If you leave the password blank, `dotsec` will prompt you for it each time it is required.

## Basic Usage

`dotsec` has two main commands/usages that you will use to manager your secrets for development. When interacting with Passbolt your secrets must be stored in a folder.

* `pull` - This command will retrieve secrets from your secret/password manager and save them to the current secret type you're using (`dotsec user-secrets` or a `.env` file).
* `push` - This command will push secrets from your secret type (`dotnet user-secrets` or `.env` file) to the secret/password manager. Use it to add/update secrets in your secret manager.

### Basic `pull`

The basic usage of `dotsec` involves using the `pull` command to retrieve secrets and save them to either a `dotnet user-secrets` file or a `.env` file. Here's a simple example:

```shell
dotsec pull "mysecrets"
```

This command pulls secrets from the "mysecrets" folder from Passbolt and then runs `dotnet user-secrets set` for each secret found. Not supplying a flag for the type will default to `dotnet user-secrets`

### Basic `push`

The basic usage of `dotsec` and the `push` command allows you to quickly push up newly created or updated secrets to your secret manager. Just like the `pull` command it will default to `dotnet user-secrets`.

```shell
dotsec push "mysecrets"
```

## Advanced Usage

### Pull Command

The `pull` command retrieves secrets from Passbolt and saves them to the env type. It takes the folder name in your secret manager as the first argument:

- `folder name` (required): The name of the Passbolt folder containing the secrets you want to retrieve.

### Push Command

The `push` command create or updates secrets inside of your secret manager from your current secrets type. It takes the folder name in your secret manager as the first argument:

- `folder name` (required): The name of the Passbolt folder containing the secrets you want to retrieve.

Both the `pull` and `push` commands take the following flags:

- `--project (-p)` (optional): The path to the dotnet project where you want to sync the secrets. Defaults to the current directory. This flag is only valid with `--type dotnet`.

- `--file (-f)` (optional): The `.env` file where you want to save the secrets. Defaults to `.env` in the current directory. This flag is only valid with `--type env`.

- `--type` (optional): Defaults to `dotnet`. Specifies the type of secrets file you want to use. Use `dotnet` to use `dotnet user-secrets` or `env` to use a `.env` file.

### Examples

#### Pull secrets for a dotnet project:

```shell
dotsec pull "mysecrets" --project /path/to/dotnet/project --type dotnet
```

This command retrieves secrets from the "mysecrets" folder in Passbolt and runs `dotnet user-secrets set` inside the project directory specified. If no `secrets.json` file is found then it will run `dotnet user-secrets init` on the directory first.

#### Pull secrets and save to a custom .env file:

```shell
dotsec pull "mysecrets" --file .env.development --type env
```

This command retrieves secrets from the "mysecrets" folder in Passbolt and saves them to the custom `.env` file named ".env.development" in the current directory.

#### Push secrets from a dotnet project:

```shell
dotsec push "mysecrets" --project /path/to/dotnet/project --type dotnet
```

This command will update the secrets inside of the "mysecrets" folder from the secrets that the dotnet project, in the project directory specified.
