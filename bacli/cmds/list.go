package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
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

		fmt.Println("Listing clients")

		return nil

	},
}
