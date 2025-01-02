#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"

VERSION_FILE="cache/master_version"

. .env.local
logfile="cache/`date '+%Y-%m-%d'`.log"
exec > >(tee -a "$logfile") 2>&1

echo "======================== job run at $(date '+%Y-%m-%d %H:%M:%S') ==========================="

if [ -f "$VERSION_FILE" ]; then
  cur_version=$(cat "$VERSION_FILE")
else
  cur_version=""
fi

if [ $# -gt 0 ] && [ "$1" == "-f" ] || [ "$1" == "-fn" ]; then
  # force update
  ./campus --db --ab --webab --putdb --forcedb --forceab
else 
  ./campus --db --ab --webab --putdb
fi

if [ $# -gt 0 ] && [ ! "$1" == "-fn" ]; then
  echo "=== run git push ==="
  ./push_master.sh "$cur_version"
fi

echo "=== run asset upload ==="
pipenv run python3 unpack_upload.py
exit

echo ">>> run_job.sh completed."
