package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dgutierrez1287/vault-util/app"
	"github.com/dgutierrez1287/vault-util/logger"
	"github.com/dgutierrez1287/vault-util/util"
	"github.com/spf13/cobra"
)

var getSecretCmd = &cobra.Command{
  Use: "get-secret",
  Short: "Gets the secret data for secret",
  Long: "Gets the secret data for secret",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.GetSecretOutput
    var vaultInstance *app.VaultInstance
    var err error

    if !machineOutput {
      fmt.Println(util.TitleString)
    }

    // Get vault configuration from settings file
    if vaultName != "" {
      logger.LogInfo("Vault name passed, getting connection details from settings file")

      logger.LogInfo("Getting the settings file path")
      settingsFilePath, err := app.ConfigFilePath()
      if err != nil {
        logger.LogErrorExit("Error getting settings file path", 200, err)
      }

      vaultInstance, err = app.GetVaultConfigFromSettings(vaultName, settingsFilePath)
      if err != nil {
        logger.LogErrorExit("Error getting the vault config from settings", 200, err)
      }
    } else {
      logger.LogInfo("No Vault name is passed getting connection details from command line")

      vaultInstance, err = app.NewVault(vaultUrl, token, skipTlsVerify, 
        caCertFile, caKeyFile)
      if err != nil {
        logger.LogErrorExit("Error creating the vault instance", 150, err)
      }
    }

    ctx := context.Background()

    logger.LogInfo("Getting vault client")
    vaultClient, err := app.NewClient(*vaultInstance, &ctx)
    if err != nil {
      logger.LogErrorExit("Error getting vault client", 250, err)
    }

    data := make(map[string]interface{})
    secret, err := app.NewSecret(secretKey, "", "", data, *vaultClient)
    if err != nil {
      logger.LogErrorExit("Error getting vault secret", 250, err)
    }

    logger.LogInfo("Checking if the secret exists")
    secretExists, err := secret.SecretExists(vaultClient)
    if err != nil {
      logger.LogErrorExit("Error checking if secret exists", 250, err)
    }

    if !secretExists {
      logger.LogInfo("Secret does not exist")
      if machineOutput {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.SecretExists = secretExists
        output, eCode := machineReadableOutput.GetOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }  else {
        app.GetSecretConsoleOutput(secret, secretExists)
        os.Exit(0)
      }
    }

    logger.LogInfo("Reading the secret")
    err = secret.ReadSecret(vaultClient)
    if err != nil {
      logger.LogErrorExit("Error reading vault secret", 250, err)
    }

    logger.LogDebug("Outputing results")
    if machineOutput {
      machineReadableOutput.ExitCode = 0
      machineReadableOutput.SecretExists = secretExists
      machineReadableOutput.VaultKey = secret.NormalizedSecretPath
      machineReadableOutput.Data = secret.SecretData
      output, ecode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(ecode)
    }

    app.GetSecretConsoleOutput(secret, secretExists)
    os.Exit(0)
  },
}

func init() {
  //Required command cli options
  getSecretCmd.MarkFlagRequired("secret-key")

  // command specific cli options

  //Add command
  RootCmd.AddCommand(getSecretCmd)
}
