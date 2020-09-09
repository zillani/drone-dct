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
    
docker build  --file docker/docker/Dockerfile.linux.amd64 --tag zshaik/drone-dct:0.2.0 . --no-cache
```

## Docker 

```bash
docker run --rm \
  -e PLUGIN_PASSPHRASE=root1234 \
  -e PLUGIN_REPOKEY= LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0Kcm9sZTogdGl0YW4KCk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRWQ5QkY4Qk9KdFFOOHRWdUxad1hnak05dGkraHAKREcxQk9CZit4M1F6a2lMZWlaRmk0dmRUdk5KQWdzTW9helY5dXpCVzJHYVUydUVkQjNYRjYydGgzQT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ \
  -e PLUGIN_ROOTKEY= LS0tLS1CRUdJTiBFTkNSWVBURUQgUFJJVkFURSBLRVktLS0tLQpyb2xlOiB0aXRhbgoKTUlIdU1Fa0dDU3FHU0liM0RRRUZEVEE4TUJzR0NTcUdTSWIzRFFFRkREQU9CQWpUcEM0Umhsc3ZEUUlDQ0FBdwpIUVlKWUlaSUFXVURCQUVxQkJDTWl6NUo1TXVxWGhRZFEvRHk0WXYxQklHZ1hYdEZ0UVllcnN4andpTlBMcy9YCnQ1MCtBdTVBUUJkOWtnVmQ5NXRXajMyeXBrQ2NucTNudHFJOWlJSXJkSVJwd2x4b3pMV251Tkp3V1VyclJZNE0KajR1ZnNyVEU4cHMwVmlaMnlsZnY5c1o0OVlpZzFpejEvUkhlZlJyblBhSzZvN0NOZEcrNVlGNURkYkV1SWtLawo5enNyQzBGVmE5QXlrRVNINGZ1Y2VIcGtXcDJLK2hqeWdPWTVNL0ZiTU9RK2JrbElpcFZZT09CVXlabG1aKytPCmRRPT0KLS0tLS1FTkQgRU5DUllQVEVEIFBSSVZBVEUgS0VZLS0tLS0 \
  -e PLUGIN_ROOTCERTNAME= harbortrust \
  -e PLUGIN_ROOTKEYNAME= af6776266b1b0e753875b7f1ee4c845b1e002671aa0d91b104aa1bc1f138110d \
  -e PLUGIN_REPO= zshaik/goldengoose \
  -e PLUGIN_TAG= v0.1.0 \
  -e PLUGIN_USERNAME= zshaik \
  -e PLUGIN_PASSWORD= 9247411984 \
  -e PLUGIN_TAG=latest \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  --privileged \
  zshaik/drone-dct --dry-run
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

## Enable Privileged mode on drone runner

you need to have this env set in drone runner, 

```bash
    spec:
      containers:
      - env:
        - name: DRONE_RUNNER_PRIVILEGED_IMAGES
          value: zshaik/drone-dct
```
[gitter reference](https://gitter.im/drone/drone?at=5f4580c8ec534f584fbfaf1f)