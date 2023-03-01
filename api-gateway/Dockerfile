FROM kong:3.1.1-alpine
USER root
COPY /kong.yaml /tmp/kong.yaml
RUN luarocks install kong-plugin-response-cache
USER kong
