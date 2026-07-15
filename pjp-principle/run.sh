docker compose --file ./docker-compose.yml --env-file development.env up -d
docker compose --file ./docker-compose.yml down
docker rmi scylla-pjp
