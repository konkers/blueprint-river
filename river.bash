#!/bin/bash
#
# Based on song.bash from
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

# Switch to the build directory
cd $(dirname "${BASH_SOURCE[0]}")

# The source directory path and operating system will get written to
# .river.bootstrap by the bootstrap script.

BOOTSTRAP=".river.bootstrap"
if [ ! -f "${BOOTSTRAP}" ]; then
    echo "Error: river script must be located in a directory created by bootstrap.bash"
    exit 1
fi

source "${BOOTSTRAP}"

# Now switch to the source directory so that all the relative paths from
# $BOOTSTRAP are correct
cd ${SRCDIR_FROM_BUILDDIR}

export GOROOT="prebuilts/go/$PREBUILTOS/"

# Run the blueprint wrapper
NINJA="./prebuilts/ninja/${PREBUILTOS}/ninja" \
    BUILDDIR="${BUILDDIR}" SKIP_NINJA=true build/blueprint/blueprint.bash

# Ninja can't depend on environment variables, so do a manual comparison
# of the relevant environment variables from the last build using the
# river_env tool and trigger a build manifest regeneration if necessary
ENVFILE="${BUILDDIR}/.river.environment"
ENVTOOL="${BUILDDIR}/.bootstrap/bin/river_env"
if [ -f "${ENVFILE}" ]; then
    if [ -x "${ENVTOOL}" ]; then
        if ! "${ENVTOOL}" "${ENVFILE}"; then
            echo "forcing build manifest regeneration"
            rm -f "${ENVFILE}"
        fi
    else
        echo "Missing river_env tool, forcing build manifest regeneration"
        rm -f "${ENVFILE}"
    fi
fi

"prebuilts/ninja/${PREBUILTOS}/ninja" -f "${BUILDDIR}/build.ninja" -w dupbuild=err "$@"
