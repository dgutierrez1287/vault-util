package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dgutierrez1287/vault-util/logger"
  "github.com/dgutierrez1287/vault-util/app"
	"github.com/dgutierrez1287/vault-util/util"

	"github.com/spf13/cobra"
)

var secretsFile string
var kvVersion string
var bulkLoadCmd = &cobra.Command {
  Use: "bulk-load",
  Short: "bulk creates/updates secrets from a json file to vault",
  Long: "bulk creates/updates secrets from a json file to vault",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.BulkActionOutput
    var secretsAdded []string
    var secretErrors []app.SecretActionError
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

    logger.LogInfo("Reading secrets from json file", "file", secretsFile)
    secrets, err := app.ReadSecretsFromJson(secretsFile, vaultClient, ctx)

    if err != nil {
      logger.LogErrorExit("Error reading secrets from json file", 250, err)
    }

    logger.LogInfo("Creating or updating secrets")
    for name, secret := range secrets.Secrets {
      logger.LogDebug("writing secret secret", "name", name)
      err = secret.WriteSecret(vaultClient)

      if err != nil {
        logger.LogError("Error writing secret", "error", err)
        secretErrors = append(secretErrors, app.SecretActionError{
          VaultKey: secret.VaultKey,
          Error: err,
        })
      } else {
        logger.LogInfo("Secret created/updated")
        secretsAdded = append(secretsAdded, name)
      }
    }

    logger.LogDebug("Outputing results")
    if machineOutput {
      machineReadableOutput.ExitCode = 0
      machineReadableOutput.SecretsAdded = secretsAdded
      machineReadableOutput.Errors = secretErrors

      output, eCode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(eCode)
    } else {
      app.BulkActionConsoleOutput(secretsAdded, secretErrors, "added")
      os.Exit(0)
    }
  },
}

func init() {
  // add command
  RootCmd.AddCommand(bulkLoadCmd)

  // secrets file
  bulkLoadCmd.PersistentFlags().StringVarP(&secretsFile, "secrets-file", "", "", "The json file that contains the secrets to be loaded/updated")
  bulkLoadCmd.MarkFlagRequired("secrets-file")
}



