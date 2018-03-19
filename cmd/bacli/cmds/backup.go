package cmds

import (
	"os"

	"github.com/damekr/backer/cmd/bacli/client"
	"github.com/spf13/cobra"
)

/*
CLI Options in backup
./bacli backup -c|--client <client_IP> -d|--dirs <path_1>...

*/

var backupPaths []string

var runBackup = &cobra.Command{
	Use:   "backup",
	Short: "Run a backup job",
	Long:  `Run backup job with specified client name`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startBackup()
	},
}

func init() {
	runBackup.Flags().StringVarP(&clientIP, "client", "c", "", "ip of client to be backed up")
	runBackup.MarkFlagRequired("client")
	runBackup.Flags().StringSliceVarP(&backupPaths, "dirs", "d", []string{}, "dirs to be backed up, might be also one file")
	runBackup.MarkFlagRequired("paths")
}

func startBackup() error {
	log.Debugln("Starting backup of client: ", clientIP)
	log.Println("Running backup of paths: ", backupPaths)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}

	err := clnt.RunBackupInSecure(clientIP, backupPaths)
	if err != nil {
		log.Error("Could not run backup of client")
		os.Exit(1)
	}
	return nil
}
