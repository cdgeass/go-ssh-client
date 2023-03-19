package cmd

import (
	"fmt"
	"github.com/cdgeass/ssh-client/config"
	"github.com/spf13/cobra"
	"os"
)

var conf config.Config

var rootCmd = &cobra.Command{
	Use:   "ssh-client",
	Short: "A ssh client written by golang",
	Long:  `A ssh client written by golang, manage multi server info connect without inputting password`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	conf = config.Load()
}

var (
	index    int
	name     string
	host     string
	port     int
	user     string
	password string
)
