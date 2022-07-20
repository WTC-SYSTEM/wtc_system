arg="$1"
echo "Stopping recipe service"
docker compose stop recipe_service
echo "Removing recipe service"
docker-compose rm -f recipe-microservice
if [[ $arg == "d" ]]; then
  docker compose up --build --detach
else
  docker compose up --build
fi
