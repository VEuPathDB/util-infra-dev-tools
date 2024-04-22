package vdi_projects

var (
	plugins = [...]Project{
		{Name: "vdi-plugin-bigwig", Images: []string{"vdi-plugin-bigwig"}},
		{Name: "vdi-plugin-biom", Images: []string{"vdi-plugin-biom"}},
		{Name: "vdi-plugin-example", Images: []string{"vdi-plugin-example"}},
		{Name: "vdi-plugin-genelist", Images: []string{"vdi-plugin-genelist"}},
		{Name: "vdi-plugin-isasimple", Images: []string{"vdi-plugin-isasimple"}},
		{Name: "vdi-plugin-rnaseq", Images: []string{"vdi-plugin-rnaseq"}},
	}
	baseImages = [...]Project{
		{Name: "docker-gus-apidb-base", Images: []string{"gus-apidb-base"}},
		{Name: "vdi-plugin-handler-server", Images: []string{"vdi-plugin-handler-server"}},
		{Name: "vdi-docker-plugin-base", Images: []string{"vdi-plugin-base"}},
	}
	rootService = [...]Project{
		{Name: "vdi-service", Images: []string{"vdi-service"}},
	}
	support = [...]Project{
		{Name: "vdi-internal-db", Images: []string{"vdi-internal-db"}},
		{Name: "docker-apache-kafka", Images: []string{"apache-kafka"}},
	}
	libraries = [...]Project{
		{Name: "vdi-component-common", Images: []string{}},
		{Name: "vdi-component-json", Images: []string{}},
	}
)

type Project struct {
	Name   string
	Images []string
}

func DeployableImageProducers() []Project {
	out := make([]Project, 0, len(plugins)+len(rootService)+len(support))

	for i := range rootService {
		out = append(out, rootService[i])
	}

	for i := range support {
		out = append(out, support[i])
	}

	for i := range plugins {
		out = append(out, plugins[i])
	}

	return out
}
