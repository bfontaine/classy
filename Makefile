BIN=classy

SOURCES=$(wildcard **/*.go)

all: $(BIN)

$(BIN): $(SOURCES)
	go build $(BIN).go

clean:
	$(RM) $(BIN)

check:
	go test -v ./...
