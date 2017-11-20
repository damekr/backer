package cmds

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacli/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	runCmd.AddCommand(runBackup)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a job",
	Long: `Run starting a new job.
Run requires a subcommand`,
	RunE: nil,
}

var runBackup = &cobra.Command{
	Use:   "backup",
	Short: "Run a backup job",
	Long:  `Run backup job with specified client name`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Error("Client(s) name not specified, exiting...")
			os.Exit(2)
		}
		log.Println("Running backup of: ", args[1:])
		clnt := client.ClientGRPC{
			Server: server,
			Port:   port,
		}
		err := clnt.RunBackupInSecure(args[1:])
		if err != nil {
			log.Error("Could not run backup of client")
			os.Exit(1)
		}
		return nil

	},
}
