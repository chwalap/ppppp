# ppppp
Projekt programu pobierającego i prezentującego pogodę

Build
docker-compose build --parallel --progress tty

Deploy
docker stack deploy --compose-file deploy/stack.yml ppppp
