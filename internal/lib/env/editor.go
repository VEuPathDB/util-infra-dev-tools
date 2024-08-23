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

type Editor interface {
	CommentOutOriginalValueOnChange() Editor

	AddOrReplace(key, value string) Editor

	AddIfAbsent(key, value string) Editor

	Remove(key string) Editor

	RemoveByRegex(regex *regexp.Regexp) Editor

	RemoveMatching(predicate KeyPredicate) Editor

	ApplyEdits(input io.Reader, output io.Writer)
}

func NewEditor() Editor {
	return &envEditor{
		hardSets: make(map[string]string, 8),
		softSets: make(map[string]string, 4),
	}
}

// //

type KeyPredicate interface {
	Matches(key string) bool
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
	e.hardSets[key] = value
	return e
}

func (e *envEditor) AddIfAbsent(key, value string) Editor {
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
