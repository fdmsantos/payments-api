# REST API Test (Form3)

REST API built in Golang with a PostgreSQL database. 

Api Design is [summarised here](design.pdf).

## Deploy

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
