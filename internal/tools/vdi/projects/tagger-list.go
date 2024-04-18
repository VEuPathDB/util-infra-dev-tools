package vdi_projects

var (
	plugins = map[string][]string{
		"vdi-plugin-bigwig":    {"vdi-plugin-bigwig"},
		"vdi-plugin-biom":      {"vdi-plugin-biom"},
		"vdi-plugin-example":   {"vdi-plugin-example"},
		"vdi-plugin-genelist":  {"vdi-plugin-genelist"},
		"vdi-plugin-isasimple": {"vdi-plugin-isasimple"},
		"vdi-plugin-rnaseq":    {"vdi-plugin-rnaseq"},
	}
	baseImages = map[string][]string{
		"docker-gus-apidb-base":     {"gus-apidb-base"},
		"vdi-plugin-handler-server": {"vdi-plugin-handler-server"},
		"vdi-docker-plugin-base":    {"vdi-plugin-base"},
	}
	rootService = map[string][]string{
		"vdi-service": {"vdi-service"},
	}
	support = map[string][]string{
		"vdi-internal-db":     {"vdi-internal-db"},
		"docker-apache-kafka": {"apache-kafka"},
	}
	libraries = map[string][]string{
		"vdi-component-common": {},
		"vdi-component-json":   {},
	}
)

func DeployableImageProducers() map[string][]string {
	out := make(map[string][]string, len(plugins)+len(rootService)+len(support))

	for k, v := range plugins {
		out[k] = v
	}
	for k, v := range rootService {
		out[k] = v
	}
	for k, v := range support {
		out[k] = v
	}

	return out
}
