package app

import (
	"path/filepath"
	"testing"
  "os"

	"github.com/dgutierrez1287/vault-util/util"
	"github.com/stretchr/testify/assert"
)

/*
   Tests for NewVault
*/

func TestNewVault(t *testing.T) {
  url := "https://testvault.com"
  token := "faketoken"
  skipTls := true
  caCertFilePath := ""
  caCertKeyPath := ""

  vault, err := NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.NoError(t, err)
  assert.Equal(t, vault.Url, url)
  assert.Equal(t, vault.Token, token)
  assert.True(t, vault.SkipTLSVerify)
  assert.Equal(t, vault.CACert, "")
  assert.Equal(t, vault.CACertKey, "")
}

func TestNewVaultNoUrl(t *testing.T) {
  url := ""
  token := "faketoken"
  skipTls := true
  caCertFilePath := ""
  caCertKeyPath := ""

  _, err := NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.Error(t, err)
}

func TestNewVaultNoToken(t *testing.T) {
  url := "https://testvault.com"
  token := ""
  skipTls := true
  caCertFilePath := ""
  caCertKeyPath := ""

  _, err := NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.Error(t, err)
}

func TestNewVaultCaProvided(t *testing.T) {
  url := "https://testvault.com"
  token := "faketoken"
  skipTls := false
  caCertFilePath := filepath.Join(util.MockHomeDir, 
    "test-ca-cert")
  caCertKeyPath := filepath.Join(util.MockHomeDir, 
    "test-ca-key")

  err := util.MockHomeSetup()
  assert.NoError(t, err)

  err = util.MockCaCertFile()
  assert.NoError(t, err)

  err = util.MockCaKeyFile()
  assert.NoError(t, err)

  vault, err := NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.NoError(t, err)
  assert.Equal(t, vault.Url, url)
  assert.Equal(t, vault.Token, token)
  assert.False(t, vault.SkipTLSVerify)
  assert.Equal(t, vault.CACert, "testcertcontent")
  assert.Equal(t, vault.CACertKey, "testkeycontent")

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestNewVaultMissingCaCertFile(t *testing.T) {
  url := "https://testvault.com"
  token := ""
  skipTls := false
  caCertFilePath := filepath.Join(util.MockHomeDir, 
    "test-ca-cert")
  caCertKeyPath := filepath.Join(util.MockHomeDir, 
    "test-ca-key")

  err := util.MockHomeSetup()
  assert.NoError(t, err)

  _, err = NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestNewVaultMissingCaKeyFile(t *testing.T) {
  url := "https://testvault.com"
  token := ""
  skipTls := false
  caCertFilePath := filepath.Join(util.MockHomeDir, 
    "test-ca-cert")
  caCertKeyPath := filepath.Join(util.MockHomeDir, 
    "test-ca-key")

  err := util.MockHomeSetup()
  assert.NoError(t, err)

  err = util.MockCaCertFile()
  assert.NoError(t, err)

  _, err = NewVault(url, token, skipTls, caCertFilePath, caCertKeyPath)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

/*
    Tests for GetVaultConfigFromSettings
*/
func TestVaultFromConfig(t *testing.T) {
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

  vaultResp, err := GetVaultConfigFromSettings("test-vault", 
    util.MockSettingsFile)

  assert.NoError(t, err)
  assert.Equal(t, vaultResp.Url, "https://testvault.com")
  assert.Equal(t, vaultResp.Token, "faketoken")

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestVaultFromConfigNoSettingsFile(t *testing.T) {
  _, err := GetVaultConfigFromSettings("test-vault",
    util.MockSettingsFile)

  assert.Error(t, err)
}

func TestVaultFromConfigErrorRead(t *testing.T) {
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

  _, err = GetVaultConfigFromSettings("test-vault",
    util.MockSettingsFile)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}

func TestVaultFromConfigNoVault(t *testing.T) {
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

  _, err = GetVaultConfigFromSettings("notexistvault",
    util.MockSettingsFile)
  assert.Error(t, err)

  err = util.MockHomeCleanup()
  assert.NoError(t, err)
}
