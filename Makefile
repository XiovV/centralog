.PHONY: logservers
logservers:
	@echo 'Restarting logserver1...'
	docker stop logserver1
	docker rm logserver1
	docker run -d --name logserver1 chentex/random-logger:latest 1000 4000
	@echo 'Restarting logserver2...'
	docker stop logserver2
	docker rm logserver2
	docker run -d --name logserver2 chentex/random-logger:latest 1000 4000

.PHONY: stop
stop:
	@echo 'Stopping logservers...'
	docker stop logserver1 logserver1

.PHONY: start
start:
	@echo 'Starting logservers...'
	docker start logserver1 logserver2
