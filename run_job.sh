#!/bin/bash
set -euo pipefail

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

./campus --db 2>&1 | tee -a "$logfile"

echo "=== run git push ===" >> "$logfile"
./push_master.sh "$cur_version" 2>&1 | tee -a "$logfile"

echo ">>> run_job.sh completed."
