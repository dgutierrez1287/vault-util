package app

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/vault-util/logger"
)

var userHomeDir = os.UserHomeDir
var statFunc = os.Stat
var marshalFunc = json.Marshal

/*
Settings - vault util settings
*/
type Settings struct {
  Vaults map[string]VaultInstance         `json:"vaults"`
}

/*
Adds a vault to the settings
*/
func (settings *Settings) AddVault(vaultName string, vaultData VaultInstance) {
  
  logger.LogDebug("checking if vault exists", "vault", vaultName)
  _, ok := settings.Vaults[vaultName]

  if ok {
    logger.LogDebug("vault already exists, updating settings")
    settings.Vaults[vaultName] = vaultData
  } else {
    logger.LogDebug("Adding vault")
    settings.Vaults[vaultName] = vaultData
  }
}

/*
deletes a vault from the settings
*/
func (settings *Settings) DeleteVault(vaultName string) {

  logger.LogDebug("checking if vault exists", "vault", vaultName)
  _, ok := settings.Vaults[vaultName]

  if ok {
    delete(settings.Vaults, vaultName)
    logger.LogDebug("vault settings deleted")
  } else {
    logger.LogDebug("vault is not present in settings")
  }
}

/*
Get the path for the app settings files
this will be at $HOME/$username/.vault-util-settings.json
and the equilvalent on windows
*/
func ConfigFilePath() (string, error) {
  userDir, err := userHomeDir()
  if err != nil {
    logger.LogError("Error getting user home dir")
    return "", err
  }

  return filepath.Join(userDir, ".vault-util-settings.json"), nil
}

/*
This will check if the settings file exists
*/
func SettingsFileExists(settingsFilePath string) (bool, error) {
  var fileExists bool

  if _, err := statFunc(settingsFilePath); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.LogDebug("Settings file does not exist")
      fileExists = false
    } else {
      logger.LogError("Error checking if file exists")
      return false, err
    }
  } else {
    logger.LogDebug("Settings file exists")
    fileExists = true
  }
  
  return fileExists, nil 
}

/*
Reads the settings file and retruns the settings
object
*/
func ReadSettingsFile(settingsFilePath string) (Settings, error) {
  
  file, err := os.Open(settingsFilePath) 
    if err != nil {
      logger.LogError("Error opening settings file")
      return Settings{}, err
    }
    defer file.Close()

    bytes, err := io.ReadAll(file)
    if err != nil {
      logger.LogError("Error reading settings file")
      return Settings{}, err
    }

    var settings Settings
    err = json.Unmarshal(bytes, &settings)
    if err != nil {
      logger.LogError("Error unmarshaling json to setting struct")
      return Settings{}, err
    }
    
    logger.LogDebug("Settings file read successfully")
    return settings, nil
}

/*
Writes the settings file
*/
func WriteSettingsFile(settingsFilePath string, settings Settings) error {

  jsonData, err := marshalFunc(&settings)
  if err != nil {
    logger.LogError("Error marshaling settings to json")
    return err
  }

  file, err := os.Create(settingsFilePath)
  if err != nil {
    logger.LogError("Error creating the settings file")
    return err
  }

  defer file.Close()

  _, err = file.Write(jsonData)
  if err != nil {
    logger.LogError("Error writing settings content to file")
    return err
  }
  return nil
}
