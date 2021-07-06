#!/usr/bin/env bash
set -e

source "$(dirname "$0")/lib.sh"

IMAGE_DIR_PATH=$(image_dir_path)
IMAGE_COUNT=$(ls -1 "$IMAGE_DIR_PATH" | wc -l)
GIF_PATH="$IMAGE_DIR_PATH/out.gif"

if [ "$IMAGE_COUNT" -eq "0" ]; then
  echo "No images from in $IMAGE_DIR_PATH"
fi

echo "Creating GIF with $IMAGE_COUNT images from $IMAGE_DIR_PATH"
# delay's unit is 1/100 of a second
convert -resize 20% -delay 10 -loop 0 `ls -v1tr "$IMAGE_DIR_PATH"/*.jpg` "$GIF_PATH"

echo "Tweeting image capture"
TWEET_URL=$(scallion tweet -i "$GIF_PATH" -c /etc/scallion/cred.json)

echo "Tweet sent at $TWEET_URL"
