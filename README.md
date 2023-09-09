## dotsec

`dotsec` is a CLI application written in Go that can be used to quickly download secrets into your dotnet projects local secrets. Secret managment can be complicated and when working on a team someone must have all those secrets somewhere that you need. The basic idea behind this tool is that your secrets live in a password manager and you give users that need access to secrets to the password manager. Running this will authenticate them with the password manager and automatically download the secrets and store them in the dotnet projects `secrets.json` by internally running `dotnet user-secrets set` to set the secret.


