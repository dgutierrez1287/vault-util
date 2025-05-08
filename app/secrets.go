package app

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/dgutierrez1287/vault-util/logger"
)

/*
VaultSecrets - map of secrets to read a bulk
set of secrets from a file
*/
type VaultSecrets struct {
  Secrets map[string]VaultSecret           `json:"secrets"`
}

/*
Reads secrets from json file this requires the 
format

"secrets": {
  "<secretName>": {
    "key": "<secretKeyName>",
    "data": {
      "<key>": "<value>",
      .
      .
      .
    }
  }
}
*/
func ReadSecretsFromJson(secretsFilePath string, 
  client *VaultClient, ctx context.Context) (VaultSecrets, error) {
  var secrets VaultSecrets

  file, err := os.Open(secretsFilePath)
  if err != nil {
    logger.LogError("Error opening secrets file")
    return secrets, err
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.LogError("Error reading secrets file")
    return secrets, err
  }

  err = json.Unmarshal(bytes, &secrets)
  if err != nil {
    logger.LogError("Error unmarshaling json to struct")
    return secrets, err
  }

  logger.LogDebug("Getting secret details for secrets in the list")
  for name, secret := range secrets.Secrets {
    logger.LogDebug("Getting details for secret", "name", name)
    secret.getSecretDetails(client)
    secrets.Secrets[name] = secret
  }

  return secrets, nil
}


