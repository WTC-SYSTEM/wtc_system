arg="$1"
echo "Stopping gateway service"
docker compose stop api_gateway
echo "Removing gateway service"
docker-compose rm -f api-gateway-microservice
if [[ $arg == "d" ]]; then
  docker compose up --build --detach
else
  docker compose up --build
fi
