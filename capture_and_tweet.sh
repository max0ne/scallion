#!/usr/bin/env bash
set -e

IMAGE_FILE="/tmp/$(date --rfc-3339=seconds).jpg"

echo "Capturing image to $IMAGE_FILE"
raspistill -o "$IMAGE_FILE"

echo "Rendering image capture with text"
MEMED_IMAGE=$(meme -i "$IMAGE_FILE" -t "|$(date "+%Y-%m-%d %H:%M")")

echo "Tweeting image capture"
TWEET_URL=$(scallion tweet -i "$MEMED_IMAGE" -c /etc/scallion/cred.json)

echo "Tweet sent at $TWEET_URL"
