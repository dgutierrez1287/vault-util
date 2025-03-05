package util

import (
  "strings"

  "github.com/dgutierrez1287/vault-util/logger"
)

type Secret struct {
  VaultKey string                     `json:"key"`
  SecretKey string
  SecretPath string
  SecretData map[string]interface{}   `json:"data"`
}

func (s *Secret) splitPath() {
  logger.Logger.Debug("splitting the vault path into secret key and secret path")
  logger.Logger.Debug("vault key", s.VaultKey)

  parts := strings.SplitN(s.VaultKey, "/", 2)
  if len(parts) < 2 {
    logger.Logger.Error("error splitting vault path")
  }

  s.SecretKey = parts[0]
  s.SecretPath = parts[1]

  logger.Logger.Debug("the resulting split vault key", "key", s.SecretKey, "path", s.SecretPath)
}

type VaultSecrets struct {
  Secrets map[string]Secret           `json:"secrets"`
}

