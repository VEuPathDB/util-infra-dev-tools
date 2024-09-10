package env

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/scanning"
	"vpdb-dev-tool/internal/lib/xio"
)

// Editor defines an API for modifying a stream of environment variable
// definitions via specific actions.
//
// Edits are applied in the following order, with no more than a single edit
// being applied to a single environment variable:
//
// 1. replace
// 2. remove
// 3. add
type Editor interface {

	// CommentOutOriginalValueOnChange configures the Editor instance to retain
	// the original values when making a change to a stream by writing out the
	// original line as a comment.
	CommentOutOriginalValueOnChange() Editor

	// AddOrReplace configures the Editor instance to set the given key to the
	// given value overwriting existing environment variable definitions on
	// conflict.
	AddOrReplace(key, value string) Editor

	// AddIfAbsent configures the Editor instance to add the given key/value pair
	// to the stream of environment variables only if doing so would not conflict
	// with an existing environment variable definition.
	AddIfAbsent(key, value string) Editor

	// Remove configures the Editor instance to remove a target environment
	// variable from the stream if it is encountered.
	Remove(key string) Editor

	// RemoveByRegex configures the Editor instance to remove any environment
	// variables from the stream whose keys match a given regular expression if
	// encountered.
	RemoveByRegex(regex *regexp.Regexp) Editor

	// RemoveMatching configures the Editor instance to remove any environment
	// variables from the stream whose keys match a given KeyPredicate if
	// encountered.
	RemoveMatching(predicate KeyPredicate) Editor

	// ApplyEdits applies all configured environment modifications to the given
	// input stream, writing the modified output to the given output stream.
	ApplyEdits(input io.Reader, output io.Writer)
}

func NewEditor() Editor {
	return &envEditor{
		hardSets: make(map[string]string, 8),
		softSets: make(map[string]string, 4),
	}
}

// //

// KeyPredicate defines an API for testing environment variable key strings.
type KeyPredicate interface {
	Matches(key string) bool
}

// KeyPredicateFunc provides a wrapper allowing for arbitrary functions to
// implement the KeyPredicate interface, provided they have a matching
// signature.
type KeyPredicateFunc func(key string) bool

func (k KeyPredicateFunc) Matches(key string) bool {
	return k(key)
}

type exactPredicate string

func (e exactPredicate) Matches(key string) bool {
	return key == string(e)
}

type regexPredicate struct{ regex *regexp.Regexp }

func (r regexPredicate) Matches(key string) bool {
	return r.regex.MatchString(key)
}

// //

type envEditor struct {
	keepOriginal bool
	hardSets     map[string]string
	softSets     map[string]string
	removals     []KeyPredicate
}

func (e *envEditor) CommentOutOriginalValueOnChange() Editor {
	e.keepOriginal = true
	return e
}

func (e *envEditor) AddOrReplace(key, value string) Editor {
	if _, ok := e.softSets[key]; ok {
		logrus.Fatalf("configured env.Editor to both AddOrReplace and AddIfAbsent environment variable \"%s\"", key)
	}
	e.hardSets[key] = value
	return e
}

func (e *envEditor) AddIfAbsent(key, value string) Editor {
	if _, ok := e.hardSets[key]; ok {
		logrus.Fatalf("configured env.Editor to both AddOrReplace and AddIfAbsent environment variable \"%s\"", key)
	}
	e.softSets[key] = value
	return e
}

func (e *envEditor) Remove(key string) Editor {
	e.removals = append(e.removals, exactPredicate(key))
	return e
}

func (e *envEditor) RemoveByRegex(regex *regexp.Regexp) Editor {
	e.removals = append(e.removals, regexPredicate{regex})
	return e
}

func (e *envEditor) RemoveMatching(predicate KeyPredicate) Editor {
	e.removals = append(e.removals, predicate)
	return e
}

func (e *envEditor) ApplyEdits(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	buffOut := xio.NewBufferedReqWriter(output)

	var key, value string
	var padding []byte

	seenKeys := make(map[string]bool, 256)

SCANNING:
	for scanner.Scan() {
		rawLine := scanner.Bytes()

		// If the line is blank,
		if len(rawLine) == 0 {
			// write a blank line to the output stream and continue to the next input
			// line.
			buffOut.WriteLineFeed()
			continue
		}

		// If the first non-space character is an octothorpe,
		if b, pos := scanning.FirstNonWhitespace(rawLine); b == '#' || pos < 0 {
			// write out the comment line and continue to the next input line.
			buffOut.Write(rawLine)
			buffOut.WriteLineFeed()
			continue
		} else

		// else, if the first non-space character is not the first character of the
		// line,
		if pos > 0 {
			// record the padding to keep the user's formatting.
			padding = rawLine[:pos]
			rawLine = rawLine[pos:]
		} else

		// else, the first non-space character is at the start of the line,
		{
			// and there is no padding.
			padding = nil
		}

		// Now we know we have a non-comment line.

		// If the line does not contain an equals sign in a valid position,
		if split := bytes.IndexByte(rawLine, '='); split < 1 {
			// then write out the invalid line and continue to the next input line.
			buffOut.Write(padding)
			buffOut.Write(rawLine)
			buffOut.WriteLineFeed()
			continue
		} else

		// else, break the line into a key and value pair using the position of the
		// equals sign.
		{
			key = string(rawLine[:split])
			value = string(rawLine[split+1:])
		}

		// Record that we have seen this key (so we can soft-set missing keys
		// later).
		seenKeys[key] = true

		// If the key appears in the map of replacement values,
		if alt, ok := e.hardSets[key]; ok {
			logrus.Debugf("replacing env key %s", key)

			// comment out the original value if requested,
			e.maybeCommentOutOriginal(buffOut, padding, rawLine)

			// write out a new line using the new value for the target key,
			writeOutPair(key, alt, buffOut)

			// then continue to the next input line.
			continue
		}

		// Check if the key is one marked for removal.
		for i := range e.removals {
			if e.removals[i].Matches(key) {
				logrus.Debugf("removing env key %s", key)

				e.maybeCommentOutOriginal(buffOut, padding, rawLine)
				continue SCANNING
			}
		}

		// Else, it's not for us and we can disregard it.
		writeOutPair(key, value, buffOut)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatalf("encountered an error while scanning the input stream: %s", err)
		panic(err) // unreachable
	}

	// Now that we've scanned through the existing keys, we can apply any softSet
	// values.

	injectedLine := false

	// Iterate through all the hard-set key/value pairs to find any we didn't
	// replace while scanning the input file.
	for key, value = range e.hardSets {
		// If the key is in the seen keys map then we replaced it previously and we
		// don't need to append it.
		if seenKeys[key] {
			continue
		}

		// insert a blank line between the original input and the new stuff (if
		// not already inserted)
		if !injectedLine {
			buffOut.WriteLineFeed()
			injectedLine = true
		}

		logrus.Debugf("appending env key %s", key)

		// write out the new key/value pair
		writeOutPair(key, value, buffOut)
	}

	// Iterate through all the soft-set key/value pairs
	for key, value = range e.softSets {
		// If the soft-set key was not seen in the input,
		if !seenKeys[key] {
			// insert a blank line between the original input and the new stuff (if
			// not already inserted)
			if !injectedLine {
				buffOut.WriteLineFeed()
				injectedLine = true
			}

			logrus.Debugf("appending env key %s", key)

			// and write out the new key/value pair
			writeOutPair(key, value, buffOut)
		}
	}

	// Flush the buffered writer.
	buffOut.Flush()
}

func (e *envEditor) maybeCommentOutOriginal(buf xio.BufferedReqWriter, padding, original []byte) {
	if e.keepOriginal {
		buf.WriteByte('#')
		buf.Write(padding)
		buf.Write(original)
		buf.WriteLineFeed()
	}
}

func writeOutPair(key, value string, writer xio.BufferedReqWriter) {
	writer.WriteString(key)
	writer.WriteByte('=')
	writer.WriteString(value)
	writer.WriteLineFeed()
}
