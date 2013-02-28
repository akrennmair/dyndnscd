TARGET=dyndnscd

all: $(TARGET)

$(TARGET): $(wildcard *.go)
	go build -o $@

clean:
	$(RM) $(TARGET)

.PHONY: clean
