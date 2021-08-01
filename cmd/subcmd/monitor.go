// Package subcmd is for all subcmds in the cmd tree
package subcmd

import (
	"pkg/monitor"

	"github.com/spf13/cobra"
)

// MonitorCmd is the main gpio monitor cmd
var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitors gpios and plays music",
	Long:  `The primary application command`,
	Run: func(cmd *cobra.Command, args []string) {
		monitorInstance := monitor.MonitorCmd{}

		monitorInstance.Run()
	},
}
