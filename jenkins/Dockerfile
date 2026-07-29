FROM jenkins/jenkins:lts-jdk17

USER root

RUN apt-get update \
  && apt-get install -y --no-install-recommends git ca-certificates curl tar \
  && curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-27.3.1.tgz \
    | tar -xz -C /usr/local/bin --strip-components=1 docker/docker \
  && mkdir -p /usr/local/lib/docker/cli-plugins \
  && curl -fsSL https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-linux-x86_64 \
    -o /usr/local/lib/docker/cli-plugins/docker-compose \
  && chmod +x /usr/local/bin/docker /usr/local/lib/docker/cli-plugins/docker-compose \
  && rm -rf /var/lib/apt/lists/*

USER jenkins

RUN jenkins-plugin-cli --plugins workflow-aggregator git docker-workflow
