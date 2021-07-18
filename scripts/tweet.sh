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

# Twitter allow maximum 30 sec video
# divide and ceil IMAGE_COUNT by 30 sec to get the framerate that can fit all images into 30 sec
MAX_VIDEO_SECONDS=30
FRAMERATE=$(( (IMAGE_COUNT + MAX_VIDEO_SECONDS - 1) / MAX_VIDEO_SECONDS ))

# Encode mp4
ffmpeg \
  -framerate "$FRAMERATE" \
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

# Cleanup
rm -r "$TEMPDIR"
for IMAGE in $IMAGES; do
  rm "$IMAGE_DIR_PATH/$IMAGE"
done
