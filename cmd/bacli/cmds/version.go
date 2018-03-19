package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of the tool",
	Long:  "0.1 version",
	RunE: func(cmd *cobra.Command, args []string) error {
		printBacliVersion()
		return nil
	},
}

func printBacliVersion() {
	fmt.Println("Version 0.1")
}
