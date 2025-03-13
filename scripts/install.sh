#!/usr/bin/env bash

declare -r organization=VEuPathDB
declare -r repository=util-infra-dev-tools
declare -r osName=$(uname | tr '[:upper:]' '[:lower:]')
declare -r tmpFileName="$(mktemp)"
declare -r binDir="${HOME}/.local/bin"
declare -r incDir="${HOME}/.local/share/vpdb"

trap 'rm -f $tmpFileName' EXIT

function checkPrereqs() {
  if ! command -v curl; then
    echo "curl is required to run the automatic install script" 1>&2
    exit 1
  fi
}

function getDownloadURL() {
  curl -LsH'X-GitHub-Api-Version: 2022-11-28' "https://api.github.com/repos/${organization}/${repository}/releases/latest" \
    | grep 'browser_download_url' \
    | grep "$1" \
    | grep -o 'http[^"]\+'
}

function getPathVar() {
  grep -m1 '^export PATH=' "$1"
}

function patchEnvRC() {
  local -r currentPathValue=$(getPathVar "$1")
  local -i exitCode=0

  if ! grep -q "source ${incDir}/autocomplete.sh" "$1"; then
    echo "source ${incDir}/autocomplete.sh" >> "$1"
    exitCode=1
  fi


  if ! echo "$currentPathValue" | grep -q "\$HOME/.local/bin\|$HOME/.local/bin"; then
    echo "export PATH=\"\${PATH}:$HOME/.local/bin\"" >> "$1"
    exitCode=1
  fi

  return $exitCode
}

declare printedSource=0

function printSourceHelp() {
  echo "To add the vpdb command to your current PATH execute:"
  echo "  source $1"
  printedSource=1
}


case $osName in
  linux|darwin)
    # do nothing
    ;;
  *)
    echo "unsupported os version" 1>&2
    exit 1
    ;;
esac


curl -Lso "${tmpFileName}" "$(getDownloadURL "$osName")"
mkdir -p "${binDir}"
unzip -qq -o "${tmpFileName}" -d "${binDir}"
mkdir -p "${incDir}"
mv "${binDir}/autocomplete.sh" "${incDir}"

if command -v vpdb >/dev/null; then
  echo "Installed vpdb-dev-tool:"
  echo ""
  vpdb -v
  exit
fi

if [ -f "$HOME/.zshrc" ]; then
  if ! patchEnvRC "$HOME/.zshrc"; then
    printSourceHelp "$HOME/.zshrc"
  fi
fi

if [ -f "$HOME/.bashrc" ]; then
  if ! patchEnvRC "$HOME/.bashrc" && [ $printedSource -eq 0 ]; then
    printSourceHelp "$HOME/.bashrc"
  fi
fi

if command -v vpdb >/dev/null; then
  echo "Installed vpdb-dev-tool:"
  echo ""
  vpdb -v
else
  echo "please add ${binDir}/vpdb to your path to get started"
fi