package app

import (
	"errors"
	"os"
	"testing"

	"github.com/dgutierrez1287/vault-util/logger"
	"github.com/dgutierrez1287/vault-util/util"
	"github.com/stretchr/testify/assert"
)

// TestMain is executed before running any tests
func TestMain(m *testing.M) {
	// Initialize the logger before running any tests
	logger.InitLogging(false, true, false)
	os.Exit(m.Run())
}

/*
    Tests for AddVault
*/
func TestAddVaultNotExist(t *testing.T) {
  settings := Settings{Vaults: make(map[string]VaultInstance)}

  vaultName := "test-vault"

  vaultToAdd := VaultInstance{
    Url: "https://testvault.com",
    Token: "faketoken",
    SkipTLSVerify: true,
  }

  settings.AddVault(vaultName, vaultToAdd)

  vault, ok := settings.Vaults["test-vault"]
  assert.True(t, ok)
  assert.Equal(t, vault.Token, "faketoken")
  assert.Equal(t, vault.Url, "https://testvault.com")
}

func TestAddVaultUpdate(t *testing.T) {
  settings := Settings{Vaults: map[string]VaultInstance{}}

  vaultName := "test-vault"

  vaultData := VaultInstance{
    Url: "https://testvault.com",
    Token: "faketoken",
    SkipTLSVerify: true,
  }

  settings.Vaults[vaultName] = vaultData

  vaultData.Token = "newfaketoken"

  settings.AddVault(vaultName, vaultData)

  vault, ok := settings.Vaults["test-vault"]
  assert.True(t, ok)
  assert.Equal(t, vault.Token, "newfaketoken")
  assert.Equal(t, vault.Url, "https://testvault.com")
}

/*
    Tests for DeleteVault
*/
func TestDeleteVault(t *testing.T) {
  settings := Settings{make(map[string]VaultInstance)}

  vaultName := "test-vault"

  vaultData := VaultInstance{
    Url: "https://testvault.com",
    Token: "faketoken",
    SkipTLSVerify: true,
  }

  settings.Vaults[vaultName] = vaultData

  settings.DeleteVault(vaultName)

  _, ok := settings.Vaults["test-vault"]
  assert.False(t, ok)
}

func TestDeleteVaultNotExists(t *testing.T) {
  settings := Settings{make(map[string]VaultInstance)}

  vaultName := "test-vault"

  settings.DeleteVault(vaultName)

  _, ok := settings.Vaults["test-vault"]
  assert.False(t, ok)
}

/*
    Tests for ConfigFilePath
*/
func TestConfigFilePath(t *testing.T) {
  originalFunc := userHomeDir
  userHomeDir = func() (string, error) {
    return "/mock/home", nil
  }
  defer func() {userHomeDir = originalFunc}()

  path, err := ConfigFilePath()
  assert.NoError(t, err)
  assert.Equal(t, path, "/mock/home/.vault-util-settings.json")
}

func TestConfigFilePathError(t *testing.T) {
  originalFunc := userHomeDir
  userHomeDir = func() (string, error) {
    return "", errors.New("Fake Error")
  }
  defer func() {userHomeDir = originalFunc}()

  path, err := ConfigFilePath()
  assert.Error(t, err)
  assert.Equal(t, path, "")
}

/*
    Tests for SettingsFileExists
*/
func TestSettingsFileExists(t *testing.T) {
  err := util.MockHomeSetup()
  assert.NoError(t, err)

  _, err = os.Create(util.MockSettingsFile)
  assert.NoError(t, err)

  exists, err := SettingsFileExists(util.MockSettingsFile)
  assert.NoError(t, err)
  assert.True(t, exists)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestSettingsFileExistsNotExists(t *testing.T) {
  err := util.MockHomeSetup()
  assert.NoError(t, err)

  exists, err := SettingsFileExists(util.MockSettingsFile)
  assert.NoError(t, err)
  assert.False(t, exists)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestSettingsFileExistsError(t *testing.T) {

  originalStatFunc := statFunc
  defer func() { statFunc = originalStatFunc }()
  expectedErr := errors.New("stat error")
	statFunc = func(name string) (os.FileInfo, error) {
		return nil, expectedErr
	}

  exists, err := SettingsFileExists(util.MockSettingsFile)
  assert.Error(t, err)
  assert.False(t, exists)
}

/*
    Tests for ReadSettingsFile
*/
func TestWriteReadSettings(t *testing.T) {
  err := util.MockHomeSetup()
  assert.NoError(t, err)

  testVaultData := VaultInstance{
    Url: "https://testvault.com",
    Token: "faketoken",
  }

  vaults := make(map[string]VaultInstance)
  vaults["test-vault"] = testVaultData

  settings := Settings{
    Vaults: vaults,
  }

  err = WriteSettingsFile(util.MockSettingsFile, settings)
  assert.NoError(t, err)

  readSettings, err := ReadSettingsFile(util.MockSettingsFile)
  assert.NoError(t, err)

  vault, ok := readSettings.Vaults["test-vault"]
  assert.True(t, ok)
  assert.Equal(t, vault.Url, "https://testvault.com")
  assert.Equal(t, vault.Token, "faketoken")

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestReadSettingsError(t *testing.T) {
  _, err := ReadSettingsFile(util.MockSettingsFile)

  assert.Error(t, err)
}

func TestReadSettingsUnmarshalError(t *testing.T) {
  err := util.MockHomeSetup()
  assert.NoError(t, err)

  badJsonText := `
  {
    "vaults": "I should not be a string"
  }
  `

  file, err := os.Create(util.MockSettingsFile)
  assert.NoError(t, err)

  _, err = file.Write([]byte(badJsonText))
  assert.NoError(t, err)

  _, err = ReadSettingsFile(util.MockSettingsFile)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

/*
    Tests for WriteSettingsFile
*/
func TestWriteSettingsError(t *testing.T) {

  vaults := make(map[string]VaultInstance)
  settings := Settings{
    Vaults: vaults,
  }

  err := WriteSettingsFile(util.MockSettingsFile, settings)
  assert.Error(t, err)
}

func TestWriteSettingsMarshalError(t *testing.T) {
  err := util.MockHomeSetup()
  assert.NoError(t, err)

  vaults := make(map[string]VaultInstance)
  settings := Settings{
    Vaults: vaults,
  }

  original := marshalFunc
  defer func() { marshalFunc = original }()

  marshalFunc = func(v any) ([]byte, error) {
    return nil, errors.New("forced marshal error")
  }

  err = WriteSettingsFile(util.MockSettingsFile, settings)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}
