#!/bin/sh
set -eu
cd "$(dirname "$0")/../.."
repo_root="$PWD"
context="$repo_root/.local-dev/docker-arm64"
# Requires Go, Python 3, Node 24.11.1, npm 11.6.2 and installed webapp dependencies.
mkdir -p "$context/mattermost/bin" "$context/mattermost/config"
# Remove cached shared assets so development JSX cannot enter the release.
rm -rf webapp/platform/shared/dist webapp/platform/shared/.parcel-cache
(cd webapp && NODE_ENV=production NODE_OPTIONS=--max-old-space-size=6144 npm run build)
(cd server && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -tags 'sourceavailable production' \
    -ldflags '-s -w -X github.com/mattermost/mattermost/server/public/model.BuildNumber=medikids -X github.com/mattermost/mattermost/server/public/model.BuildEnterpriseReady=false -X github.com/mattermost/mattermost/server/public/model.BuildForkFeatures=true' \
    -o "$context/mattermost/bin/mattermost" ./cmd/mattermost)
rm -f "$context/mattermost/config/config.json"
(cd server && OUTPUT_CONFIG="$context/mattermost/config/config.json" go run ./scripts/config_generator)
python3 deploy/docker/prepare.py
docker build --platform linux/arm64 -f deploy/docker/Dockerfile -t "${1:-edmundog/comunicacion:latest}" "$context"
