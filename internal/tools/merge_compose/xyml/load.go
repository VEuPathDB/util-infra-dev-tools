package xyml

import (
	"gopkg.in/yaml.v3"
	"os"
	"vpdb-dev-tool/internal/lib/must"
	"vpdb-dev-tool/internal/lib/xos"
)

func LoadDocument(path string) *yaml.Node {
	file := xos.MustOpen(path, os.O_RDONLY, 0644)
	defer xos.MustClose(file)

	var node yaml.Node

	dec := yaml.NewDecoder(file)
	must.NotError(dec.Decode(&node))

	return &node
}
