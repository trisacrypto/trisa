FROM alpine:latest

RUN apk --no-cache add curl
RUN addgroup -g 1001 -S nonroot && adduser -u 1001 -S nonroot -G nonroot

ARG GOHUGO_RELEASE
RUN curl -Ls https://github.com/gohugoio/hugo/releases/download/v${GOHUGO_RELEASE}/hugo_${GOHUGO_RELEASE}_Linux-64bit.tar.gz | tar xvz -C /bin hugo

USER nonroot

ENTRYPOINT [ "hugo" ]