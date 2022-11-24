.PHONY: logservers
logservers:
	@echo 'Restarting logserver1...'
	docker kill logserver1
	docker rm logserver1
	docker run -d --name logserver1 chentex/random-logger:latest 1000 4000
	@echo 'Restarting logserver2...'
	docker kill logserver2
	docker rm logserver2
	docker run -d --name logserver2 chentex/random-logger:latest 1000 4000

.PHONY: stop
stop:
	@echo 'Stopping logservers...'
	docker kill logserver1 logserver2

.PHONY: start
start:
	@echo 'Starting logservers...'
	docker start logserver1 logserver2

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/centralog.proto