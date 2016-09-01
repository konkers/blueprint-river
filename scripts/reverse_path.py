#!/usr/bin/env python
#
# From https://android.googlesource.com/platform/build/soong
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

from __future__ import print_function

import os
import sys

# Find the best reverse path to reference the current directory from another
# directory. We use this to find relative paths to and from the source and build
# directories.
#
# If the directory is given as an absolute path, return an absolute path to the
# current directory.
#
# If there's a symlink involved, and the same relative path would not work if
# the symlink was replace with a regular directory, then return an absolute
# path. This handles paths like out -> /mnt/ssd/out
#
# For symlinks that can use the same relative path (out -> out.1), just return
# the relative path. That way out.1 can be renamed as long as the symlink is
# updated.
#
# For everything else, just return the relative path. That allows the source and
# output directories to be moved as long as they stay in the same position
# relative to each other.
def reverse_path(path):
    if path.startswith("/"):
        return os.path.abspath('.')

    realpath = os.path.relpath(os.path.realpath('.'), os.path.realpath(path))
    relpath = os.path.relpath('.', path)

    if realpath != relpath:
        return os.path.abspath('.')

    return relpath


if __name__ == '__main__':
    print(reverse_path(sys.argv[1]))
