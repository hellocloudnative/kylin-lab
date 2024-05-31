package bin

import (
	"errors"
	"github.com/spf13/cobra"
	log "github.com/wonderivan/logger"
	"kylin-lab/server"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "loopy",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `loopy`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `loopy V1.1.0, -h for more help`
		log.Info("%s\n", usageStr)
	},
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
