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
	git push origin main

setup-git:
	git config --global user.name "github-actions[bot]"
	git config --global user.email "github-actions[bot]@users.noreply.github.com"
	git config --global push.default simple
	git config --global pull.rebase false
	git config --global pull.ff only
