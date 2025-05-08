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

//detail flag
var outputMountDetail bool

var listMountsCmd = &cobra.Command{
  Use: "list-mounts",
  Short: "Lists secret mounts",
  Long: "Lists secrets mounts",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput app.MountListOutput
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

    logger.LogInfo("Getting secret mounts")
    mounts, err := app.GetSecretMounts(vaultClient)
    if err != nil {
      logger.LogErrorExit("Error getting the secret mounts", 250, err)
    }

    logger.LogDebug("Outputing results")
    if machineOutput {
      machineReadableOutput.ExitCode = 0
      if outputMountDetail {
        mountMap := app.MountstoMap(mounts)
        machineReadableOutput.MountsWithData = mountMap
      } else {
        mountNames := app.GetMountNames(mounts)
        machineReadableOutput.MountNames = mountNames
      }
      output, eCode := machineReadableOutput.GetOutputJson()
      fmt.Println(output)
      os.Exit(eCode)
    } else {
      if outputMountDetail {
        app.ListMountsWithDetailConsoleOutput(mounts)
        os.Exit(0)
      } else {
        app.ListMountNamesConsoleOutput(mounts)
        os.Exit(0)
      }
    }

  },
}

func init() {
  //command specific flags
  listMountsCmd.PersistentFlags().BoolVarP(&outputMountDetail, "detail", 
    "", false, "Output detail for all the secrets mounts")

  //Add Command
  RootCmd.AddCommand(listMountsCmd)
}
