# uauth

Unauth is a reverse proxy that allowes IP based access to a docker registry.

# Setup

Create a `/etc/harbor-unauth.json` like this:

```json
{
  "url": "https://127.0.0.1:8080/",
  "proxy_count": 1,
  "allowlist": [
    "127.0.0.1/32",
    "172.17.0.1/32",
    "::1/128"
  ],
  "auth": {
    "user": "<robot user>",
    "password": "<robot token>"
  },
  "projects": [
    "<list of harbor projects...>"
  ]
}
```

Per default the the configuration is stored in `/etc/harbor-unauth.json`,
overwrite using the `UNAUTH_CONFIG` environment variable.

Note: The `HOST` header will not be rewritten to prevent issues with the dynamic
URLS generated by a Docker registry. This means that between unauth and Harbor a
loadbalancer using the Hostname will fail.

# Development

Setup a `config.json` and run `UNAUTH_CONFIG=./config.json go run .` Make sure
your configuration name is included in `.gitignore`.
