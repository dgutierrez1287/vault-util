package app

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

/*
Console output for bulk action
*/
func BulkActionConsoleOutput(successList []string, 
  errorList []SecretActionError, action string) {

  fmt.Println("Bulk Action Results")
  fmt.Println("===========================")
  
  if action == "delete" {
    fmt.Printf("%d secrets removed\n",len(successList))
  } else {
    fmt.Printf("%d secrets added/updated\n",len(successList))
  }
  fmt.Println("")
  fmt.Println("The following secrets had errors")
  for _, errorSecret := range errorList {
    fmt.Printf("key: %s, error: %s", errorSecret.VaultKey, errorSecret.Error)
  }
}

/*
Console output for listing vaults
*/
func ListVaultsConsoleOutput(settingsExist bool, vaults []string) {
  fmt.Println("Vaults Configured")
  fmt.Println("============================")
  
  if !settingsExist{
    fmt.Println("No settings file is present, no vaults configured")
    return
  }

  fmt.Println("Vaults:")
  for _, name := range vaults {
    fmt.Println(name)
  }
}

/*
Console output for Get secret
*/
func GetSecretConsoleOutput(secret VaultSecret, secretExists bool) {
  fmt.Println("Get Secret Results")
  fmt.Println("==============================")

  if !secretExists {
    fmt.Println("Secret does not exist")
    return
  }

  fmt.Println("Key: " + secret.NormalizedSecretPath)
  fmt.Println("")
  fmt.Println("data:")
  
  for key, val := range secret.SecretData {
    fmt.Println("Key: " + key)
    fmt.Println("Value: " + val.(string))
    fmt.Println("")
  }
}

/*
Console output for listing secrets
*/
func ListSecretsConsoleOutput(secrets []string, mountName string) {
  fmt.Printf("Secrets for mount %s\n", mountName)
  fmt.Println("============================")

  fmt.Println("Secrets:")
  for _, key := range secrets {
    fmt.Println(key)
  }
}

/*
Console output for listing mount names
*/
func ListMountNamesConsoleOutput(mounts []SecretMount) {
  fmt.Println("Secret Mounts")
  fmt.Println("===============================")

  fmt.Println("Mounts:")
  for _, mount := range mounts {
    fmt.Println(mount.Mount)
  }
}

/*
Console output for listing mounts with detail
*/
func ListMountsWithDetailConsoleOutput(mounts []SecretMount) {
  fmt.Println("Secret Mounts")
  fmt.Println("===========================")
  fmt.Println()

  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{"Mount", "Description", "Type", "Version"})
  table.SetAlignment(tablewriter.ALIGN_LEFT)
  table.SetRowLine(true)
  table.SetAutoWrapText(false)

  for _, mount := range mounts {
    if mount.Type == "kv" {
      table.Append([]string{mount.Mount, mount.Description, mount.Type, mount.KvVersion})
    } else {
      table.Append([]string{mount.Mount, mount.Description, mount.Type, "N/A"})
    }
  }

  table.Render()
}
