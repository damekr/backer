package cmds

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacli/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	listCmd.AddCommand(listClients)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing out various types of content",
	Long: `Listing out various types of content.
List requires a subcommand,`,
	RunE: nil,
}

var listClients = &cobra.Command{
	Use:   "clients",
	Short: "List all clients",
	Long:  `List all of the clients`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Listing clients")
		clnt := client.ClientGRPC{
			Server: server,
			Port:   port,
		}
		clients, err := clnt.ListAllInSecure()
		if err != nil {
			log.Error("Could not list clients")
			os.Exit(1)
		}
		fmt.Println(clients)
		return nil

	},
}
