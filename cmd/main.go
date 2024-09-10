package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"

	"vpdb-dev-tool/internal/lib/cmd"
	"vpdb-dev-tool/internal/lib/logging"
	"vpdb-dev-tool/internal/tools/ssh_compose"
	"vpdb-dev-tool/internal/tools/stack"
	"vpdb-dev-tool/internal/tools/vdi"
)

var (
	Version   = "dev"
	BuildDate = "none"
	Commit    = "unknown"
)

const vString = "" +
	"   Version: %s\n" +
	"Build Date: %s\n" +
	"    Commit: %s\n"

func main() {
	tree := cli.Tree()
	opts := cmd.Opts{}

	cmd.RegisterOpts(tree, &opts)

	tree.
		WithFlag(cli.ComboFlag('v', "version").
			WithCallback(func(argo.Flag) {
				fmt.Printf(vString, Version, BuildDate, Commit)
				os.Exit(0)
			})).
		WithCallback(func(_ argo.CommandTree) {
			logging.SetupLogging(opts.LogLevel)
		})

	ssh_compose.Init(tree)
	vdi.Init(tree)
	stack.Init(tree)

	_, err := tree.Parse(os.Args)

	if err != nil {
		log.Fatalln(err)
	}
}
