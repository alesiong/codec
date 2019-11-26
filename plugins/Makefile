ALL = ${shell find . -type d}

.PHONY: all
all: $(ALL)

%: %.so
	go build -buildmode=plugin $@/*