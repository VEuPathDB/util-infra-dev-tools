package xyml

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func MergeNodes(first, second *yaml.Node) *yaml.Node {
	if first.Kind != second.Kind {
		return second
	}

	switch first.Kind {
	case yaml.ScalarNode:
		return second
	case yaml.MappingNode:
		return mergeMappings(first, second)
	case yaml.SequenceNode:
		return mergeSequences(first, second)
	case yaml.AliasNode:
		return second
	case yaml.DocumentNode:
		return mergeDocuments(first, second)
	default:
		logrus.Fatalf("unrecognized node kind %d", first.Kind)
		return nil // unreachable
	}
}

func mergeDocuments(first, second *yaml.Node) *yaml.Node {
	out := copyNode(second)

	if len(first.Content) == 0 {
		out.Content = second.Content
	} else if len(second.Content) == 0 {
		out.Content = first.Content
	} else {
		out.Content = []*yaml.Node{MergeNodes(first.Content[0], second.Content[0])}
	}

	return out
}

func mergeSequences(first, second *yaml.Node) *yaml.Node {
	out := copyNode(second)
	ord := make([]string, 0, max(len(first.Content), len(second.Content)))
	set := make(map[string]*yaml.Node, len(ord))

	for _, node := range first.Content {
		key := nodeKey(node)

		ord = append(ord, key)
		set[key] = node
	}
	for _, node := range second.Content {
		key := nodeKey(node)

		if _, ok := set[key]; !ok {
			ord = append(ord, key)
		}

		set[key] = node
	}

	out.Content = make([]*yaml.Node, 0, len(set))

	for _, key := range ord {
		out.Content = append(out.Content, set[key])
	}

	return out
}

func mergeMappings(first, second *yaml.Node) *yaml.Node {
	out := copyNode(second)
	ord := make([]string, 0, max(len(first.Content), len(second.Content))/2)
	set := make(map[string]struct{ f, s int }, len(ord))

	i := 0
	for i < len(first.Content) {
		key := nodeKey(first.Content[i])

		ord = append(ord, key)
		set[key] = struct{ f, s int }{f: i, s: -1}

		i += 2
	}

	i = 0
	for i < len(second.Content) {
		key := nodeKey(second.Content[i])

		if v, ok := set[key]; ok {
			v.s = i
			set[key] = v
		} else {
			ord = append(ord, key)
			set[key] = struct{ f, s int }{f: -1, s: i}
		}

		i += 2
	}

	out.Content = make([]*yaml.Node, 0, len(set)*2)

	for _, k := range ord {
		v := set[k]

		if v.s > -1 {
			if v.f > -1 {
				logrus.Tracef("merging key %s", k)
				out.Content = append(out.Content, second.Content[v.s], MergeNodes(first.Content[v.f+1], second.Content[v.s+1]))
			} else {
				logrus.Tracef("overwriting key %s", k)
				out.Content = append(out.Content, second.Content[v.s], second.Content[v.s+1])
			}
		} else {
			logrus.Tracef("keeping key %s", k)
			out.Content = append(out.Content, first.Content[v.f], first.Content[v.f+1])
		}
	}

	return out
}

func copyNode(node *yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:   node.Kind,
		Style:  node.Style,
		Tag:    node.Tag,
		Value:  node.Value,
		Anchor: node.Anchor,
		Alias:  node.Alias,
	}
}

var nodeKeyDisc uint32 = 0

func nodeKey(node *yaml.Node) string {
	if node.Kind == yaml.ScalarNode {
		return node.Value
	} else {
		nodeKeyDisc++
		return fmt.Sprintf("node-%d", nodeKeyDisc)
	}
}
