#!/bin/bash
set -eo pipefail
cd "$(dirname "$0")"

VERSION_FILE="cache/master_version"

. .env.local
logfile="cache/`date '+%Y-%m-%d'`.log"
exec > >(tee -a "$logfile") 2>&1

echo "======================== job run at $(date '+%Y-%m-%d %H:%M:%S') ==========================="

if [ -f "$VERSION_FILE" ]; then
  pre_version=$(cat "$VERSION_FILE")
else
  pre_version=""
fi

if [ $# -gt 0 ] && [ "$1" == "-f" ] || [ "$1" == "-fn" ]; then
  # force update
  ./campus --db --ab --webab --putdb --forcedb --forceab
else 
  ./campus --db --ab --webab --putdb
fi

if [ $# -gt 0 ] && [ ! "$1" == "-fn" ] || [ $# -eq 0 ]; then
  echo "=== run git push ==="
  ./push_master.sh "$pre_version"
fi

if [ -f "$VERSION_FILE" ]; then
  new_version=$(cat "$VERSION_FILE")
else
  new_version=""
fi

echo "=== run asset upload ==="
pipenv run python3 unpack_upload.py

if [[ "$pre_version" != "$new_version" ]]; then
  cd script && pnpm i && pnpm tsx main.ts
fi

echo ">>> run_job.sh completed."
