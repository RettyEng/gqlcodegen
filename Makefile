BUILD_DIR:=build
STRINGER_SRC:=$(shell grep -lr '//go:generate stringer' | grep '\.go$$')
STRING_GO:=$(patsubst %.go, %_string.go, $(STRINGER_SRC))
SRC:=$(shell find . -name '*.go')
.PHONY: all clean delete-generated update-sample

all: $(BUILD_DIR)/gqlcodegen

clean:
	rm $(STRING_GO)
	rm -rf $(BUILD_DIR)

delete-generated:
	find sample/enum -mindepth 1 -type d | xargs --no-run-if-empty rm -rf
	find sample -name '*_gql.go' | xargs --no-run-if-empty rm

update-sample: delete-generated
	go generate sample/resolver.go
	go generate sample/enum/enum.go

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

%_string.go: %.go
	go generate $^

$(BUILD_DIR)/gqlcodegen: $(BUILD_DIR) $(SRC) $(STRING_GO)
	go build -o $@ cmd/gqlcodegen.go
