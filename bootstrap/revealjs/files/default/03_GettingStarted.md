## Getting Started

Run uphy/reveal.js container

```console
$ docker run --name reveal.js -d \
    -v "$(pwd)/data:/reveal.js/data" \
    -p "8000:8000" \
    -p "35729:35729" \
    uphy/reveal.js:3.7.0
```

1. Open http://localhost:8000 with your browser
2. Presentation files are stored on `data` directory.  
   Edit them and make your own presentation.