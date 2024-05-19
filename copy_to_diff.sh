#!/bin/bash
set -e

if [ -d "cache/masterYaml" ]; then
  cp cache/masterYaml/* ../gakumasu-diff/
fi
