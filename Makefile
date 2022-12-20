all: docker-clean docker-build docker-create

docker-build:
	docker build --tag restful-html-games .

docker-create:
	docker create \
		-p 8082:8082 \
		--name restful-html-games \
		restful-html-games

docker-start:
	docker start restful-html-games

docker-stop:
	docker stop restful-html-games

docker-clean:
	-docker rm restful-html-games
	-docker rmi restful-html-games
