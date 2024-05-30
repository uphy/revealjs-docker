revealjs := docker run --rm -v $(PWD)/data:/reveal.js/data -v $(PWD)/build:/reveal.js/build --net host ghcr.io/uphy/reveal.js:5.1.0

.PHONY: init
init:
	@mkdir -p data build

.PHONY: demo
demo: init
	@$(revealjs) init demo --overwrite

.PHONY: start
start: init
	@$(revealjs) start

.PHONY: build
build: init
	@$(revealjs) build
