.PHONY: clear-swag swag-init test

clear-swag:
	@echo "clearing docs dirrectory.."
	# windows
	# del /Q docs\* 2>null || exit 0
	# mac / linux
	rm -rf docs/*

swag-init: clear-swag
	@echo "Generate swagger docs.."
	swag init -g main.go

# run all tests
test:
	go test ./tests/... -v