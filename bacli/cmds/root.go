package cmds

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "cmds"})

const mgmtPort = "8090"

var (
	logging     bool
	verbose     bool
	verboseLog  bool
	quiet       bool
	disableAuth bool
)

var (
	logFile  string
	server   string
	user     string
	password string
	port     string
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

// AddCommands adds child commands to the root command.
func AddCommands() {
	//RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(pingCmd)
	RootCmd.AddCommand(runCmd)
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&logging, "logger", false, "enable Logging")
	RootCmd.PersistentFlags().StringVar(&logFile, "logFile", "", "logger File path (if set, logging enabled automatically)")
	RootCmd.PersistentFlags().BoolVar(&verboseLog, "verboseLog", false, "verbose logging")
	RootCmd.PersistentFlags().StringVarP(&server, "server", "s", "127.0.0.1", "ip or host of backer server")
	RootCmd.PersistentFlags().StringVar(&user, "user", "admin", "username to authenticate to server")
	RootCmd.PersistentFlags().StringVar(&password, "password", "admin", "password to authenticate to server")
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", mgmtPort, "port for management interface of bacsrv")
	RootCmd.PersistentFlags().BoolVar(&disableAuth, "disableAuth", false, "disable authentication to server")
}
