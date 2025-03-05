package util

import (
	"context"

	"github.com/dgutierrez1287/vault-util/logger"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

func GetVaultClient(settings *VaultSettings) *vault.Client {
  logger.Logger.Info("Getting Vault Tls configuration (if needed)")
  getVaultTlsConfig(settings)

  var client *vault.Client
  var err error

  logger.Logger.Info("Creating vault client")
  if settings.UseTlsConfig {
    client, err = vault.New(
      vault.WithAddress(settings.VaultUrl),
      vault.WithTLS(settings.TlsConfig),
    )
  } else {
    client, err = vault.New(
      vault.WithAddress(settings.VaultUrl),
    )
  }

  if err != nil {
    logger.Logger.Error("Error getting vault client", "error", err)
  }

  client.SetToken(settings.Token)

  logger.Logger.Debug("client value", "client", client)

  return client
}
  
func LoadSecrets(client *vault.Client, kvVersion string ,secrets *VaultSecrets) {

  ctx := context.Background()

  for name, secret := range secrets.Secrets {
    logger.Logger.Info("Processing secret", "name", name)
    logger.Logger.Debug("Secret Name", "name", name)
    logger.Logger.Debug("Secret value", "secret", secret)
    logger.Logger.Debug("Secret vault path", "key", secret.VaultKey)

    secret.splitPath()

    var normalizedSecretPath string
    if kvVersion == "v2" {
      normalizedSecretPath = secret.SecretKey + "/data" + secret.SecretPath
    } else {
      normalizedSecretPath = secret.VaultKey
    }

    logger.Logger.Debug("The secret path for this kv is", "path", normalizedSecretPath, "version" ,kvVersion)

    writeReq := schema.KvV2WriteRequest{
      Data: secret.SecretData,
    }
    _, err := client.Secrets.KvV2Write(ctx, normalizedSecretPath, writeReq, vault.WithMountPath(secret.SecretKey))

    if err != nil {
      logger.Logger.Error("error writing secret", "error" ,err)
    } else {
      logger.Logger.Info("secret was created/updated")
    }
  }
}

func getVaultTlsConfig(settings *VaultSettings) {
  if settings.SkipTlsVerify {
    logger.Logger.Info("skipping all tls verification")
    vaultTls := vault.TLSConfiguration {
      InsecureSkipVerify: true,
    }

    logger.Logger.Debug("tls configuration", "config", vaultTls)

    settings.UseTlsConfig = true
    settings.TlsConfig = vaultTls

  } else if settings.CaCertFilePath != "N/A" {
    logger.Logger.Info("using custom ca cert and key for tls verification")
    vaultTls := vault.TLSConfiguration{
      ClientCertificate: vault.ClientCertificateEntry{
        FromFile: settings.CaCertFilePath,
      },
      ClientCertificateKey: vault.ClientCertificateKeyEntry{
        FromFile: settings.CaCertKeyFilePath,
      },
    }

    logger.Logger.Debug("tls configuration", "config", vaultTls)

    settings.UseTlsConfig = true
    settings.TlsConfig = vaultTls
  } else {
    
    logger.Logger.Info("No tls information provided, assuming know cert")
    settings.UseTlsConfig = false
  }
}
