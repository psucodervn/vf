package cmd

import (
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var (
  cfgFile string
  homeDir string
  depth   int
)

const (
  DefaultDepth = 3
)

func readConfig(cmd *cobra.Command, args []string) {
  if depth <= 0 {
    depth = viper.GetInt("depth")
  }
  if depth == 0 {
    depth = DefaultDepth
  }
}
