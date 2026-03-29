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


# prompt the developer for version
read -p "Enter CLI version to release (e.g., v1.0.0 - INCLUDE THE "v"!!): " VERSION

if [[ -z "$VERSION" ]]; then
    echo "No version entered. Aborting."
    exit 1
fi

echo "Preparing to release CLI version $VERSION..."

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
read -p "Tag and push firmware/$VERSION? (y/n) " CONFIRM
if [ "$CONFIRM" != "y" ]; then
    echo "Aborted"
    exit 0
fi

git tag -a cli/$VERSION -m "Release cli/$VERSION"
git push origin cli/$VERSION

echo "Completed! The GitHub action should take care of the rest."