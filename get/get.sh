#!/bin/sh
set -e

TAR_FILE="/tmp/helmreleaser.tar.gz"
RELEASES_URL="https://github.com/helmreleaser/helmreleaser/releases"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" |
    rev |
    cut -f1 -d'/'|
    rev
}

download() {
  test -z "$VERSION" && VERSION="$(last_version)"
  test -z "$VERSION" && {
    echo "Unable to get helmreleaser version." >&2
    exit 1
  }

  CLEANED_VERSION=${VERSION#?}

  ARCH=$(uname -m)
  if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
  fi

  rm -f "$TAR_FILE"
  curl -s -L -o "$TAR_FILE" \
    "$RELEASES_URL/download/$VERSION/helmreleaser_$CLEANED_VERSION""_$(uname -s | tr '[:upper:]' '[:lower:]')_$ARCH-$CLEANED_VERSION.tar.gz"
}

download
tar -xf "$TAR_FILE" -C "$TMPDIR"
"${TMPDIR}/helmreleaser" "$@"
