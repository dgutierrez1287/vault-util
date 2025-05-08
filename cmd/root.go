package cmd

import (
	"fmt"

  "github.com/dgutierrez1287/vault-util/logger"
	"github.com/spf13/cobra"
)

// output flags
var debug bool
var logColorize bool
var machineOutput bool

//common command line options
// mount options
var mountName string

// Secret options
var secretKey string

// vault connection flags
var caCertFile string
var caKeyFile string
var skipTlsVerify bool
var vaultUrl string
var token string
var vaultName string


var RootCmd = &cobra.Command{
  Use: "vault-util",
  Short: "A vault utility cli",
  Long: "A vault utility to add functionality and add ease of use",
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    logger.InitLogging(debug, logColorize, machineOutput)
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("vault-util, Use --help for help")
  },
}

func Execute() error {
  return RootCmd.Execute()
}

func init() {

  /*
  Logging and output control options
  */
  // log debug settings
  RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

  // output colorization
  RootCmd.PersistentFlags().BoolVarP(&logColorize, "colorize", "", true, "Enable output colorization")

  /*
  machine output
  This will supress all other output and only output json at the end
  that is machine readable
  */
  RootCmd.PersistentFlags().BoolVarP(&machineOutput, "machine-output", "m", false, "Enables machine output for this to be run by another script")

  /*
  options that will be used for multiple command 
  line endpoints
  */
  // secret mount name
  RootCmd.PersistentFlags().StringVarP(&mountName, "secret-mount", "", "", "Secret mount")

  // secret key
  RootCmd.PersistentFlags().StringVarP(&secretKey, "secret-key", "", "", "Secret key")

  /*
  Vault connection options
  */
  // vault name: needed to reference a vault in the settings
  RootCmd.PersistentFlags().StringVarP(&vaultName, "vault-name", "", "", "(Optional) The name of the vault to use in the settings file")

  // ca cert file; only needed if the cert for vault is self signed
  RootCmd.PersistentFlags().StringVarP(&caCertFile, "ca-cert-file", "", "", "(Optional) The cacert file for access to the vault server")

  // ca cert key file: only needed if the cert for vault is self signed
  RootCmd.PersistentFlags().StringVarP(&caKeyFile, "ca-key-file", "", "", "(Optional) the cacert key file for access to the vault server")
  
  // skip tls verification
  RootCmd.PersistentFlags().BoolVarP(&skipTlsVerify, "skip-tls-verify", "", false, "(Optional) to skip tls verification")

  // vault url
  RootCmd.PersistentFlags().StringVarP(&vaultUrl, "vault-url", "", "", "The url for the vault server")

  // root token for vault 
  RootCmd.PersistentFlags().StringVarP(&token, "token", "", "", "The root token for access to vault")
}


