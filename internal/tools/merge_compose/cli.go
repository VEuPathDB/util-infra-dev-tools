package merge_compose

import (
	"fmt"
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
	"github.com/sirupsen/logrus"
	"os"
	"vpdb-dev-tool/internal/tools/merge_compose/conf"
)

const (
	envFlagDesc  = "One or more environment files to use when interpolating compose files."
	yamlFlagDesc = "One or more compose files to merge together.\n\nCompose files are merged as overlays with the" +
		" first specified file being the bottom layer and the last being the top layer.  This means that on property" +
		" conflict, the value from the later file will be used."
	versionFlagDesc = "Prints the tool's version number and exits."
	commandDesc     = "Merges compose files and applies environment substitutions recursively."
)

func Init(tree argo.CommandTreeBuilder) {
	options := conf.Options{}

	tree.WithLeaf(cli.Leaf("merge-compose").
		WithDescription(commandDesc).
		WithFlag(cli.ComboFlag('e', "env-file").
			WithDescription(envFlagDesc).
			WithArgument(cli.Argument().
				Require().
				WithBinding(&options.EnvFiles).
				WithDefault([]string{".env"}).
				WithValidator(envListValidator))).
		WithFlag(cli.ComboFlag('f', "compose-file").
			WithDescription(yamlFlagDesc).
			WithArgument(cli.Argument().
				Require().
				WithBinding(&options.ComposeFiles).
				WithDefault([]string{"docker-compose.yml"}).
				WithValidator(composeListValidator))).
		WithCallback(func(_ argo.CommandLeaf) { run(options) }))
}

func RunStandalone(version func() string) {
	options := conf.Options{}

	cli.Command().
		WithDescription(commandDesc).
		WithFlag(cli.ComboFlag('e', "env-file").
			WithDescription(envFlagDesc).
			WithArgument(cli.Argument().
				Require().
				WithBinding(&options.EnvFiles).
				WithDefault([]string{".env"}).
				WithValidator(envListValidator))).
		WithFlag(cli.ComboFlag('f', "compose-file").
			WithDescription(yamlFlagDesc).
			WithArgument(cli.Argument().
				Require().
				WithBinding(&options.ComposeFiles).
				WithDefault([]string{"docker-compose.yml"}).
				WithValidator(composeListValidator))).
		WithFlag(cli.ComboFlag('v', "version").
			WithDescription(versionFlagDesc).
			WithCallback(func(_ argo.Flag) {
				fmt.Println(version())
				os.Exit(0)
			})).
		WithCallback(func(_ argo.Command) { run(options) }).
		MustParse(os.Args)
}

func composeListValidator(value []string, _ string) error {
	if len(value) == 0 {
		logrus.Fatalf("at least one compose file is required")
	}

	return nil
}

func envListValidator(value []string, _ string) error {
	if len(value) == 0 {
		logrus.Fatalf("at least one environment file is required")
	}

	return nil
}
