# uauth

Unauth is a reverse proxy that allowes IP based access to a docker registry.

# Setup

Create a `/etc/harbor-unauth.json` like this:

```json
{
  "url": "https://harbor.example.com/",
  "auths": {
    "apps": {
      "user": "<robot user>",
      "password": "<robot token>"
    },
    "images": {
      "user": "<robot user>",
      "password": "<robot token>"
    }
  }
}
```

Per default the the configuration is stored in `/etc/harbor-unauth.json`,
overwrite using the `UNAUTH_CONFIG` environment variable.

# Development

Setup a `config.json` and run `UNAUTH_CONFIG=./config.json go run .` Make sure
your configuration name is included in `.gitignore`.
