BUILD_DIR:=build
STRINGER_SRC:=$(shell grep -lr '//go:generate stringer' | grep '\.go$$')
STRING_GO:=$(patsubst %.go, %_string.go, $(STRINGER_SRC))
SRC:=$(shell find . -name '*.go')
.PHONY: all clean delete-generated update-example install

all: $(BUILD_DIR)/gqlcodegen

install:
	go install cmd/gqlcodegen/main.go

clean:
	rm $(STRING_GO)
	rm -rf $(BUILD_DIR)

delete-generated:
	find example/enum -mindepth 1 -type d | xargs --no-run-if-empty rm -rf
	find example -name '*_gql.go' | xargs --no-run-if-empty rm

update-example: delete-generated
	go generate example/resolver.go
	go generate example/enum/enum.go

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

%_string.go: %.go
	go generate $^

$(BUILD_DIR)/gqlcodegen: $(BUILD_DIR) $(SRC) $(STRING_GO)
	go build -o $@ cmd/gqlcodegen/main.go
