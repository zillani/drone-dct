# drone-docker

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-docker/status.svg)](http://cloud.drone.io/drone-plugins/drone-docker)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/docker.svg)](https://microbadger.com/images/plugins/docker "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-docker?status.svg)](http://godoc.org/github.com/drone-plugins/drone-docker)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-docker)](https://goreportcard.com/report/github.com/drone-plugins/drone-docker)

Drone DCT plugin uses Docker-in-Docker to sign Docker images to a container registry.

## Build

Build the binaries with the following commands:

```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-dct ./cmd/drone-dct
    
docker build  --file docker/docker/Dockerfile.linux.amd64 --tag zshaik/drone-dct:0.1.6 .
```

## Usage

```
steps:
  - name: docker_build
    image: zshaik/drone-dct:0.1.1
    pull: always
    settings:
      passphrase:
        from_secret: passphrase
      repokey:
        from_secret: repokey
      rootkey:
        from_secret: rootkey
      rootcertname: harbortrust
      rootkeyname:
        from_secret: rootkeyname
      repo: zshaik/busybox
      tag: signed-tag
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
```
