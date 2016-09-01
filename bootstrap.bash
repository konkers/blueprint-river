#!/bin/bash
#
# Based on bootstrap.bash from
# https://android.googlesource.com/platform/build/soong
# Copyright 2015 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

ORIG_SRCDIR=$(dirname "${BASH_SOURCE[0]}")
if [[ "$ORIG_SRCDIR" != "." ]]; then
  if [[ ! -z "$BUILDDIR" ]]; then
    echo "error: To use BUILDDIR, run from the source directory"
    exit 1
  fi
  export BUILDDIR=$("${ORIG_SRCDIR}/build/river/scripts/reverse_path.py" "$ORIG_SRCDIR")
  cd $ORIG_SRCDIR
fi
if [[ -z "$BUILDDIR" ]]; then
  echo "error: Run ${BASH_SOURCE[0]} from the build output directory"
  exit 1
fi
export SRCDIR="."
export BOOTSTRAP="${SRCDIR}/bootstrap.bash"

export BOOTSTRAP_MANIFEST="${SRCDIR}/build/river/build.ninja.in"
export RUN_TESTS="-t"

case $(uname) in
    Linux)
	export GOOS="linux"
  export PREBUILTOS="linux-$(uname -m)"
	;;
    Darwin)
	export GOOS="darwin"
  export PREBUILTOS="darwin-$(uname -m)"
	;;
    *) echo "unknown OS:" $(uname) && exit 1;;
esac
export GOROOT="${SRCDIR}/prebuilts/go/$PREBUILTOS/"
export GOARCH="amd64"
export GOCHAR="6"

if [[ $# -eq 0 ]]; then
    mkdir -p $BUILDDIR

    if [[ $(find $BUILDDIR -maxdepth 1 -name Android.bp) ]]; then
      echo "FAILED: The build directory must not be a source directory"
      exit 1
    fi

    export SRCDIR_FROM_BUILDDIR=$(build/river/scripts/reverse_path.py "$BUILDDIR")

    sed -e "s|@@BuildDir@@|${BUILDDIR}|" \
        -e "s|@@SrcDirFromBuildDir@@|${SRCDIR_FROM_BUILDDIR}|" \
        -e "s|@@PrebuiltOS@@|${PREBUILTOS}|" \
        "$SRCDIR/build/river/river.bootstrap.in" > $BUILDDIR/.river.bootstrap
    ln -sf "${SRCDIR_FROM_BUILDDIR}/build/river/river.bash" $BUILDDIR/river
fi

"$SRCDIR/build/blueprint/bootstrap.bash" "$@"
