package cmds

import (
	"os"

	"github.com/damekr/backer/cmd/bacli/client"
	"github.com/spf13/cobra"
)

/*
CLI Options in restore

./bacli restore -c|--client <clientIP> -i|--backupid <backupID> -- Restore whole backup to given client into the same place
./bacli restore -c|--client <clientIP> -i|--backupid <backupID> -r|--remote <restore to remote path> -- Restore whole backup to given client into specified path
./bacli restore -c|--client <clientIP> -i|--backupid <backupID> -d|--dir <path to be restored> -- Restore dir or whole path to the same location
./bacli restore -c|--client <clientIP> -i|--backupid <backupID> -d|--dir <path to be restored> -r|--remote <restore to remote path> --
Restore dir or whole path to given directory of the client

*/
var (
	backupID  int64
	dir       string
	remoteDir string
)

var runRestore = &cobra.Command{
	Use:   "restore",
	Short: "Run restore job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if clientIP == "" || backupID == 0 {
			log.Errorln("Please specify at least clientIP and backupID")
			os.Exit(1)
		}
		switch {
		case dir != "" && remoteDir != "":
			log.Debugf("Running restore of one dir: %s to different place: %s\n", dir, remoteDir)
			return restoreDirIntoDifferentPlace(clientIP, dir, remoteDir, backupID)
		case dir != "" && remoteDir == "":
			log.Debugf("Running restore of one dir: %s to the same place\n", dir)
			return restoreDirIntoSamePlace(clientIP, dir, backupID)
		case remoteDir != "":
			log.Debugln("Running restore of whole backup to different place:", remoteDir)
			return restoreWholeBackupIntoDifferentPlace(clientIP, remoteDir, backupID)
		case clientIP != "" && backupID != 0:
			log.Debugln("Running restore of whole backup of to the same place")
			return restoreWholeBackupIntoSamePlace(clientIP, backupID)
		}
		return nil

	},
}

func init() {
	runRestore.Flags().StringVarP(&clientIP, "client", "c", "", "ip of client to be restored")
	runRestore.MarkFlagRequired("client")
	runRestore.Flags().Int64VarP(&backupID, "backupid", "i", 0, "id of existing backup")
	runRestore.MarkFlagRequired("backupid")
	runRestore.Flags().StringVarP(&dir, "dir", "d", "", "filename or directory from created backup to be restored")
	runRestore.Flags().StringVarP(&remoteDir, "remote", "r", "", "remote path to restore dir or directory")
}

func restoreWholeBackupIntoSamePlace(clientIP string, backupID int64) error {
	log.Debugln("Running restore of client: ", clientIP)
	log.Println("Running restore of backupID: ", backupID)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}
	err := clnt.RunRestoreWholeBackupInSecure(clientIP, backupID)
	if err != nil {
		log.Error("Could not run backup of client")
		os.Exit(1)
	}

	return nil
}

func restoreWholeBackupIntoDifferentPlace(clientIP, remoteDir string, backupID int64) error {
	log.Debugln("Running restore of client: ", clientIP)
	log.Println("Running restore of backupID: ", backupID)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}
	err := clnt.RunRestoreWholeBackupDifferentPlaceInSecure(clientIP, remoteDir, backupID)
	if err != nil {
		log.Error("Could not run backup of client")
		os.Exit(1)
	}

	return nil
}

func restoreDirIntoSamePlace(clientIP, dir string, backupID int64) error {
	log.Debugln("Running restore of client: ", clientIP)
	log.Println("Running restore of backupID: ", backupID)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}
	err := clnt.RunRestoreOfDirInSecure(clientIP, dir, backupID)
	if err != nil {
		log.Error("Could not run backup of client")
		os.Exit(1)
	}
	return nil
}

func restoreDirIntoDifferentPlace(clientIP, dir, remoteDir string, backupID int64) error {
	log.Debugln("Running restore of client: ", clientIP)
	log.Println("Running restore of backupID: ", backupID)
	clnt := client.ClientGRPC{
		Server: server,
		Port:   port,
	}
	err := clnt.RunRestoreOfDirDifferentPlaceInSecure(clientIP, dir, remoteDir, backupID)
	if err != nil {
		log.Error("Could not run backup of client")
		os.Exit(1)
	}
	return nil
}
