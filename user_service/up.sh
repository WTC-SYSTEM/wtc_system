arg="$1"

if [[ $arg == "d" ]]; then
  docker compose down && docker compose up --build --detach
else
  docker compose down && docker compose up --build
fi
