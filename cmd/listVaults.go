package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/vault-util/logger"
  "github.com/dgutierrez1287/vault-util/app"
	"github.com/dgutierrez1287/vault-util/util"
	"github.com/spf13/cobra"
)

var listVaultsCmd = &cobra.Command{
  Use: "list-vaults",
  Short: "Lists the names of the vaults in the config",
  Long: "Lists the names of the vaults in the config",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.VaultListOutput

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
    if !exists {
      logger.LogDebug("Settings file does not exist outputing")
      if machineOutput {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.Vaults = []string{}
        machineReadableOutput.Message = "no settings file, no vaults configured"

        output, eCode := machineReadableOutput.GetOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
      
      app.ListVaultsConsoleOutput(exists, []string{})
      os.Exit(0)
    }
      
    logger.LogInfo("Getting the list of vaults from settings")
    appSettings, err = app.ReadSettingsFile(settingsFilePath)
    if err != nil {
      logger.LogErrorExit("Error reading the settings file", 100, err)
    }

    vaultNames := []string{}

    for name := range appSettings.Vaults {
      vaultNames = append(vaultNames, name)
    }

    if machineOutput {
      machineReadableOutput.ExitCode = 0
      machineReadableOutput.Vaults = vaultNames

      output, eCode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(eCode)
    }

    app.ListVaultsConsoleOutput(true, vaultNames)
    os.Exit(0)
  },
}

func init() {
  // add command
  RootCmd.AddCommand(listVaultsCmd)
}
