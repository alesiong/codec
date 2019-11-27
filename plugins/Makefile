ALL = tr.so charset.so

.PHONY: all
all: $(ALL)

%.so: %
	go build -buildmode=plugin $</*