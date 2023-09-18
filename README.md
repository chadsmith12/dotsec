
# Overview

Dotsec is a command-line interface (CLI) tool written in Go to simplify the process of downloading secrets from Passbolt during development for your project. It is designed to streamline the sharing of secrets within your team, supporting either `dotnet user-secrets` or a `.env` file as the storage mechanism.

## Installation

You can install Dotsec by downloading the latest release from the [GitHub repository](https://github.com/chadsmith12/dotsec) or by building it from source. Make sure to add the Dotsec binary to your system's PATH for easy access.

```shell
# Install Dotsec from source
git clone https://github.com/chadsmith12/dotsec.git
cd dotsec
go build
mv dotsec /usr/local/bin/ # Move to a directory in your PATH
```

## Configuration

Before using Dotsec, you need to configure it by providing the Passbolt server details, your private key file, and an optional password. You can configure Dotsec using the following command:

```shell
dotsec configure
```

This command will prompt you for the Passbolt server URL, private key file path, and password. If you leave the password blank, Dotsec will prompt you for it each time it is required.

## Basic Usage

The basic usage of Dotsec involves using the `pull` command to retrieve secrets and save them to either a `dotnet user-secrets` file or a `.env` file. Here's a simple example:

```shell
dotsec pull "mysecrets"
```

This command pulls secrets from the "mysecrets" folder from Passbolt and then runs `dotnet user-secrets set` for each secret found.

## Advanced Usage

### Pull Command

The `pull` command retrieves secrets from Passbolt and saves them to the env type. It accepts the following arguments and flags:

- `folder name` (required): The name of the Passbolt folder containing the secrets you want to retrieve.

Flags:

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


