PKG_SRC_PATH := $(PKG_SRC_PATH)

all:
ifeq ($(PKG_SRC_PATH),)
	$(error environment variable PKG_SRC_PATH is not set)
endif
	@echo "PKG_SRC_PATH=$(PKG_SRC_PATH)"

deploy: all
	@rm -rf $(PKG_SRC_PATH)
	@mkdir -p $(PKG_SRC_PATH)
	@cp -r ./ $(PKG_SRC_PATH)
	@rm -rf $(PKG_SRC_PATH)/.git

install: deploy
	@cd "$(PKG_SRC_PATH)" && go install
