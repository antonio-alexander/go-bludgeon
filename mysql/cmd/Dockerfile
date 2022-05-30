#---------------------------------------------------------------------------------------------------
# bludgeon-mysql [Dockerfile]
# 
# Reference: https://stackoverflow.com/questions/63178036/how-to-find-commit-hash-from-within-a-running-docker-image
# commit: git rev-parse HEAD
# 
# https://stackoverflow.com/questions/6245570/how-to-get-the-current-branch-name-in-git
# branch: git rev-parse --abbrev-ref HEAD
# 
# Sample docker build commands:
#  docker build -f ./cmd/Dockerfile . -t ghcr.io/antonio-alexander/go-bludgeon-mysql:amd64_latest \
#   --build-arg GIT_COMMIT=$GITHUB_SHA --build-arg GIT_BRANCH=$GITHUB_REF --build-arg PLATFORM=linux/amd64
#
#---------------------------------------------------------------------------------------------------

ARG GIT_BRANCH=no_git_branch_provided
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64

FROM --platform=${PLATFORM} alpine:latest

ENV MYSQL_ROOT_PASSWORD=bludgeon
ENV MYSQL_DATABASE=bludgeon
ENV MYSQL_USER=bludgeon
ENV MYSQL_PASSWORD=bludgeon

RUN \
    apk update \
    && apk add --no-cache mariadb mariadb-client \
    && mkdir -p /run/mysqld \
    && chown -R root:root /run/mysqld \
    && mkdir -p /var/lib/mysql \
    && chown -R root:root /var/lib/mysql \
    && mysql_install_db --user=root --ldata=/var/lib/mysql

RUN mkdir /bludgeon

COPY ./cmd/init_bludgeon_db.sh /bludgeon
COPY ./sql /bludgeon
COPY ./cmd/docker_entrypoint.sh /

RUN \
    chmod +x /bludgeon/init_bludgeon_db.sh \
    && chmod +x /docker_entrypoint.sh \
    && sh /bludgeon/init_bludgeon_db.sh

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

ENTRYPOINT ["/docker_entrypoint.sh"]
