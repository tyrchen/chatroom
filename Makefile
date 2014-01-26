CHECK=\033[32m✔\033[39m
DONE="\n$(CHECK) Done.\n"

GO=$(shell which go)
CLOC=$(shell which cloc)
BIN=./bin
TARGETS=$(patsubst %.go,$(BIN)/%,$(wildcard *.go))

build: $(TARGETS)
	@echo $(DONE)

$(TARGETS): $(BIN)/%: %.go
	@echo "building $<..."
	@$(GO) build -o $@ $<

clean:
	@$(RM) -f $(TARGETS)
	@echo $(DONE)

cloc:
	@$(CLOC) . --exclude-dir=webclient/assets
