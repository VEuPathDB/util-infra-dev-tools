package merge_compose

import (
	"fmt"
	"regexp"
	"strings"
	"vpdb-dev-tool/internal/lib/must"
	"vpdb-dev-tool/internal/tools/merge_compose/conf"
	"vpdb-dev-tool/internal/tools/merge_compose/env"
	"vpdb-dev-tool/internal/tools/merge_compose/xyml"
)

func run(options conf.Options) {
	fmt.Println(injectEnvironment(buildEnvMap(options.EnvFiles), buildMergedYAML(options.ComposeFiles)))
}

func buildEnvMap(files []string) map[string]string {
	output := make(map[string]string, 256)

	env.ProcessUserEnv(output)
	env.LoadMany(files, output)

	return output
}

func buildMergedYAML(files []string) string {
	merged := xyml.LoadDocument(files[0])

	for i := 1; i < len(files); i++ {
		merged = xyml.MergeNodes(merged, xyml.LoadDocument(files[i]))
	}

	return xyml.Stringify(merged)
}

var (
	dumbPattern  = must.Return1(regexp.Compile(`\$\{(\w+)}|\$(\w+)`))
	smartPattern = must.Return1(regexp.Compile(`\$\{(\w+)(:-|:\?)([^}]*)}`))
)

func injectEnvironment(env map[string]string, doc string) string {
	buf := strings.Builder{}
	buf.Grow(len(doc) * 2)

	// Track the last position in the source document that we have written to the
	// buffer.
	lastPos := 0

	// first pass, replace the dumb patterns
	hits := dumbPattern.FindAllStringSubmatchIndex(doc, -1)
	for _, hit := range hits {
		// Write whatever was between the last hit and here.
		buf.WriteString(doc[lastPos:hit[0]])

		var key string
		if hit[2] == -1 {
			key = doc[hit[4]:hit[5]]
		} else {
			key = doc[hit[2]:hit[3]]
		}

		// try and find a value for the env var in the env map
		if value, ok := env[key]; ok {
			// we have a match; write out the value instead of the env var.
			buf.WriteString(value)
		} else {
			// we don't have a match
		}

		// update the last position to the next char after whatever we wrote to the
		// buffer for the env var.
		lastPos = hit[1]
	}
	// Write the remainder of the doc
	buf.WriteString(doc[lastPos:])

	// Replace the doc string with the updated string that has the env vars
	// replaced.
	doc = buf.String()

	// Reset for round 2
	buf.Reset()
	lastPos = 0

	// second pass, replace the complex patterns
	hits = smartPattern.FindAllStringSubmatchIndex(doc, -1)
	for _, hit := range hits {
		// hit = [
		//   0 -> start of full pattern hit
		//   1 -> end of full pattern hit (index 1 PAST last char)
		//   2 -> start wanted env var name
		//   3 -> end of wanted env var name (index 1 PAST last char)
		//   4 -> start of switch type
		//   5 -> end of switch type (index 1 PAST last char)
		//   6 -> start of alternate
		//   7 -> end of alternate (index 1 PAST last char)
		// ]

		// Write whatever was between the last hit and here.
		buf.WriteString(doc[lastPos:hit[0]])

		// try and find a value for the env var in the env map
		if value, ok := env[doc[hit[2]:hit[3]]]; ok {
			// we have a match; write out the value instead of the env var.
			buf.WriteString(value)
		} else if doc[hit[4]:hit[5]] == ":-" {
			// we have a fallback, use that
			buf.WriteString(doc[hit[6]:hit[7]])
		} else {
			// it was required, use the error message if present, or inject one.
			msg := strings.TrimSpace(doc[hit[6]:hit[7]])
			if len(msg) == 0 {
				msg = fmt.Sprintf("missing required env var %s", doc[hit[2]:hit[3]])
			}

			buf.WriteString("${")
			buf.WriteString(doc[hit[2]:hit[3]])
			buf.WriteString(":?")
			buf.WriteString(msg)
			buf.WriteByte('}')
		}

		lastPos = hit[1]
	}
	buf.WriteString(doc[lastPos:])

	return buf.String()
}
