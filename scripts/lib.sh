#!/usr/bin/env bash
set -e

function image_dir_path {
  echo "/var/lib/scallion"
}

function image_file_name {
  echo "$(date "+%H:%M:%S").jpg"
}

function image_file_path {
  echo "$(image_dir_path)/$(image_file_name)"
}
