#!/bin/bash

#   Copyright The containerd Authors.

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

function verify {
  echo "Verifying go.mod/go.sum in $(pwd)"
  go mod tidy
  diff=`git diff`
  if [[ -n "${diff}" ]]; then
    echo "${diff}"
    echo
    echo "error"
    exit 1
  fi
}

if ! git diff --quiet; then \
  printf "Working tree is not clean, can't proceed\n"
  exit 1
fi

# verify the root go.mod/go.sum file
verify

# verify the api go.mod/go.sum file
pushd $ROOT/api >> /dev/null
verify
popd >> /dev/null

# verify the integration client go.mod/go.sum file
pushd $ROOT/integration/client >> /dev/null
verify
popd >> /dev/null

echo "Done"
