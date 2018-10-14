# Reveal.js Docker

Simplified Reveal.js server

- All you need is Markdown slide files  
  `index.html` is not required
- YAML formatted config

## Getting Started

Run uphy/reveal.js container

```console
$ docker run --name reveal.js -d \
    -v "$(pwd)/data:/reveal.js/data" \
    -p "8000:8000" \
    -p "35729:35729" \
    uphy/reveal.js:3.7.0
```

Open http://localhost:8000 with your browser

Presentation files are stored on `data` directory.  
Edit them and make your own presentation.

## Config File

You can change the presentation theme, title, and many [reveal.js configs](https://github.com/hakimel/reveal.js/#configuration)  
Config file(`config.yml`) is located on your `data` directory. 

## Generate demo data

Generate demo data.

```console
$ docker run --name reveal.js -d \
    -v "$(pwd)/data:/reveal.js/data" \
    -p "8000:8000" \
    -p "35729:35729" \
    uphy/reveal.js:3.7.0 init demo
```

Start server

```console
$ docker run --name reveal.js -d \
    -v "$(pwd)/data:/reveal.js/data" \
    -p "8000:8000" \
    -p "35729:35729" \
    uphy/reveal.js:3.7.0
```

Build presentation as static html files.  

```console
$ docker run --rm \
    -v "$(pwd)/data:/reveal.js/data" \
    -v "$(pwd)/doc:/reveal.js/build" \
    uphy/reveal.js:3.7.0 build
```

Files are stored in `doc` directory.  
This is useful for hosting the presentation on GitHub Pages.

## docker-compose

Create docker-compose.yml.

```yaml
version: "2"

services:
  revealjs:
    image: uphy/reveal.js:3.7.0
    ports:
      - "8000:8000"
      - "35729:35729"
    volumes:
      - "./data:/reveal.js/data"
      - "./doc:/reveal.js/build"
```

Generate demo data.

```console
$ docker-compose run --rm revealjs init demo
```

Start server.

```console
$ docker-compose up -d
```

Build presentation as static html files.  

```console
$ docker-compose run --rm revealjs build
```

## Advanced

If you want to edit `index.html` directly, edit `index.html.tmpl` located on your `data` directory.