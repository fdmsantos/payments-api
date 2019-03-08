# REST API Test (Form3)

REST API built in Golang with a PostgreSQL database. 

Api Design is [summarised here](design.pdf).

## Running the tests

```sh
git clone https://github.com/fdmsantos/test
cd test
docker run -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword postgres -d
export DB_HOST=localhost
export DB_NAME=payments-api
export DB_USER=postgres
export DB_PASS=mysecretpassword
export DB_PORT=5432
go test ./...
```

## Deployment

### Docker-Compose

**Requirements**

- Docker and Docker-compose
- Nothing running on port 8000 and 5432

```sh
git clone https://github.com/fdmsantos/test
cd test
docker-compose build
docker-compose up -d 
```

Connect to localhost:8000/v1/user To create new User.

## Future Improvements

* Add Log  [logrus](https://github.com/sirupsen/logrus)
* Add Documentation [godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)
* Create Api Documentation with [swagger](https://github.com/go-swagger/go-swagger)
* Deploy in AWS ECS with Terraform