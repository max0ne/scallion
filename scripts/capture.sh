#!/usr/bin/env bash
set -e

source "$(dirname "$0")/lib.sh"

echo "Capturing image to $(image_file_path)"

mkdir -p "$(image_dir_path)"
raspistill -o "$(image_file_path)"