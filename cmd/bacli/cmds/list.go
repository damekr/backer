package cmds

import (
	"os"

	"github.com/damekr/backer/cmd/bacli/client"
	"github.com/spf13/cobra"
)

/*
CLI Options in list
./bacli list -o|--operation <operation> -c|--client <client_IP>
operations:
 - backups
 - clients
 - ...
*/

func init() {
	backupsList.Flags().StringVarP(&clientIP, "client", "c", "", "ip of client to be listed backups")
	backupsList.MarkFlagRequired("client")
	listCmd.AddCommand(backupsList)
	listCmd.AddCommand(clientsList)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing out various types of content",
	Long: `Listing out various types of content.
List requires a subcommand,`,
	RunE: nil,
}

var backupsList = &cobra.Command{
	Use:   "backups",
	Short: "Listing out backups of given client",
	Long:  "Listing out backups of given client, client must be specified",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listBackups()
	},
}

var clientsList = &cobra.Command{
	Use:   "clients",
	Short: "Listing out clients",
	Long:  "Listing out clients which have any backups visible by this server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listClients()
	},
}

func listBackups() error {
	log.Println("Listing backups of client: ", clientIP)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}

	err := clnt.ListBackupsInSecure(clientIP)
	if err != nil {
		log.Error("Could not list backups of client, err: ", err)
		os.Exit(1)
	}
	return nil
}

func listClients() error {
	log.Println("Listing clients")
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}

	err := clnt.ListClientsInSecure()
	if err != nil {
		log.Error("Could not list backups of client, err: ", err)
		os.Exit(1)
	}
	return nil
}

func listClient() {

}
