#!/bin/bash

TOKEN=$(cat .token)

set -e

go build

MESSAGE=$(cat << END
a test message

enhancements for OCTO-12, CAT-567
END
)

MATTERMOST_URL=https://mattermost.brankas.dev/ \
MATTERMOST_TOKEN="$TOKEN" \
MATTERMOST_TEAM=dev \
MATTERMOST_CHANNEL=town-square \
MATTERMOST_REPLACE='[{"regexp": "(OCTO-[0-9]+)", "replace": "[$1](https://github.atlassian.net/browse/$1)"},{"regexp": "(CAT-[0-9]+)", "replace": "[$1](https://github.atlassian.net/browse/$1)"}]' \
DRONE_REPO_OWNER=octocat \
DRONE_REPO_NAME=hello-world \
DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
DRONE_COMMIT_BRANCH=master \
DRONE_COMMIT_AUTHOR=octocat \
DRONE_COMMIT_MESSAGE="$MESSAGE" \
DRONE_COMMIT_AUTHOR_EMAIL=octocat@github.com \
DRONE_COMMIT_AUTHOR_AVATAR="https://avatars0.githubusercontent.com/u/583231?s=460&v=4" \
DRONE_COMMIT_AUTHOR_NAME="The Octocat" \
DRONE_BUILD_NUMBER=1 \
DRONE_BUILD_STATUS=success \
DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
DRONE_TAG=1.0.0 \
./drone-mattermost
