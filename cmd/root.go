/*
Copyright © 2022 GRISARD Dimitri dimitri.grisard03@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host, user, password, typeOf, cfgFile string
	port                                  uint16
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-ucce",
	Version: "1.0.0",
	Short:   "go-ucce is a CLI tool to query UCCE Cisco appliance.",
	Long: `go-ucce is a CLI tool to query UCCE Cisco appliance.

This CLI send command via SSH, you need to specify : host, port, user, password and instance type`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ucce-cisco.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&host, "host", "a", "", "Hostname or IP address targeted")
	rootCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 22, "Ssh port used")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "User used to login")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "x", "", "Password used to login")
	rootCmd.PersistentFlags().StringVarP(&typeOf, "typeOf", "t", "", "Type of UCCE Instance (Finesse, Cuic...)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".uzo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ucce")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
