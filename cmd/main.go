package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ESP_Radio",
		Short: "A faux radio for escape rooms",
		Long:  `Vintage radio track playing program triggered by GPIOS`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			logrus.Info("starting up")
		},
	}
	rootCmd.AddCommand(subcmd.MonitorCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("error executing cmd: %v", err)
	}
}
