generate-mocks:
	minimock -i ./providers.Sender -o ./providers -s _mock.go -g

test:
	go test ./... -race -failfast

.PHONY: deploy
deploy:
	cd deploy && docker-compose up


docker:
	docker build -t email_sender . && docker run -p 5678:5678 email_sender

go-run:
	go run ./cmd/ -config "./configuration"
