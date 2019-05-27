package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/llvm"
	"github.com/spf13/cobra"

	gollvm "github.com/llvm-mirror/llvm/bindings/go/llvm"
)

var llvmCmd = &cobra.Command{
	Use:   "llvm",
	Short: "Compile a Flux script into its llvm IR",
	Long:  "Compile a Flux script into its llvm IR",
	Args:  cobra.ExactArgs(1),
	RunE:  llvmE,
}

func init() {
	rootCmd.AddCommand(llvmCmd)
}

func llvmE(cmd *cobra.Command, args []string) error {
	scriptBytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}
	script := string(scriptBytes)

	astPkg, err := flux.Parse(script)
	if err != nil {
		return err
	}

	mod, err := llvm.Build(astPkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mod.Dump()

	engine, err := gollvm.NewExecutionEngine(mod)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	funcResult := engine.RunFunction(mod.NamedFunction("main"), []gollvm.GenericValue{})
	fmt.Printf("%d\n", funcResult.Int(false))
	return nil
}
