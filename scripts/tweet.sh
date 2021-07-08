#!/usr/bin/env bash
set -e

source "$(dirname "$0")/lib.sh"

IMAGE_DIR_PATH=$(image_dir_path)

IMAGES=$(ls -v1tr "$IMAGE_DIR_PATH")
IMAGE_COUNT=$(echo "$IMAGES" | wc -l)

if [ ! "$(ls $IMAGE_DIR_PATH)" ]; then
  echo "No images from in $IMAGE_DIR_PATH"
  exit
fi

echo "Found $IMAGE_COUNT images in $IMAGE_DIR_PATH"

TEMPDIR="$(mktemp -d)"
echo "Using tempdir $TEMPDIR"

# Copy images to TEMPDIR
IDX=0
for IMAGE in $IMAGES; do
  cp "$IMAGE_DIR_PATH/$IMAGE" "$TEMPDIR/$IDX.jpg"
  IDX=$(( IDX + 1 ))
done

# Encode mp4
ffmpeg \
  -framerate 10 \
  -i "file:$TEMPDIR/%d.jpg" \
  -c:v libx264 \
  -profile:v high \
  -crf 20 \
  -pix_fmt yuv420p \
  -vf scale=640:-1 \
  "$TEMPDIR/output.mp4"

# Tweet it
echo "Tweeting image capture"
TWEET_URL=$(scallion tweet -i "$TEMPDIR/output.mp4" -c /etc/scallion/cred.json)

echo "Tweet sent at $TWEET_URL"

# Cleanup remaining images
for IMAGE in $IMAGES; do
  rm "$IMAGE"
done
rm -r "$TEMPDIR"
