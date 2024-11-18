package main

import (
	"fmt"
	"vpdb-dev-tool/internal/tools/merge_compose"
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
	merge_compose.RunStandalone(func() string { return fmt.Sprintf(vString, Version, BuildDate, Commit) })
}
