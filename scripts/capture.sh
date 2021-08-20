#!/usr/bin/env bash
set -e

source "$(dirname "$0")/lib.sh"

echo "Capturing image to $(image_file_path)"

raspistill -w 480 -h 640 -rot 270 -o "$(image_file_path)"
