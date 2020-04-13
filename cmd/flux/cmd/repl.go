package cmd

import (
	"context"

	"github.com/wolffcm/flux"
	_ "github.com/wolffcm/flux/builtin"
	"github.com/wolffcm/flux/dependencies/filesystem"
	"github.com/wolffcm/flux/repl"
	"github.com/spf13/cobra"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Launch a Flux REPL",
	Long:  "Launch a Flux REPL (Read-Eval-Print-Loop)",
	Run: func(cmd *cobra.Command, args []string) {
		deps := flux.NewDefaultDependencies()
		deps.Deps.FilesystemService = filesystem.SystemFS
		// inject the dependencies to the context.
		// one useful example is socket.from, kafka.to, and sql.from/sql.to where we need
		// to access the url validator in deps to validate the user-specified url.
		ctx := deps.Inject(context.Background())
		r := repl.New(ctx, deps)
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
