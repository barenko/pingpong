build:
	go build

dev:
	go run ./pingpong.go

dev-ping:
	PORT=3000 PONG=http://0.0.0.0:3001 go run ./pingpong.go

dev-pong:
	PORT=3001 PING=http://0.0.0.0:3000 go run ./pingpong.go

ping:
	LOG=false PORT=3000 PONG=http://0.0.0.0:3001 ./pingpong

pong:
	LOG=false PORT=3001 PING=http://0.0.0.0:3000 ./pingpong
