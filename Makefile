# include .env
# export

.SILENT:
run:
	go run main.go

build:
	fyne package -os windows --app-build 1 --release