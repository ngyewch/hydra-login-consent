![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/ngyewch/hydra-login-consent/build.yml)
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/ngyewch/hydra-login-consent)
[![Go Reference](https://pkg.go.dev/badge/github.com/ngyewch/go-pqssh.svg)](https://pkg.go.dev/github.com/ngyewch/hydra-login-consent)

# hydra-login-consent

Golang http middleware for implementing the User Login and Consent flow of Ory OAuth2
service ([Hydra](https://github.com/ory/hydra)).

## Example implementation

1. Start a local instance of Hydra (OIDC OP).
   ```
   ./start-hydra.sh
   ```
2. Start a test app (OIDC RP).
   ```
   ./start-test-client.sh
   ```
3. [OPTIONAL] Edit the config file (`config.toml`) for the example implementation.
4. Start the example implementation.
   ```
   task run
   ```
5. Open the browser to `http://127.0.0.1:8080`
    * Test user credentials can be found in `config.toml`.
