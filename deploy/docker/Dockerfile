FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl tzdata mime-support poppler-utils wv unrtf tidy \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 2000 mattermost \
    && useradd --uid 2000 --gid 2000 --home-dir /mattermost mattermost
COPY --chown=2000:2000 mattermost /mattermost
COPY --chmod=755 entrypoint.sh /usr/local/bin/comunicacion-entrypoint
RUN mkdir -p /mattermost/data /mattermost/logs /mattermost/plugins /mattermost/client/plugins \
    && chown -R 2000:2000 /mattermost
WORKDIR /mattermost
ENV MM_INSTALL_TYPE=docker
USER 2000:2000
EXPOSE 8065
VOLUME ["/mattermost/data", "/mattermost/logs", "/mattermost/config", "/mattermost/plugins", "/mattermost/client/plugins"]
HEALTHCHECK --interval=15s --timeout=5s --start-period=60s CMD curl -fsS http://localhost:8065/api/v4/system/ping || exit 1
ENTRYPOINT ["comunicacion-entrypoint"]
CMD ["/mattermost/bin/mattermost", "server", "--config", "/mattermost/config/config.json"]
