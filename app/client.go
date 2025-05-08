package app

import (
	"context"

	"github.com/dgutierrez1287/vault-util/logger"
	vaultGo "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

/*
Interface for the vault client
*/
type VaultClientInterface interface {
  //Secrets
  WriteKvSecret(secret VaultSecret) error
  ReadKvSecret(secret VaultSecret) (map[string]interface{}, error)
  ListKvSecret(mount string, kvVersion string) ([]string, error)
  //System
  getSecretMountsData() (map[string]interface{}, error)
}

/*
Vault client
*/
type VaultClient struct {
  secrets   *vaultGo.Secrets
  system    *vaultGo.System
  ctx       *context.Context
}

/*
Returns a VaultClient with any custom tls configuration enabled
and token set
*/
func NewClient(v VaultInstance, ctx *context.Context) (*VaultClient, error) {
  
  logger.LogDebug("Checking for custom tls settings")
  tlsConfigEnabled, tlsConfig := getVaultTlsConfig(v)

  var err error
  var client *vaultGo.Client

  if tlsConfigEnabled {
    logger.LogDebug("Creating vault client with custom tls config")
    client, err = vaultGo.New(
      vaultGo.WithAddress(v.Url),
      vaultGo.WithTLS(tlsConfig),
    )
  } else {
    logger.LogDebug("Creating vault client without custom tls config")
    client, err = vaultGo.New(
      vaultGo.WithAddress(v.Url),
    )
  }

  if err != nil {
    logger.LogError("Error getting the vault client")
    return nil, err
  }

  logger.LogDebug("Setting client token")
  client.SetToken(v.Token)

  return &VaultClient{
    secrets: &client.Secrets,
    system: &client.System,
    ctx: ctx,
  }, nil
}

/*
wrapper for kv write secret
*/
func (c *VaultClient) WriteKvSecret(s VaultSecret) error {

  if s.KvVersion == "2" {
    logger.LogDebug("Writing kv v2 secret")

    writeReq := schema.KvV2WriteRequest {
      Data: s.SecretData,
    }
    _, err := c.secrets.KvV2Write(*c.ctx, s.NormalizedSecretPath, 
      writeReq, vaultGo.WithMountPath(s.MountName))
    return err
  }

  logger.LogDebug("Writing kv v1 secret")
  _, err := c.secrets.KvV1Write(*c.ctx, s.NormalizedSecretPath, 
    s.SecretData, vaultGo.WithMountPath(s.MountName))

  return err
}

/*
wrapper for kv list secret
*/
func (c *VaultClient) ListKvSecrets(mount string, path string, kvVersion string) ([]string, error) {

  if kvVersion == "2" {
    logger.LogDebug("Getting a list of kv v2 secrets for", "mount", mount)
    resp, err := c.secrets.KvV2List(*c.ctx, path, vaultGo.WithMountPath(mount))
    return resp.Data.Keys, err
  }

  logger.LogDebug("Getting a list of kv v1 secrets for", "mount", mount)
  resp, err := c.secrets.KvV1List(*c.ctx, path, vaultGo.WithMountPath(mount))
  return resp.Data.Keys, err
}


/*
wrapper for kv read secret
*/
func (c *VaultClient) ReadKvSecret(s VaultSecret) (map[string]interface{}, error) {
  data := make(map[string]interface{})

  if s.KvVersion == "2" {
    logger.LogDebug("Reading kv v2 secret")

    resp, err := c.secrets.KvV2Read(*c.ctx, s.NormalizedSecretPath, 
      vaultGo.WithMountPath(s.MountName))
    if err != nil {
      logger.LogError("Error reading the v2 secret")
      return data, err
    }
    return resp.Data.Data, nil
  }

  logger.LogDebug("Reading kv v1 secret")
  resp, err := c.secrets.KvV1Read(*c.ctx, s.NormalizedSecretPath, 
    vaultGo.WithMountPath(s.MountName))
  if err != nil {
    logger.LogError("Error reading the v1 secret")
    return data, err
  }
  return resp.Data, nil
}

/*
wrapper for MountsListSecretsENgines
*/
func (c *VaultClient) GetSecretMountsData() (map[string]interface{}, 
  error) {

    mounts, err := c.system.MountsListSecretsEngines(*c.ctx)
    if err != nil {
      logger.LogError("Error getting a list of secrets engines")
      return nil, err
    }

    return mounts.Data, nil
}
 
/*
Checks if any custom tls configuration is needed and returns if 
that custom configuration is enabled and what that configuration is
*/
func getVaultTlsConfig(v VaultInstance) (bool, vaultGo.TLSConfiguration) {
  if v.SkipTLSVerify {
    logger.LogDebug("Skipping all tls verification")
    vaultTls := vaultGo.TLSConfiguration {
      InsecureSkipVerify: true,
    }

    logger.LogDebug("tls configuration", "config", vaultTls)
    return true, vaultTls

  } else if v.CACert != "" {
    logger.LogDebug("Using custom ca cert and key for tls verification")
    vaultTls := vaultGo.TLSConfiguration {
      ClientCertificate: vaultGo.ClientCertificateEntry {
        FromBytes: []byte(v.CACert),
      },
      ClientCertificateKey: vaultGo.ClientCertificateKeyEntry {
        FromBytes: []byte(v.CACertKey),
      },
    }

    logger.LogDebug("tls configuration", "config", vaultTls)
    return true, vaultTls

  } else {
    logger.LogDebug("No tls information provided, assuming known cert")
    return false, vaultGo.TLSConfiguration{}
  }
}



