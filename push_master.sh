#!/bin/bash
set -euo pipefail

echo ">>> Begin push_master.sh"

cd "$(dirname "$0")"

. .env.local
echo ">>> Local environment file is loaded."

REPO_NAME="gakumasu-diff"
ARTIFACT_DIR_NAME="cache/masterYaml"
VERSION_FILE="cache/master_version"
REPO_SSH="git@github.com:vertesan/gakumasu-diff.git"

old_version=$1

if [ ! -v SSH_KEY_PATH ]; then
  echo ">>> SSH_KEY_PATH is not set, will be exiting."
  exit 201
fi

if [ ! -f "$SSH_KEY_PATH" ]; then
  echo ">>> $SSH_KEY_PATH does not exists, will be exiting."
  exit 202
fi

if [ ! -e "$VERSION_FILE" ]; then 
  echo ">>> cache/master_version does not exist, will be exiting."
  exit 203
fi

new_version=$(cat "$VERSION_FILE")

if [[ "$old_version" == "$new_version" ]]; then
  echo ">>> Nothing updated, will be exiting."
  exit 0
fi

# clone repository if does not exist
if [ ! -d "$REPO_NAME" ]; then
  echo ">>> Cloning repo from remote..."
  GIT_SSH_COMMAND="ssh -o IdentitiesOnly=yes -i ${SSH_KEY_PATH}" git clone --depth 1 "$REPO_SSH" "$REPO_NAME"
fi

# Set git configurations
git -C "$REPO_NAME" config user.name vts-server
git -C "$REPO_NAME" config user.email 169537433+vts-server@users.noreply.github.com
git -C "$REPO_NAME" config core.sshCommand "ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i $SSH_KEY_PATH -F /dev/null"

# Copy database files to repository directory
cp $ARTIFACT_DIR_NAME/*.yaml $REPO_NAME

git -C "$REPO_NAME" add .
git -C "$REPO_NAME" commit -m "$new_version"
echo ">>> Pushing to remote repository..."
git -C "$REPO_NAME" push

echo ">>> push_master.sh completed."
