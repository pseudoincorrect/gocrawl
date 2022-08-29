# This Makefile has been replaced by build.sh

MAKEFLAGS += --no-print-directory

SERVICES = worker server

all:
	@ echo "all"

docker-deploy:
	@ cd docker ; \
	docker-compose up -d && docker-compose logs -f

docker-down:
	@ cd docker ; \
	docker-compose down

docker-remove:
	@ docker-compose rm $(SERVICE)

docker-build-all:
	@ for i in $(SERVICES); do \
		$(MAKE) docker-build-one SERVICE="$$i" ; \
	done

docker-build-one:
	@ echo "building service $(SERVICE)"
	@ docker build . -t $(SERVICE) -f ./docker/dockerfile --build-arg SVC=$(SERVICE)

.PHONY: all docker-deploy docker-remove docker-build-all docker-build-one