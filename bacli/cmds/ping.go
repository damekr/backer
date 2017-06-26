package cmds

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping server if is available",
	Long:  "This command sends just ping to server to check if server respond",
	RunE: func(cmd *cobra.Command, args []string) error {
		sendServerPing()
		return nil
	},
}

func sendServerPing() {
	log.Info("Sending ping to server: ", server)
}
