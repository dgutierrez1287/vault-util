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

var listSecretsCmd = &cobra.Command{
  Use: "list-secrets",
  Short: "Lists secrets for a mount",
  Long: "Lists secrets for mount",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.SecretListOutput
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

    logger.LogInfo("Getting secret mount")
    secretMount, err := app.NewSecretMount(mountName, "", "", "", vaultClient)
    if err != nil {
      logger.LogErrorExit("Error getting secret mount details", 250, err)
    }

    logger.LogInfo("Getting secrets")
    secrets, err := secretMount.ListSecrets(vaultClient)
    if err != nil {
      logger.LogErrorExit("Error getting secrets for mount", 250, err)
    }

    logger.LogDebug("Outputing results")
    if machineOutput {
      machineReadableOutput.ExitCode = 0
      machineReadableOutput.Secrets = secrets
      output, eCode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(eCode)
    } 

    app.ListSecretsConsoleOutput(secrets, secretMount.Mount)
    os.Exit(0)
  },
}

func init() {
  // Required command cli options
  listSecretsCmd.MarkFlagRequired("secret-mount")

  // Command specific cli options

  // Add command 
  RootCmd.AddCommand(listSecretsCmd)
}
