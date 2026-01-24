# Database URLs
DB_MAIN := postgres://postgres:test@localhost:5432/messenger
DB_TEST := postgres://postgres:test@localhost:5433/messenger
MIG_DIR := migrations
GOOSE := goose

.PHONY: mig mig-test mig-up mig-down mig-reset mig-test-up mig-test-down mig-test-reset

# top-level convenience targets: make mig <cmd>
mig: mig-$(cmd)
	@:

mig-test: mig-test-$(cmd)
	@:

# default cmd if none provided
cmd ?= up

# main DB handlers
mig-up:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_MAIN)" up

mig-down:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_MAIN)" down

mig-reset:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_MAIN)" reset
# test DB handlers
mig-test-up:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_TEST)" up

mig-test-down:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_TEST)" down

mig-test-reset:
	$(GOOSE) -dir $(MIG_DIR) postgres "$(DB_TEST)" reset