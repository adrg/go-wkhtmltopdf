PKG_SRC_PATH := $(GOPATH)/src/github.com/leandrosilva/go-wkhtmltopdf

copy:
	@rm -rf $(PKG_SRC_PATH)
	@mkdir -p $(PKG_SRC_PATH)
	@cp -r ./ $(PKG_SRC_PATH)
	@rm -rf $(PKG_SRC_PATH)/.git

install: copy
	@ls $(PKG_SRC_PATH)
	@cd "$(PKG_SRC_PATH)" && go install

build-example:
	@cd ./example && go build -o run-example.exe main.go

run-example:
	@cd ./example && ./run-example.exe

clean-example:
	@rm ./example/run-example.exe
	@rm ./example/example.pdf
