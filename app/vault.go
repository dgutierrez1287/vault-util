package app

import (
	"errors"
	"os"

	"github.com/dgutierrez1287/vault-util/logger"
)

/*
VaultInstance - settings for a vault instance
in the settings file
*/
type VaultInstance struct {
  Url string                                        `json:"url"`
  Token string                                      `json:"token"`
  SkipTLSVerify bool                                `json:"skipTlsVerify"`
  CACert string                                     `json:"caCert,omitempty"`
  CACertKey string                                  `json:"caCertKey,omitempty"`
}

func NewVault(url string, token string, skipTlsVerify bool, caCertFilePath string,
caCertKeyFilePath string) (*VaultInstance, error) {

  if url == "" {
    logger.LogError("Error vault url cannot be empty")
    return &VaultInstance{}, errors.New("vault url is empty")
  }

  if token == "" {
    logger.LogError("Error vault token cannot be empyt")
    return &VaultInstance{}, errors.New("vault token is empty")
  }
  
  logger.LogDebug("Creating a new vault instance")
  vInst := VaultInstance{
    Url: url,
    Token: token,
    SkipTLSVerify: skipTlsVerify,
  }

  logger.LogDebug("Getting the ca cert data if needed")
  err := vInst.getCertDataFromFile(caCertFilePath, caCertKeyFilePath)
  
  if err != nil {
    logger.LogError("Error getting the cacert data")
    return &VaultInstance{}, err
  }

  return &vInst, nil 
}

/*
This will read in the ca cert and ca cert key into the 
settings from the files
*/
func (v *VaultInstance) getCertDataFromFile(caCertFile string, caCertKeyFile string) error {

  if caCertFile == "" || caCertKeyFile == "" {
    logger.LogDebug("caCert file or key option is empty, skipping")
    return nil
  }
  
  logger.LogDebug("Reading the cacert file")
  caCertData, err := os.ReadFile(caCertFile) 

  if err != nil {
    logger.LogError("Error reading caCert file")
    return err
  }

  v.CACert = string(caCertData)

  logger.LogDebug("Reading the cacert key file")
  caCertKeyData, err := os.ReadFile(caCertKeyFile)

  if err != nil {
    logger.LogError("Error reading caCert key file")
    return err
  }

  v.CACertKey = string(caCertKeyData)
  return nil
}

/*
This will get a vault instance from the settings
file given a vault name
*/
func GetVaultConfigFromSettings(vaultName string, 
  settingsFilePath string) (*VaultInstance, error) {
  
  logger.LogDebug("Checking if the settings file exists")
  exists, err := SettingsFileExists(settingsFilePath)

  if err != nil {
    logger.LogError("Error checking for the settings file")
    return &VaultInstance{}, err
  }

  if !exists {
    logger.LogError("Error settings file doesn't exist")
    return &VaultInstance{}, errors.New("settings file doesn't exist")
  }

  appSettings, err := ReadSettingsFile(settingsFilePath)
  
  if err != nil {
    logger.LogError("Error reading the settings file")
    return &VaultInstance{}, err
  }

  vaultInst, ok := appSettings.Vaults[vaultName]

  if !ok {
    logger.LogError("Error vault does not exist in settings file")
    return &VaultInstance{}, errors.New("vault does not exist in settings")
  }

  return &vaultInst, nil
}


