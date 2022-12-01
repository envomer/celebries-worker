fetch:
	go run main.go people:fetch

fuse:
	go run main.go people:fuse

build:
	make fetch
	make fuse
