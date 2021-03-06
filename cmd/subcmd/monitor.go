// Package subcmd is for all subcmds in the cmd tree
package subcmd

import (
	"github.com/milligan22963/radio/pkg/monitor"
	"github.com/milligan22963/radio/pkg/util"

	"github.com/spf13/cobra"
)

// MonitorCmd is the main gpio monitor cmd
var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitors gpios and plays music",
	Long:  `The primary application command`,
	Run: func(cmd *cobra.Command, args []string) {
		utilities := util.Util{}
		utilities.SetupConfiguration("config.yaml")
		utilities.SetupLogging()

		monitorInstance := monitor.MonitorCmd{}

		monitorInstance.Run(utilities)
	},
}
