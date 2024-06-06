package bin

import (
	"errors"
	"github.com/spf13/cobra"
	log "github.com/wonderivan/logger"
	"kylin-lab/cmd/migrate"
	"kylin-lab/server"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "kylin-lab",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `kylin-lab`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `kylin-lab V1.1.0, -h for more help`
		log.Info("%s\n", usageStr)
	},
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
