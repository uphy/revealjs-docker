.PHONY: init
init:
	@git submodule update && cd bootstrap/reveal.js && npm install

.PHONY: demo
demo:
	@cd bootstrap && go run main.go --dir reveal.js init demo --overwrite

.PHONY: start
start:
	@cd bootstrap && go run main.go --dir reveal.js start

.PHONY: build
build:
	@cd bootstrap && go run main.go --dir reveal.js build
	@echo See bootstrap/reveal.js/build

.PHONY: clean
clean:
	@git submodule update
	@cd bootstrap/reveal.js && git clean -dfx && git checkout HEAD . && git status
	@rm -rf data docs

.PHONY: init-docker
init-docker:
	@mkdir -p data && mkdir -p docs

.PHONY: demo-docker
demo-docker: init-docker
	@docker compose build
	@docker compose run --rm revealjs init demo --overwrite

.PHONY: start-docker
start-docker: init-docker
	@docker compose up

.PHONY: build-docker
build-docker: init-docker
	@docker compose run --rm revealjs build
