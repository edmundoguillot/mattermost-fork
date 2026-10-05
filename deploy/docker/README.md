# Comunicación Medikids (ARM64)

Build the modified server and production frontend from the repository root:

```sh
./deploy/docker/build-arm64.sh edmundog/comunicacion:latest
docker push edmundog/comunicacion:latest
```

Stop the local webapp development watcher before building; it shares generated output directories with production builds. Packaging rejects shared JavaScript that contains development JSX.

Requires Go, Python 3, Node 24.11.1, npm 11.6.2, Docker and installed webapp dependencies (`cd webapp && SHARP_IGNORE_GLOBAL_LIBVIPS=1 npm ci`). Build output is kept in ignored `.local-dev/docker-arm64`; development configuration, credentials and instance data are excluded.

The image contains the clinic branding and Spanish defaults. Supply `MM_SQLSETTINGS_DATASOURCE` with your PostgreSQL connection string and `MM_SERVICESETTINGS_SITEURL` with your public URL. PostgreSQL runs separately. Persist `/mattermost/data`, `/mattermost/config`, `/mattermost/logs`, `/mattermost/plugins` and `/mattermost/client/plugins` with volumes writable by UID/GID 2000. SMTP is configured separately using Mattermost environment variables. This image is Linux ARM64 only.
