package cmd

import (
	"fmt"

	"github.com/labbsr0x/health-checker/checker"

	"github.com/labbsr0x/health-checker/config"
	"github.com/labbsr0x/health-checker/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the serve command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the health checker and its HTTP REST APIs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := new(config.Builder).InitFromViper(viper.GetViper())

		server := new(web.Server).InitFromBuilder(builder)
		checker := new(checker.Checker).InitFromBuilder(builder)

		go checker.Run()

		err := server.Run()
		if err != nil {
			return fmt.Errorf("An error occurred while setting up the Health Checker Server: %v", err)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	config.AddFlags(startCmd.Flags())

	err := viper.GetViper().BindPFlags(startCmd.Flags())
	if err != nil {
		panic(err)
	}
}
