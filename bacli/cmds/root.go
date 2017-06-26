package cmds

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	logging    bool
	verbose    bool
	verboseLog bool
	quiet      bool
)

var (
	logFile  string
	server   string
	user     string
	password string
)
var RootCmd = &cobra.Command{
	Use:   "bacmgr",
	Short: "Backer Server CLI Tool",
	Long:  `A tool dedicated to manage Backer Server`,
}

// Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
func Execute() {
	// RootCmd.SetGlobalNormalizationFunc(helpers.NormalizeHugoFlags)

	RootCmd.SilenceUsage = true

	AddCommands()

	if c, err := RootCmd.ExecuteC(); err != nil {
		c.Println("")
		c.Println(c.UsageString())
		os.Exit(-1)
	}
}

// AddCommands adds child commands to the root command HugoCmd.
func AddCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(pingCmd)
	RootCmd.AddCommand(listCmd)
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&logging, "log", false, "enable Logging")
	RootCmd.PersistentFlags().StringVar(&logFile, "logFile", "", "log File path (if set, logging enabled automatically)")
	RootCmd.PersistentFlags().BoolVar(&verboseLog, "verboseLog", false, "verbose logging")
	RootCmd.PersistentFlags().StringVarP(&server, "server", "s", "127.0.0.1", "ip or host of backer server")
	RootCmd.PersistentFlags().StringVarP(&user, "user", "u", "admin", "username to authenticate to server")
	RootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "password to authenticate to server")
}
