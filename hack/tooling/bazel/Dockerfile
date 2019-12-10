FROM ubuntu:18.04

RUN apt-get update && apt-get install -y \
    # Generic requirements
    curl bash git \
    # Bazel requirements
    pkg-config zip g++ zlib1g-dev unzip python3 \
    # Docker requirements
    apt-transport-https ca-certificates gnupg-agent software-properties-common \
    # Cleanup
    && rm -rf /var/lib/apt/lists/*

# Docker CE
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io \
    && rm -rf /var/lib/apt/lists/*

# Bazel
ARG BAZEL_RELEASE
RUN curl -Lso /usr/local/bazel-installer.sh \
    https://github.com/bazelbuild/bazel/releases/download/${BAZEL_RELEASE}/bazel-${BAZEL_RELEASE}-installer-linux-x86_64.sh  \
    && chmod +x /usr/local/bazel-installer.sh \
    && /usr/local/bazel-installer.sh \
    && rm /usr/local/bazel-installer.sh \
    && ln -s /usr/bin/python3 /usr/local/bin/python

# Skaffold
ARG SKAFFOLD_RELEASE
RUN curl -Lso /usr/local/bin/skaffold https://storage.googleapis.com/skaffold/releases/${SKAFFOLD_RELEASE}/skaffold-linux-amd64 \
    && chmod +x /usr/local/bin/skaffold

# Setup bazel user
RUN useradd -u 2000 -m -g docker -s /bin/bash bazel
USER bazel

WORKDIR /workspace

ENTRYPOINT [ "/bin/bash" ]