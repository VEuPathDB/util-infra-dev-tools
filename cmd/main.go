package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
	"vpdb-dev-tool/internal/tools/ssh_compose"
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

	tree.WithFlag(cli.ComboFlag('v', "version").
		WithCallback(argo.FlagCallback(func(argo.Flag) {
			fmt.Printf(vString, Version, BuildDate, Commit)
			os.Exit(0)
		})))

	ssh_compose.Init(tree)

	_, err := tree.Parse(os.Args)

	if err != nil {
		log.Fatalln(err)
	}
}
