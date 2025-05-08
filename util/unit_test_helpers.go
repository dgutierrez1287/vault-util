package util

import (
	"os"
	"path/filepath"
)

/*
Helper functions for testing
*/
var MockHomeDir = "./mock/home"
var MockSettingsFile = filepath.Join(MockHomeDir,
  ".vault-util-settings.json")

func MockHomeSetup() error {
  err := os.MkdirAll(MockHomeDir, 0755)
  if err != nil {
    return err
  }
  return nil
}

func MockHomeCleanup() error {
  err := os.RemoveAll("./mock")
  if err != nil {
    return err
  }
  return nil
}

func MockCaCertFile() error {
  certText := "testcertcontent"

  file, err := os.Create(filepath.Join(MockHomeDir,
    "test-ca-cert"))
  if err != nil {
    return err
  }

  _, err = file.Write([]byte(certText))
  if err != nil {
    return err
  }

  return nil
}

func MockCaKeyFile() error {
  keyText := "testkeycontent"

  file, err := os.Create(filepath.Join(MockHomeDir,
    "test-ca-key"))
  if err != nil {
    return err
  }

  _, err = file.Write([]byte(keyText))
  if err != nil {
    return err
  }

  return nil
}


