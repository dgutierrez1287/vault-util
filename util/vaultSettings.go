package util

import (
  vault "github.com/hashicorp/vault-client-go"
)

type VaultSettings struct {
  VaultUrl string
  Token string
  CaCertFilePath string
  CaCertKeyFilePath string
  SkipTlsVerify bool
  UseTlsConfig bool
  TlsConfig vault.TLSConfiguration
}
