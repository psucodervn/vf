package cmd

import (
  "github.com/spf13/viper"
)

func saveConfig() {
  err := viper.SafeWriteConfig()
  if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
    err = viper.WriteConfig()
  }
  if err != nil {
    panic(err)
  }
}
