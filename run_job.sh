#!/bin/bash
set -euo pipefail
. .env.local

VERSION_FILE="cache/master_version"

echo ">>> Begin run_job.sh"
cd "$(dirname "$0")"
logfile="cache/`date '+%Y-%m-%d'`.log"

echo "======================== job run at $(date '+%Y-%m-%d %H:%M:%S') ===========================" >> $logfile

if [ -f "$VERSION_FILE" ]; then
  cur_version=$(cat "$VERSION_FILE")
else
  cur_version=""
fi

./campus --db --ab --webab --putdb 2>&1 | tee -a "$logfile"

echo "=== run git push ===" >> "$logfile"
./push_master.sh "$cur_version" 2>&1 | tee -a "$logfile"

echo "=== run asset upload ===" >> "$logfile"
pipenv run python3 unpack_upload.py 2>&1 | tee -a "$logfile"
exit

echo ">>> run_job.sh completed."
