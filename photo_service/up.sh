arg="$1"
echo "Stopping photo service"
docker compose stop photo_service
echo "Removing photo service"
docker-compose rm -f photo-microservice
if [[ $arg == "d" ]]; then
  docker compose up --build --detach
else
  docker compose up --build
fi
