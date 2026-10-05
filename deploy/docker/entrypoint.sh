#!/bin/sh
set -eu
# Keep the clinic logo available when the data directory is a fresh volume.
if [ ! -f /mattermost/data/brand/image.png ]; then
    mkdir -p /mattermost/data/brand
    cp /mattermost/branding/logo.png /mattermost/data/brand/image.png
fi
exec "$@"
