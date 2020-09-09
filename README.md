# drone-DCT

Drone DCT (Docker Container Trust) plugin uses Docker-in-Docker to sign Docker images to a container registry.

## Build

Build the binaries with the following commands:

```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-dct ./cmd/drone-dct
    
docker build  --file docker/docker/Dockerfile.linux.amd64 --tag zshaik/drone-dct:0.1.6 . --no-cache
```

## Usage

```
steps:
  - name: docker_build
    image: zshaik/drone-dct:0.1.6
    pull: always
    settings:
      passphrase:
        from_secret: passphrase
      repokey:
        from_secret: repokey
      rootkey:
        from_secret: rootkey
      rootcertname: <cert-name>
      rootkeyname:
        from_secret: rootkeyname
      repo: <repo-name>
      tag: <image-tag>
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
```
