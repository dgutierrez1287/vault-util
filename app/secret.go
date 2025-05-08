package app

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/dgutierrez1287/vault-util/logger"
)

/*
Secret - The components of a vault kv secret
*/
type VaultSecret struct {
  VaultKey string                     `json:"key"`
  SecretData map[string]interface{}   `json:"data"`
  SecretType string                   `json:"secretType"`
  KvVersion string                    `json:"kvVersion,omitempty"`
  NormalizedSecretPath string         `json:"normalizedSecretPath"`
  MountName string                    `json:"mountName"`
}

/*
SecretActionError - a type to hold the secret
and also the error that happened when the 
action was run, This is so that errors can
be stored for output
*/
type SecretActionError struct {
  VaultKey string           `json:"secretKey"`
  Error error               `json:"error"`
}

/*
Create a new secret 
*/
func NewSecret(key string, secretType string, kvVersion string,
  data map[string]interface{}, client VaultClient) (VaultSecret, error) {

  secret := VaultSecret{
    VaultKey: key,
    SecretType: secretType,
    KvVersion: kvVersion,
    SecretData: data,
  }

  err := secret.getSecretDetails(&client)

  if err != nil {
    logger.LogError("Error creating new secret")
    return secret, err
  }
  
  return secret, nil
}

/*
This will add additional details to the secret to include 
secret mount type, kv version and the normalized secret path 
these additional details shouldn't needed to be passed by the user
but will be useful when interacting with the secret
*/
func (s *VaultSecret) getSecretDetails(client *VaultClient) error{

  if s.SecretType != ""  {
    if s.SecretType == "kv" && s.KvVersion != "" {
      logger.LogDebug("Secret type was already provided", "name", s.VaultKey,
      "type", s.SecretType, "kvVersion", s.KvVersion)
      return nil
    }
    logger.LogDebug("Secret type was already provided", "name", s.VaultKey,
      "type", s.SecretType)
      return nil
  }
  
  logger.LogDebug("Processing the vault key to get the mountpoint")
  parts := strings.SplitN(s.VaultKey, "/", 2)
  
  if len(parts) < 2 {
    logger.LogError("Error splitting path, please check vault key for", 
      "key", s.VaultKey) 
    return errors.New("error splitting vault key path")
  }

  secretsMount := fmt.Sprintf("%s/", parts[0])
  logger.LogDebug("Setting mount name to secret for later use", "mount", secretsMount)
  s.MountName = secretsMount

  logger.LogDebug("Getting the mount data to figure out type")
  mountType, kvVersion, err := GetMountType(client, secretsMount)

  if err != nil {
    logger.LogError("Error getting secret mount data", "mountName", secretsMount)
    return err
  }

  logger.LogDebug("Setting secret type", "type", mountType)
  s.SecretType = mountType

  if mountType == "kv" {
    logger.LogDebug("Secret type is kv, setting kv version", "kvVersion", kvVersion)
    s.KvVersion = kvVersion
  }
  
  logger.LogDebug("Getting the normalized secret path")
  err = s.getNormalizedSecretPath()

  if err != nil {
    logger.LogError("Error getting the normalized secret path")
    return err
  }

  logger.LogDebug("Additional secret details", "type", s.SecretType, 
    "kvVersion", s.KvVersion, "secretPath", s.NormalizedSecretPath)
  return nil
}

/*
This will get the normalized secret path based on the 
secret type
*/
func (s *VaultSecret) getNormalizedSecretPath() error {

  logger.LogDebug("Splitting vault path")
  parts := strings.SplitN(s.VaultKey, "/", 2)
  
  if len(parts) < 2 {
    logger.LogError("Error splitting path, please check vault key for", "key", s.VaultKey) 
    return errors.New("error splitting vault key path")
  }

  if s.SecretType != "kv" {
    logger.LogDebug("Secret type not KV, secret path is the default")
    s.NormalizedSecretPath = s.VaultKey
    return nil
  }

  if s.KvVersion == "2" {
    logger.LogDebug("Kv Version is 2 setting secret path")
    
    // Check if the vault key provided contains the data part for 
    // kv version 2
    dataCheckParts := strings.SplitN(parts[1], "/", 2)
    
    if dataCheckParts[0] != "data" {
      s.NormalizedSecretPath = fmt.Sprintf("%s/data/%s", parts[0], parts[1])
      logger.LogDebug("normalized secret path", "path", s.NormalizedSecretPath)
    } else {
      s.NormalizedSecretPath = s.VaultKey
      logger.LogDebug("normalized secret path", "path", s.NormalizedSecretPath)
    }

  } else {
    logger.LogDebug("Kv version is 1 setting secret path")
    s.NormalizedSecretPath = s.VaultKey
    logger.LogDebug("normalized secret path", "path", s.NormalizedSecretPath)
  }
  return nil
}

/*
write a secret
*/
func (s VaultSecret) WriteSecret(client *VaultClient) error {
  if s.SecretType == "kv" {
    logger.LogDebug("Secret is kv type")
    
    logger.LogDebug("writing secret", "path", s.NormalizedSecretPath, "data", s.SecretData)
    err := client.WriteKvSecret(s)

    if err != nil {
      logger.LogError("Error writing the kv secret")
      return err
    }
  }
  return nil
}

/*
Read a secret, this will put the data back into the secret object
*/
func (s *VaultSecret) ReadSecret(client *VaultClient) error {
  if s.SecretType == "kv" {
    logger.LogDebug("Secret is kv type")

    logger.LogDebug("Reading secret", "path", s.NormalizedSecretPath)
    data, err := client.ReadKvSecret(*s)

    if err != nil {
      logger.LogError("Error reading the kv secret")
      return err
    }

    s.SecretData = data
  }
  return nil
}

/*
Check if a secret exists
*/
func (s VaultSecret) SecretExists(client *VaultClient) (bool, error) {
  var secrets []string
  var err error

  secretName := path.Base(s.NormalizedSecretPath)
  secretDir := path.Dir(s.NormalizedSecretPath) + "/"

  logger.LogDebug("secret details", "name", secretName, "dir", secretDir)

  if s.SecretType == "kv" {
    logger.LogDebug("Secret is kv type")

    secrets, err = client.ListKvSecrets(s.MountName, secretDir, s.KvVersion)

    if err != nil {
      logger.LogError("Error getting list of kv secrets")
      return false, err
    }
  }

  logger.LogDebug("secrets", secrets)

  if slices.Contains(secrets, secretName) {
    logger.LogDebug("Secret found")
    return true, nil
  }
  logger.LogDebug("Secret not found")
  return false, nil
}


