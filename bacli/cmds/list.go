package cmds

import (
	"os"

	"github.com/damekr/backer/bacli/client"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listBackup)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing out various types of content",
	Long: `Listing out various types of content.
List requires a subcommand,`,
	RunE: nil,
}

var listBackup = &cobra.Command{
	Use:   "backups",
	Short: "List all backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Listing clients")
		clnt := client.ClientGRPC{
			Server: server,
			Port:   port,
		}
		clientName := ""
		if len(args) > 0 {
			// Just one first client name
			clientName = args[0]
		}
		err := clnt.ListBackupsInSecure(clientName)
		if err != nil {
			log.Error("Could not list clients")
			os.Exit(1)
		}
		return nil

	},
}
