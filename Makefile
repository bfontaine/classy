BIN=classy

all: $(BIN)

%: %.go
	go build $<

clean:
	$(RM) $(BIN)
