package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/rossi1/ensync-cli/internal/api"
	"github.com/rossi1/ensync-cli/internal/config"
)

var (
	cfgFile string
	debug   bool
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "ensync",
		Short: "EnSync CLI tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize logger based on debug flag
			var logger *zap.Logger
			var err error
			if debug {
				logger, err = zap.NewDevelopment()
			} else {
				logger, err = zap.NewProduction()
			}
			if err != nil {
				panic(err)
			}
			zap.ReplaceGlobals(logger)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ensync/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")

	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}

	client := api.NewClient(
		cfg.BaseURL,
		cfg.APIKey,
		api.WithLogger(zap.L()),
		api.WithRateLimit(10, 20),
	)

	rootCmd.AddCommand(
		newEventCmd(client),
		newAccessKeyCmd(client),
		newVersionCmd(),
	)

	return rootCmd.Execute()
}
