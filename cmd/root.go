package cmd

import (
	"fmt"

  "github.com/dgutierrez1287/vault-util/logger"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var caCertFile string
var caKeyFile string
var skipTlsVerify bool
var vaultUrl string
var token string
var debug bool

var RootCmd = &cobra.Command{
  Use: "vault-util",
  Short: "A vault util program",
  Long: "A vault util program ",
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    setupLogging()
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("vault-util, Use --help for help")
  },
}

func Execute() error {
  return RootCmd.Execute()
}

func init() {
  
  // ca cert file; only needed if the cert for vault is self signed
  RootCmd.PersistentFlags().StringVarP(&caCertFile, "ca-cert-file", "", "N/A", "(Optional) The cacert file for access to the vault server")

  // ca cert key file: only needed if the cert for vault is self signed
  RootCmd.PersistentFlags().StringVarP(&caKeyFile, "ca-key-file", "", "N/A", "(Optional) the cacert key file for access to the vault server")
  
  // skip tls verification
  RootCmd.PersistentFlags().BoolVarP(&skipTlsVerify, "skip-tls-verify", "", false, "(Optional) to skip tls verification")

  // log debug settings
  RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

  // vault url
  RootCmd.PersistentFlags().StringVarP(&vaultUrl, "vault-url", "", "", "The url for the vault server")
  RootCmd.MarkFlagRequired("vault-url")

  // root token for vault 
  RootCmd.PersistentFlags().StringVarP(&token, "token", "", "", "The root token for access to vault")
  RootCmd.MarkFlagRequired("token")

}

func setupLogging() {
  // set up logger and set logging based on options
  if debug {
    logger.LogLevel = "DEBUG"
    fmt.Println("Debugging Enabled")
  } else {
    logger.LogLevel = "INFO"
  }

  fmt.Printf("loglevel: %s", logger.LogLevel)

  logger.Logger = hclog.New(&hclog.LoggerOptions{
    Name: "vault-util",
    Level: hclog.LevelFromString(logger.LogLevel),
  })
}
