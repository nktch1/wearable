up:
	source build/.env && \
		docker compose -f build/docker-compose.yaml up --force-recreate