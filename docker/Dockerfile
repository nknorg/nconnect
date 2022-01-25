ARG base
FROM ${base}debian:stretch-slim
ARG build_dir
ADD $build_dir /nConnect/
WORKDIR /nConnect/data/
ENTRYPOINT ["/nConnect/nConnect", "--web-root-path", "/nConnect/web/dist"]
