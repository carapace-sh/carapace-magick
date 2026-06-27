#!/bin/bash
# Generate test images using magick itself.
# Run via: go generate ./testdata/

set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

# Small 16x16 PNG
magick -size 16x16 xc:red "$DIR/small_red.png"

# Small 16x16 grayscale PNG
magick -size 16x16 grayscale: "$DIR/small_gray.png"

# Animated GIF with 3 frames
magick -size 16x16 xc:red xc:green xc:blue "$DIR/animated.gif"

# Multi-page TIFF
magick -size 16x16 xc:red xc:green xc:blue "$DIR/multi_page.tiff"

echo "Test images generated in $DIR"
