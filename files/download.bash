#!/bin/bash

# Ensure the output directory exists
mkdir -p ./audio

# Loop over each line in output.txt
while IFS= read -r song || [ -n "$song" ]; do
    # Skip empty lines
    [ -z "$song" ] && continue

    echo "Downloading: $song"

    yt-dlp "ytsearch1:${song}" \
        --extract-audio \
        --audio-format mp3 \
        --output "./audio/${song}.%(ext)s" \
        --no-playlist

done <output.txt
