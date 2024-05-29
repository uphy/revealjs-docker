# How to develop

## Initialization

Fetch reveal.js and install npm packages.

```sh
$ make init
```

## Development

There are 2 ways to develop.

- Local
- Docker

### Local

#### Generate demo data

Generate demo data for reveal.js.

```sh
$ make demo
```

#### Start server

Start reveal.js server with hot reload.

```sh
$ make start
```

#### Build presentation directory

Extract the presentation data as a directory.

```sh
$ make build
```

### Docker

#### Generate demo data

```sh
$ make demo-docker
```

#### Start server

```sh
$ make start-docker
```

#### Build presentation directory

```sh
$ make build-docker
```

## Clean

```sh
$ make clean
```

## Upgrade reveal.js

- Update submodule revision of `bootstrap/reveal.js`
- Update demo slides `bootstrap/revealjs/files/demo` based on `bootstrap/reveal.js/demo.html`
- Execute `scripts/update-version.sh x.y.z`
- Push image to Docker Hub
