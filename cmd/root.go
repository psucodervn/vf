package cmd

import (
  "fmt"
  "os"
  "strings"

  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"

  "vf/internal"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "vf",
  Short: "Quickly navigate between your projects",
  PersistentPreRun: readConfig,
  Run: func(cmd *cobra.Command, args []string) {
    root := homeDir
    if len(args) > 0 {
      root = strings.TrimSpace(args[0])
    }
    f := &internal.Finder{
      Depth: depth,
    }
    f.Run(root)
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vf.yaml)")
  rootCmd.PersistentFlags().IntVarP(&depth, "depth", "d", 0, "max depth")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

  // Find home directory.
  home, err := os.UserHomeDir()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  homeDir = home

  if cfgFile != "" {
    viper.SetConfigFile(cfgFile)
  } else {
    viper.AddConfigPath(home)
    viper.SetConfigType("yaml")
    viper.SetConfigName(".vf")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    // fmt.Println("Using config file:", viper.ConfigFileUsed())
    log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
  }
}
