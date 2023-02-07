SHELL=/bin/bash -euo pipefail

MOCKS_PATH="$(shell pwd)/test/wiremock"
CERTS_PATH="$(shell pwd)/certs"
CONTRACTS_PATH="$(shell pwd)/deps"

compose-run:
	source build/.env && \
		docker compose -f build/docker-compose.yaml up --force-recreate --build

docker-run:
	docker run \
	  -p 9000:9000 \
	  -p 8000:8000 \
	  -v ${MOCKS_PATH}:/home/mock \
	  -v ${CERTS_PATH}:/etc/ssl/mock/share \
	  -v ${CONTRACTS_PATH}:/proto \
	  ${GRPC_WIREMOCK_IMAGE}