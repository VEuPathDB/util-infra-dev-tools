package vdi_projects

var (
	plugins = [...]Project{
		{Name: "vdi-plugin-bigwig", Images: []string{"vdi-plugin-bigwig"}},
		{Name: "vdi-plugin-biom", Images: []string{"vdi-plugin-biom"}},
		{Name: "vdi-plugin-example", Images: []string{"vdi-plugin-example"}},
		{Name: "vdi-plugin-genelist", Images: []string{"vdi-plugin-genelist"}},
		{Name: "vdi-plugin-wrangler", Images: []string{"vdi-plugin-wrangler"}},
		{Name: "vdi-plugin-rnaseq", Images: []string{"vdi-plugin-rnaseq"}},
	}
	rootService = [...]Project{
		{Name: "vdi-service", Images: []string{"vdi-service", "vdi-internal-db"}},
	}
	support = [...]Project{
		{Name: "docker-apache-kafka", Images: []string{"apache-kafka"}},
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
