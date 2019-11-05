FROM debian:stable-slim

RUN groupadd -r nonroot && useradd -m -g nonroot -s /bin/bash nonroot

RUN apt-get update && apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/*

ARG CFSSL_RELEASE
RUN curl -Ls -o /bin/cfssl       https://github.com/cloudflare/cfssl/releases/download/v${CFSSL_RELEASE}/cfssl_${CFSSL_RELEASE}_linux_amd64       && chmod +x /bin/cfssl
RUN curl -Ls -o /bin/cfssljson   https://github.com/cloudflare/cfssl/releases/download/v${CFSSL_RELEASE}/cfssljson_${CFSSL_RELEASE}_linux_amd64   && chmod +x /bin/cfssljson
RUN curl -Ls -o /bin/multirootca https://github.com/cloudflare/cfssl/releases/download/v${CFSSL_RELEASE}/multirootca_${CFSSL_RELEASE}_linux_amd64 && chmod +x /bin/multirootca

USER nonroot
