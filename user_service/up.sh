arg="$1"
echo "Stopping user service"
docker compose stop user_service
echo "Removing user service"
docker-compose rm -f user-microservice
if [[ $arg == "d" ]]; then
  docker compose up --build --detach
else
  docker compose up --build
fi
