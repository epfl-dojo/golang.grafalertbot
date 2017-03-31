build: get_dep
	go build -o grafalert
get_dep:
	go get "github.com/deckarep/golang-set"
	go get "gopkg.in/telegram-bot-api.v4"
