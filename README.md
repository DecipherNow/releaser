# Releaser
Implement the decipher release process.


## Usage
### Help
```bash
NAME:
   Releaser - Facilitate the release process of artifacts

USAGE:
   releaser [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     docker   Do the docker job
     github   Do the github release
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --symver value          Symantic Version of the release to prepare
   --dockerImage value     Source Docker image to release
   --dockerUser value      Username for docker push/pulls
   --dockerPassword value  Password for docker push/pulls
   --githubToken value     Access token for github releases
   --githubUser value      Username for github releases
   --githubPassword value  Password for github releases
   --asset value           File[s] to be uploaded to the github release
   --help, -h              show help
   --version, -v           print the version

```

### Symver tagging/push docker images
`./releaser docker --symver $VERSION --dockerImage deciphernow/gm-proxy --dockerUser $DOCKER_USER --dockerPassword $DOCKER_PASS`

### Create release and upload an asset for the release
``./releaser github --symver $VERSION --githubToken $GITHUB_TOKEN --githubOrg deciphernow --asset ./binary-asset`

## Build
- `dep ensure -v`
- `go build`