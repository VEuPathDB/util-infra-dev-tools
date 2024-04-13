#!/usr/bin/env sh

organization=VEuPathDB
repository=util-vpdb-dev-tool

checkPrereqs() {
  if ! command -v curl; then
    echo "curl is required to run the automatic install script" 1>&2
    exit 1
  fi
}

getDownloadURL() {
  curl -LsH'X-GitHub-Api-Version: 2022-11-28' https://api.github.com/repos/${organization}/${repository}/releases/latest \
    | grep 'browser_download_url' \
    | grep $1 \
    | grep -o 'http[^"]\+'
}

getPathVar() {
  grep -m1 '^export PATH=' "$1"
}

patchEnvRC() {
  currentPathValue=$(getPathVar "$1")
  if ! echo "$currentPathValue" | grep -q "\$HOME/.local/bin\|$HOME/.local/bin"; then
    echo "export PATH=\"\${PATH}:$HOME/.local/bin\"" >> "$1"
    echo "To add the vpdb command to your current PATH execute:"
    echo "  source $1"
  fi
}

osName=$(uname | tr '[:upper:]' '[:lower:]')

case $osName in
  linux|darwin)
    # do nothing
    ;;
  *)
    echo "unsupported os version" 1>&2
    exit 1
    ;;
esac

tmpFileName=/tmp/${repository}
binDir="${HOME}/.local/bin"

curl -Lso ${tmpFileName} "$(getDownloadURL "$osName")"
mkdir -p "${binDir}"
unzip -qq ${tmpFileName} -d "${binDir}"

trap 'rm -f $tmpFileName' EXIT

if command -v vpdb; then
  exit
fi

if [ -f "$HOME/.zshrc" ]; then
  patchEnvRC "$HOME/.zshrc"
fi

if [ -f "$HOME/.bashrc" ]; then
  patchEnvRC "$HOME/.bashrc"
fi
