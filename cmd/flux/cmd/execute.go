package cmd

import (
	"context"
	"fmt"

	"github.com/wolffcm/flux"
	_ "github.com/wolffcm/flux/builtin"
	"github.com/wolffcm/flux/dependencies/filesystem"
	"github.com/wolffcm/flux/repl"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a Flux script",
	Long:  "Execute a Flux script from string or file (use @ as prefix to the file)",
	Args:  cobra.ExactArgs(1),
	RunE:  execute,
}

func init() {
	rootCmd.AddCommand(executeCmd)
}

func execute(cmd *cobra.Command, args []string) error {
	deps := flux.NewDefaultDependencies()
	deps.Deps.FilesystemService = filesystem.SystemFS
	ctx := deps.Inject(context.Background())
	r := repl.New(ctx, deps, querier{})
	if err := r.Input(args[0]); err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	return nil
}
