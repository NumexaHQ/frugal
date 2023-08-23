IMAGE_TAG=dev
IMAGE_NAME=numexa-auth:$(IMAGE_TAG)
auth:
	docker build -f auth/Dockerfile -t numexa-auth:dev .

monger: 
	docker build -f monger/Dockerfile -t numexa-monger:dev .

vibe: 
	docker build -f vibe/Dockerfile -t numexa-vibe:dev .

up:
	docker-compose up -d

all: auth monger vibe
.PHONY: auth monger vibe up

