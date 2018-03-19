package cmds

import (
	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Run a job",
	Long: `Run starting a new job.
Run requires a subcommand`,
	RunE: nil,
}
