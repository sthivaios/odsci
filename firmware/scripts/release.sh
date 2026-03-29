#!/bin/bash
#
# Copyright 2026 Stratos Thivaios
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
#

set -e

# read version from header
VERSION=$(grep FIRMWARE_VERSION_STR App/Inc/version.h | grep -o '".*"' | tr -d '"')

if [ -z "$VERSION" ]; then
    echo "Could not read version from version.h"
    exit 1
fi

echo "Releasing firmware v$VERSION..."

# ensure the working dir is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Working directory is not clean, commit your changes first"
    exit 1
fi

# ensure we are on the main branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$BRANCH" != "main" && "$BRANCH" != "master" ]]; then
    echo "You must be on main or master to release"
    exit 1
fi

# confirm with developer
read -p "Tag and push firmware/v$VERSION? (y/n) " CONFIRM
if [ "$CONFIRM" != "y" ]; then
    echo "Aborted"
    exit 0
fi

git tag firmware/v$VERSION
git push origin firmware/v$VERSION

echo "Completed! The GitHub action should take care of the rest."
