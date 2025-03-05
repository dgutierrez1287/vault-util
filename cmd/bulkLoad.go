package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

  "github.com/dgutierrez1287/vault-util/util"
  "github.com/dgutierrez1287/vault-util/logger"

	"github.com/spf13/cobra"
)

var secretsFile string
var kvVersion string
var bulkLoadCmd = &cobra.Command {
  Use: "bulk-load",
  Short: "bulk Loads secrets from a json file to vault",
  Long: "bulk Loads secrets from a json file to vault",
  Run: func(cmd *cobra.Command, args []string) {

    fmt.Print(util.TitleString)

    // get secrets from the json file 
    logger.Logger.Info("Reading secrets from json file")
    secrets := readSecretsFromJson(secretsFile)
    
    logger.Logger.Debug("Command line args")
    logger.Logger.Debug("secrets file", "path", secretsFile)
    logger.Logger.Debug("cacert file path", "path" ,caCertFile)
    logger.Logger.Debug("cacert key file path", "path" ,caKeyFile)
    logger.Logger.Debug("vault url", "url" ,vaultUrl)
    logger.Logger.Debug("token", "token" ,token)
    logger.Logger.Debug("skip tls verify", "setting", skipTlsVerify)
    logger.Logger.Debug("KV Version", "version", kvVersion)

    settings := util.VaultSettings{
      VaultUrl: vaultUrl,
      Token: token,
      CaCertFilePath: caCertFile,
      CaCertKeyFilePath: caKeyFile,
      SkipTlsVerify: skipTlsVerify,
    }
    
    // Get vault client
    logger.Logger.Info("Getting vault client")
    client := util.GetVaultClient(&settings)
    
    // load the secrets into vault 
    logger.Logger.Info("Loading secrets into vault")
    util.LoadSecrets(client, kvVersion ,&secrets)

  },
}

func init() {
  RootCmd.AddCommand(bulkLoadCmd)

  // secrets file
  bulkLoadCmd.PersistentFlags().StringVarP(&secretsFile, "secrets-file", "", "", "The json file that contains the secrets to be loaded/updated")
  bulkLoadCmd.MarkFlagRequired("secrets-file")

  // secret kv version
  bulkLoadCmd.PersistentFlags().StringVarP(&kvVersion, "kv-version", "", "v2", "The KV version for vault, defaults to v2 (v1 or v2 is allowed)")
}

// Read secrets from the secrets json file, this requires a map in the 
// json file as follows
// "secrets": {}
func readSecretsFromJson(secrtsFile string) util.VaultSecrets {
  file, err := os.Open(secretsFile)
  if err != nil {
    logger.Logger.Error("Error opening secrets file", "error", err)
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.Logger.Error("Error reading secrets file", "error", err)
  }

  var secrets util.VaultSecrets
  err = json.Unmarshal(bytes, &secrets)
  if err != nil {
    logger.Logger.Error("Error unmarshaling json to struct", "error", err)
  }

  return secrets
}


