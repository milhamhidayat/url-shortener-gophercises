# url-shortener-gophercises

Url Shortener Gophercises
# Local Installation

1. Clone this project

2. Install all package

	```sh
	go get ./...
	```
3. Run project

	```sh
	go run main/main.go flag(choose one flag)
	```
	| Flag | Description | Example |
	|--|--| -- |
	| yaml | Yaml file default (urls.yaml) | ```yaml=urls.yaml``` |
	| json | Json file default (urls.json) | ```json=urls.json```
	| db | bolt db option | boltdb

