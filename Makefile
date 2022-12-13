fetch:
	go run main.go people:fetch

fuse:
	go run main.go people:fuse

build:
	make fetch
	make fuse

update:
	make build
	git add data/*.json
	git add api/*.json
	git commit -m "Update data"
	git push