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

## Advanced

If you want to edit `index.html` directly, edit `index.html.tmpl` located on your `data` directory.