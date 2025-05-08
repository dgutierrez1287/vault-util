package app

import (
	"encoding/json"
)

/*
Interface to to cover all different types
of machine output
*/
type MachineOutput interface {
  GetOutputJson() (string, int)
}

/*
VaultListOutput - Machine output for 
listing vaults that are in settings
*/
type VaultListOutput struct {
  ExitCode int              `json:"exitCode"`
  Vaults []string           `json:"vaults"`
  Message string            `json:"message,omitempty"`
}

func (v VaultListOutput) GetOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(v)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

/*
MountListOutput - Machine output for 
listing secret mounts 
*/
type MountListOutput struct {
  ExitCode int                            `json:"exitCode"`
  MountNames []string                     `json:"mountNames,omitempty"`
  MountsWithData map[string]interface{}   `json:"mounts,omitempty"`
}

func (m MountListOutput) GetOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(m)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

/*
GetSecretOutput - Machine output for 
get secret
*/
type GetSecretOutput struct {
  ExitCode int                  `json:"exitCode"`
  SecretExists bool             `json:"secretExists"`
  VaultKey string               `json:"secretKey,omitempty"`
  Data map[string]interface{}   `json:"secretData,omitempty"`
}

func (s GetSecretOutput) GetOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(s)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

/*
SecretListOutput - machine output for 
secret list
*/
type SecretListOutput struct {
  ExitCode int              `json:"exitCode"`
  Secrets []string          `json:"secrets"`
}

func (s SecretListOutput) GetOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(s)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

/*
AddRemoveOuput - Machine output for
adding and removing vaults from config
*/
type AddRemoveOutput struct {
  ExitCode int                `json:"exitCode"`
  Message string              `json:"message"`
}

func (a AddRemoveOutput) GetOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(a)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

/*
BulkActionOutput - Machine output for
bulk actions
*/
type BulkActionOutput struct {
  ExitCode int                  `json:"exitCode"`
  SecretsAdded []string         `json:"secretsAdded,omitempty"`
  SecretsRemoved []string       `json:"secretsRemoved,omitempty"`
  Errors []SecretActionError    `json:"Errors,omitempty"`
}

func (b BulkActionOutput) GetOutputJson() (string, int) {
  jsonBtyes, err := json.Marshal(b)
  if err != nil {
    return "\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100
  }
  return string(jsonBtyes), 0
}


