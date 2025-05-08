package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/vault-util/logger"
	"github.com/dgutierrez1287/vault-util/app"
	"github.com/dgutierrez1287/vault-util/util"
	"github.com/spf13/cobra"
)

var addVaultCmd = &cobra.Command{
  Use: "add-vault",
  Short: "Adds a vault to the config",
  Long: "Adds a vault to the config so connection settings are saved",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.AddRemoveOutput
    
    if !machineOutput {
      fmt.Println(util.TitleString)
    }

    logger.LogInfo("Getting the settings file path")
    settingsFilePath, err := app.ConfigFilePath()
    if err != nil {
      logger.LogErrorExit("Error getting settings file path", 100, err)
    }

    logger.LogInfo("Checking if the settings file exists")
    exists, err := app.SettingsFileExists(settingsFilePath)
    if err != nil {
      logger.LogErrorExit("Error checking for the settings file", 100, err)
    }

    var appSettings app.Settings
    if exists {

      logger.LogInfo("Settings file exists, getting current settings")
      appSettings, err = app.ReadSettingsFile(settingsFilePath)
      if err != nil {
        logger.LogErrorExit("Error reading the settings file", 100, err)
      }
    } else {

      logger.LogInfo("Settings file doesn't exist, creating blank settings")
      appSettings = app.Settings{Vaults: make(map[string]app.VaultInstance)}
    }

    logger.LogDebug("Creating vault instance")
    vaultInst, err := app.NewVault(vaultUrl, token, skipTlsVerify, caCertFile, caKeyFile)

    if err != nil {
      logger.LogErrorExit("Error creating new vault instance", 100, err)
    }

    logger.LogInfo("Adding vault to config", "name", vaultName)
    appSettings.AddVault(vaultName, *vaultInst)

    logger.LogInfo("Updating config with new vault")
    err = app.WriteSettingsFile(settingsFilePath, appSettings)

    if err != nil {
      logger.LogErrorExit("Error updating the config", 100, err)
    }

    if !machineOutput {
      logger.LogInfo("Vault instance successfully added", "name", vaultName)
      os.Exit(0)
    } else {
      machineReadableOutput.ExitCode = 0
      machineReadableOutput.Message = fmt.Sprintf("%s vault successfully created", vaultName)
      output, eCode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(eCode)
    }
  },
}

func init() {
  // Required common cli options
  addVaultCmd.MarkFlagRequired("vault-name")
  
  // Command specific cli options

  // Add command 
  RootCmd.AddCommand(addVaultCmd)
}
