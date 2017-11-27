package cmds

import (
	"os"

	"github.com/damekr/backer/bacli/client"
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
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}
	response, err := clnt.PingInSecure()
	if err != nil {
		log.Error("Could not ping server: ", server)
		os.Exit(1)
	}
	log.Info("Server response: ", response)
}
