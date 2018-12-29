B=\033[0;1m
G=\033[0;92m
R=\033[0m

NAME ?= homo

.PHONY: help attach auto up down
# Show this help prompt
help:
	@echo '  Usage:'
	@echo ''
	@echo '    make <target>'
	@echo ''
	@echo '  Targets:'
	@echo ''
	@awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | grep -v 'IGNORE' | sort | uniq

# Auto Tillich-Zémor hasher demo
auto: down
	@echo "\n${B}${G}build container${R}\n"
	@time docker build -t poc-demo .
	@echo "\n${B}${G}Bootup container:${R}\n"
	@time docker run -d --rm -it --name hash-demo poc-demo:latest sh
	@bash ./auto.sh

# Stop demo container
down:
	@echo "\n${B}${G}Stop container${R}\n"
	@docker stop hash-demo || true
	@docker rm hash-demo || true

# Run Tillich-Zémor hasher demo
up: down
	@echo "\n${B}${G}build container${R}\n"
	@time docker build -t poc-demo .
	@echo "\n${B}${G}enter inside container:${R}\n"
	@time docker run --rm -it --name hash-demo poc-demo:latest sh

# Attach to existing container
attach:
	@echo "\n${B}${G} attach to hash-container ${R}\n"
	@time docker exec -it --name hash-demo /bin/sh
