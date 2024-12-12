# Gateway Service


### mod tidy
```bash
export GOPROXY=direct && export GOSUMDB=off  && go mod tidy
```

### run
```bash
export APP_ENV=development && go run ./main.go
```

docker-compose -p gateway-service up --build
