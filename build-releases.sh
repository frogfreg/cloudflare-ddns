#!/bin/bash

# Simple Cloudflare DDNS Release Builder
# Builds binaries for multiple platforms and generates checksums

set -e

VERSION=${1:-"0.0.1"}
BINARY_NAME="cloudflare-ddns"
RELEASE_DIR="releases"

# Supported platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)


# Create release directory
mkdir -p "$RELEASE_DIR"

echo "Building binaries for version $VERSION..."

# Build binaries
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${BINARY_NAME}-${GOOS}-${GOARCH}.exe"
    else
        OUTPUT_NAME="${BINARY_NAME}-${GOOS}-${GOARCH}"
    fi
    
    echo "Building $OUTPUT_NAME..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$RELEASE_DIR/$OUTPUT_NAME" main.go
done

echo "All binaries built successfully"

# Generate checksums
cd "$RELEASE_DIR"
echo "Generating checksums..."

for file in *; do
    if [ -f "$file" ]; then
        if command -v sha256sum >/dev/null 2>&1; then
            sha256sum "$file" > "$file.sha256"
        elif command -v shasum >/dev/null 2>&1; then
            shasum -a 256 "$file" > "$file.sha256"
        fi
    fi
done

cd ..
echo "Done. Files in $RELEASE_DIR:"
ls -la "$RELEASE_DIR" 
