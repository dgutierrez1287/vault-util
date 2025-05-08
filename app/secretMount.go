package app

import (
	"errors"
	"strings"

	"github.com/dgutierrez1287/vault-util/logger"
)

/*
SecretMount - The components of a secret
mount
*/
type SecretMount struct {
  Mount string            `json:"mount"`
  Type string             `json:"type"`
  Description string      `json:"description,omitempty"`
  KvVersion string        `json:"kvVersion,omitempty"`
}

func NewSecretMount(name string, mountType string, 
description string, version string, client *VaultClient) (SecretMount, error) {
  var err error

  logger.LogDebug("Verifying the mount is in the correct format")
  if name[len(name)-1] != '/' {
    name = name + "/"
  }

  if mountType == "" {
    logger.LogDebug("Getting mount type since it was not supplied")
    mountType, version, err = GetMountType(client, name)

    if err != nil {
      logger.LogError("Error getting information on the mount type")
      return SecretMount{}, err
    }
  }

  mount := SecretMount{
    Mount: name,
    Type: mountType,
    KvVersion: version,
    Description: description,
  }

  return mount, nil
}

/*
This will get a list of all the secrets for a 
given mount
*/
func (sm SecretMount) ListSecrets(client *VaultClient) ([]string,
  error) {

  var secrets []string

  if sm.Type == "kv" {
    logger.LogDebug("Secret mount is a kv")

    var walk func(string) error
    walk = func(path string) error {
      respKeys, err := client.ListKvSecrets(sm.Mount, path, sm.KvVersion)
      if err != nil {
        // check for 404 to skip non folder
        if strings.Contains(err.Error(), "404") {
          return nil
        }
        return err
      }

      for _, key := range respKeys {
        if strings.HasSuffix(key, "/") {
          err := walk(path + key)
          if err != nil {
            return err
          }
        } else {
          //leaf secret, add to list
          fullPath := path + key
          secrets = append(secrets, fullPath)
        }
      }
      return nil
    }

    err := walk(sm.Mount)
    if err != nil {
      return nil, err
    }
  }
  return secrets, nil 
}

/*
This will get the type for a certain secrets 
engine, if the type is kv then it will also return
the kv version
*/
func GetMountType(client *VaultClient, mountName string) (string, 
  string, error) {

  logger.LogDebug("Getting data for all mounts")
  mounts, err := client.GetSecretMountsData()

  if err != nil {
    logger.LogError("Error getting secret mounts data")
    return "", "", err
  }

  logger.LogDebug("Getting mount data for", "mount", mountName)
  mountData, ok := mounts[mountName].(map[string]interface{})

  if !ok {
    logger.LogError("Error mount does not exist")
    return "", "", errors.New("secrets mount doesn't exist")
  }

  logger.LogDebug("Mount data", "data", mountData)
  
  mountType, ok := mountData["type"].(string)
  
  if !ok {
    logger.LogError("Error getting secret mount type")
    return "", "", errors.New("cannot get secret mount type")
  }

  if mountType != "kv" {
    logger.LogDebug("Secret mount type is not a KV type")
    return mountType, "", nil
  }
  
  logger.LogDebug("Secret mount type is KV, getting the KV version")
  options, ok := mountData["options"].(map[string]interface{})

  if !ok {
    logger.LogError("Error getting the options for the secret mount")
    return mountType, "", errors.New("cannot get secret mount options")
  }

  kvVersion, ok := options["version"].(string)

  if !ok {
    logger.LogError("Error getting kv version type")
    return mountType, "", errors.New("cannot get kv version") 
  }

  return mountType, kvVersion, nil
}

/*
This will get a list of secret mounts for output
*/
func GetSecretMounts(client *VaultClient) ([]SecretMount, error) {
  var mounts []SecretMount

  logger.LogDebug("Getting data for all mounts")
  rawMounts, err := client.GetSecretMountsData()

  if err != nil {
    logger.LogError("Error getting secret mounts data")
    return mounts, err
  }

  logger.LogDebug("Processing mounts")
  for rawMount, rawData := range rawMounts {
    mount := SecretMount{}

    mount.Mount = rawMount
    data := rawData.(map[string]interface{})

    mount.Type = data["type"].(string)
    mount.Description = data["description"].(string)

    if mount.Type == "kv" {
      logger.LogDebug("Secret mount is kv, getting the kv version")
      options, ok := data["options"].(map[string]interface{})

      if !ok {
        logger.LogError("Error getting mount options")
        return mounts, errors.New("cannot get secret mount options")
      }

      mount.KvVersion = options["version"].(string)
    }

    logger.LogDebug("mount", "mount", mount.Mount, "type", mount.Type, 
      "kvVersion", mount.KvVersion, "description", mount.Description)

    mounts = append(mounts, mount)
  }
  return mounts, nil
}

/*
This will convert an array of secret mounts to a map of
[string]interface{} to make json output cleaner without
another type
*/
func MountstoMap(mounts []SecretMount) (map[string]interface{}) {
  outputMap := make(map[string]interface{})

  logger.LogDebug("converting array of mounts to map")
  for _, mount := range mounts {
    logger.LogDebug("processing mount", "mount", mount.Mount)

    mountDataMap := make(map[string]string)
    mountDataMap["type"] = mount.Type
    if mount.Type == "kv" {
      mountDataMap["version"] = mount.KvVersion
    }
    mountDataMap["description"] = mount.Description

    outputMap[mount.Mount] = mountDataMap
  }
  logger.LogDebug("output mount map", "output", outputMap)
  return outputMap
} 

/*
This will return just the mount names from a array
of secret mounts
*/
func GetMountNames(mounts []SecretMount) ([]string) {
  var mountNames []string
  
  logger.LogDebug("Getting list of mount names")
  for _, mount := range mounts {
    mountNames = append(mountNames, mount.Mount)
  }
  logger.LogDebug("mount names", "names", mountNames)
  return mountNames
}



