FROM golang:1.10-alpine
COPY /src /workspace/
RUN \
  cd /workspace &&\
  go build -o clawer

FROM ljnelson/docker-calibre-alpine
ENTRYPOINT ["clawer"]
COPY --from=0 /workspace/clawer /usr/local/bin/
