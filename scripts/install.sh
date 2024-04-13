#!/usr/bin/env sh

organization=VEuPathDB
repository=util-vpdb-dev-tool

repoUrl="https://github.com/${organization}/${repository}"
apiUrl=
releaseSuffix="/releases/latest"

checkPrereqs() {
  if ! command -v curl; then
    echo "curl is required to run the automatic install script" 1>&2
    exit 1
  fi
}

getDownloadURL() {
  curl -L -H'X-GitHub-Api-Version: 2022-11-28' https://api.github.com/repos/${organization}/${repository}/releases/latest \
    | grep 'browser_download_url' \
    | grep $1 \
    | grep -o 'http[^"]\+'
}

case $(uname | tr '[:upper:]' '[:lower:]') in
  linux*)

    ;;
  darwin*)
    ;;
  *)
    echo "unsupported os" 1>&2
    exit 1
    ;;
esac