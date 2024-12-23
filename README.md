# Gateway Service


### mod tidy
```bash
export GOPROXY=direct && export GOSUMDB=off  && go mod tidy
```

### run
```bash
export APP_ENV=development && go run ./main.go
```

docker-compose -p gateway-service up --build -d

>>> note
if gateway container was rebuilt make sure to add it to the beauticket_backend network ,
preferably using portainer's dashboard