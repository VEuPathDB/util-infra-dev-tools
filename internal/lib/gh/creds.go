package gh

import (
	"bufio"
	"bytes"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"vpdb-dev-tool/internal/lib/scanning"
	"vpdb-dev-tool/internal/lib/util"
	"vpdb-dev-tool/internal/lib/xbytes"
	"vpdb-dev-tool/internal/lib/xos"
)

const (
	likelyUserEnvKey  = "GITHUB_USERNAME"
	likelyTokenEnvKey = "GITHUB_TOKEN"

	userPropKey  = "gpr.user"
	tokenPropKey = "gpr.key"
)

type Credentials struct {
	Username string
	Token    string
}

func RequireCredentials() (creds Credentials) {
	creds.Username, _ = os.LookupEnv(likelyUserEnvKey)
	creds.Token, _ = os.LookupEnv(likelyTokenEnvKey)

	if len(creds.Username) == 0 || len(creds.Token) == 0 {
		logrus.Debugf("environment was missing github username or token, checking for global gradle.properties")

		if !tryGradlePropsForCreds(&creds) {
			logrus.Fatal("need GitHub credentials for this operation but none were found.\n\nPlease provide them on the environment or via a global gradle.properties file.")
		}
	}

	return
}

func tryGradlePropsForCreds(creds *Credentials) bool {
	home, err := os.UserHomeDir()

	// no $HOME env var?
	if err != nil {
		return false
	}

	propsFile := path.Join(home, ".gradle/gradle.properties")

	// no gradle props file :(
	if !util.MustReturn(xos.PathExists(propsFile)) {
		return false
	}

	file := xos.MustOpen(propsFile)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	userKeyBytes := []byte(userPropKey)
	tokenKeyBytes := []byte(tokenPropKey)

	var hasUser, hasToken bool

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		if b, pos := scanning.FirstNonWhitespace(line); pos < 0 {
			continue
		} else if b == '#' || b == '!' {
			continue
		} else {
			line = line[pos:]
		}

		if bytes.HasPrefix(line, userKeyBytes) {
			creds.Username, hasUser = eatPropsValue(line[len(userKeyBytes):])
		} else if bytes.HasPrefix(line, tokenKeyBytes) {
			creds.Token, hasToken = eatPropsValue(line[len(tokenKeyBytes):])
		}

		if hasUser && hasToken {
			return true
		}
	}

	util.Must(scanner.Err())

	return false
}

// TODO: this does not handle the case where someone has a multiline property
//       definition.
func eatPropsValue(line []byte) (string, bool) {
	if len(line) == 0 {
		return "", false
	}

	i := 0

	// the next character MUST be a divider to be valid.
	if line[i] != '=' && !xbytes.IsWhitespace(line[i]) && line[i] != ':' {
		return "", false
	}

	for len(line) > i && xbytes.IsWhitespace(line[i]) {
		i++
	}

	// it was all whitespace
	if i >= len(line) {
		return "", false
	}

	if line[i] == '=' || line[i] == ':' {
		i++

		if i >= len(line) {
			return "", false
		}

		for len(line) > i && xbytes.IsWhitespace(line[i]) {
			i++
		}

		if i >= len(line) {
			return "", false
		}
	}

	return string(line[i:]), true
}
