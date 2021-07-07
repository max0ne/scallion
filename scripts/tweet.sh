#!/usr/bin/env bash
set -e

source "$(dirname "$0")/lib.sh"

IMAGE_DIR_PATH=$(image_dir_path)

# clean up old _small images if there's any
rm -f "$IMAGE_DIR_PATH"/*_small.jpg

IMAGE_COUNT=$(ls -1 "$IMAGE_DIR_PATH" | wc -l)
GIF_PATH="$IMAGE_DIR_PATH/out.gif"

if [ "$IMAGE_COUNT" -eq "0" ]; then
  echo "No images from in $IMAGE_DIR_PATH"
fi

echo "Found $IMAGE_COUNT images in $IMAGE_DIR_PATH"

echo "Shrinking images in $IMAGE_DIR_PATH"
for IMAGE in $(ls "$IMAGE_DIR_PATH"/*.jpg); do
  echo "Shrinking image $IMAGE"
  FILE_NAME_NO_EXT=$(basename $IMAGE .jpg)
  SMALL_IMAGE="$FILE_NAME_NO_EXT"_small.jpg

  convert -resize 16% "$IMAGE" "$IMAGE_DIR_PATH"/"$SMALL_IMAGE"
done

echo "Creating GIF with $IMAGE_COUNT images from $IMAGE_DIR_PATH"

# delay's unit is 1/100 of a second
convert -delay 10 -loop 0 $(ls -v1tr "$IMAGE_DIR_PATH"/*_small.jpg) "$GIF_PATH"

echo "Tweeting image capture"
TWEET_URL=$(scallion tweet -i "$GIF_PATH" -c /etc/scallion/cred.json)

echo "Tweet sent at $TWEET_URL"
