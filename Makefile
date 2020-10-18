PKG_SRC_PATH := $(GOPATH)/src/github.com/adrg/go-wkhtmltopdf

deploy:
	@rm -rf $(PKG_SRC_PATH)
	@mkdir -p $(PKG_SRC_PATH)
	@cp -r ./ $(PKG_SRC_PATH)
	@rm -rf $(PKG_SRC_PATH)/.git

install: deploy
	@cd "$(PKG_SRC_PATH)" && go install
