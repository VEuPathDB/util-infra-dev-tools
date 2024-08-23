package stack

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"

	"vpdb-dev-tool/internal/tools/stack/use"
)

const description = "Docker Compose stack related commands and utilities."

func Init(builder argo.CommandTreeBuilder) {
	branch := cli.Branch("stack").
		WithDescription(description)

	use.Init(branch)

	builder.WithBranch(branch)
}
