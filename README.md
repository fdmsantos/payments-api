# REST API Payments

REST API built in Golang with a PostgreSQL database. 

Api Design is [summarised here](design.pdf).

## Running the tests

```sh
git clone https://github.com/fdmsantos/payments-api
cd payments-api
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword postgres
export DB_HOST=localhost
export DB_NAME=paymentsapi
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
git clone https://github.com/fdmsantos/payments-api
cd payments-api
docker-compose build
docker-compose up -d 
```

Connect to localhost:8000/v1/user To create new User.

### Aws with terraform (Work in Progress)

**Requirements**

- Docker
- Terraform
- Aws Cli with login configured

**Steps**

1. Containerize Payment REST API application
2. Use AWS CLI to create Amazon ECR repository
3. Build docker image and push to ECR
4. Create VPC, Subnets, InternetGateway, Route tables
5. Create IAM role
6. Create ECS Cluster, Loadbalancer & Listener, Security groups etc
7. Deploy docker container

```sh
git clone https://github.com/fdmsantos/payments-api
cd payments-api/tf
docker build -t payments-api .
aws ecr create-repository --repository-name payments-api
aws ecr get-login --no-include-email | sh
IMAGE_REPO=$(aws ecr describe-repositories --repository-names payments-api --query 'repositories[0].repositoryUri' --output text)
docker tag payments-api:latest $IMAGE_REPO:v1
docker push $IMAGE_REPO:v1
terraform init
# To get Image ID
echo $IMAGE_REPO
terraform apply
```

## Usage

### New User

```sh
curl --request POST \
  --url http://localhost:8000/v1/user \
  --header 'content-type: application/json' \
  --data '{
	"email": "fabiosantos@gmail.com",
	"password": "secretpassword"

}'
```

### Login

```sh
curl --request POST \
  --url http://localhost:8000/v1/user/login \
  --header 'content-type: application/json' \
  --data '{
	"email": "fabiosantos@gmail.com",
	"password": "secretpassword"

}'
```

### Login

```sh
curl --request POST \
  --url http://localhost:8000/v1/user/login \
  --header 'content-type: application/json' \
  --data '{
	"email": "fabiosantos@gmail.com",
	"password": "secretpassword"

}'
```

### New Payment

```sh
curl --request POST \
  --url http://localhost:8000/v1/payments \
  --header 'authorization: Bearer $token' \
  --data '{
      "type": "Payment",
      "id": "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
      "version": 0,
      "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
      "attributes": {
        "amount": "100.21",
        "beneficiary_party": {
          "account_name": "W Owens",
          "account_number": "31926819",
          "account_number_code": "BBAN",
          "account_type": 0,
          "address": "1 The Beneficiary Localtown SE2",
          "bank_id": "403000",
          "bank_id_code": "GBDSC",
          "name": "Wilfred Jeremiah Owens"
        },
        "charges_information": {
          "bearer_code": "SHAR",
          "sender_charges": [
            {
              "amount": "5.00",
              "currency": "GBP"
            },
            {
              "amount": "10.00",
              "currency": "USD"
            }
          ],
          "receiver_charges_amount": "1.00",
          "receiver_charges_currency": "USD"
        },
        "currency": "GBP",
        "debtor_party": {
          "account_name": "EJ Brown Black",
          "account_number": "GB29XABC10161234567801",
          "account_number_code": "IBAN",
          "address": "10 Debtor Crescent Sourcetown NE1",
          "bank_id": "203301",
          "bank_id_code": "GBDSC",
          "name": "Emelia Jane Brown"
        },
        "end_to_end_reference": "Wil piano Jan",
        "fx": {
          "contract_reference": "FX123",
          "exchange_rate": "2.00000",
          "original_amount": "200.42",
          "original_currency": "USD"
        },
        "numeric_reference": "1002001",
        "payment_id": "123456789012345678",
        "payment_purpose": "Paying for goods/services",
        "payment_scheme": "FPS",
        "payment_type": "Credit",
        "processing_date": "2017-01-18",
        "reference": "Payment for Em\u0027s piano lessons",
        "scheme_payment_sub_type": "InternetBanking",
        "scheme_payment_type": "ImmediatePayment",
        "sponsor_party": {
          "account_number": "56781234",
          "bank_id": "123123",
          "bank_id_code": "GBDSC"
        }
      }
    }'
```

### Update Payment

```sh
curl --request PUT \
  --url http://localhost:8000/v1/payments/216d4da9-e59a-4cc6-8df3-3da6e7580b77 \
  --header 'authorization: Bearer $token' \
  --data '
{
	"type": "Payment Update",
	"id": "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
	"version": 0,
	"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"attributes": {
		"amount": "100.21",
		"beneficiary_party": {
			"account_name": "W Oddasdwens",
			"account_number": "31926819",
			"account_number_code": "BBAN",
			"account_type": 0,
			"address": "1 The Beneficiary Localtown SE2",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Wilfred Jeremiah Owens"
		},
		"charges_information": {
			"bearer_code": "SHAR",
			"sender_charges": [
				{
					"amount": "5.00",
					"currency": "GBP"
				},
				{
					"amount": "10.00",
					"currency": "USD"
				}
			],
			"receiver_charges_amount": "1.00",
			"receiver_charges_currency": "USD"
		},
		"currency": "GBP",
		"debtor_party": {
			"account_name": "EJ Brown Black",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "10 Debtor Crescent Sourcetown NE1",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "Emelia Jane Brown"
		},
		"end_to_end_reference": "Wil piano Jan",
		"fx": {
			"contract_reference": "FX123",
			"exchange_rate": "2.00000",
			"original_amount": "200.42",
			"original_currency": "USD"
		},
		"numeric_reference": "1002001",
		"payment_id": "123456789012345678",
		"payment_purpose": "Paying for goods/services",
		"payment_scheme": "FPS",
		"payment_type": "Credit",
		"processing_date": "2017-01-18",
		"reference": "Payment for Em\u0027s piano lessons",
		"scheme_payment_sub_type": "InternetBanking",
		"scheme_payment_type": "ImmediatePayment",
		"sponsor_party": {
			"account_number": "56781234",
			"bank_id": "123123",
			"bank_id_code": "GBDSC"
		}
	}
}'
```


### Get Payment

```sh
curl --request POST \
  --url http://localhost:8000/v1/user/login \
  --header 'content-type: application/json' \
  --data '{
	"email": "fabiosantos@gmail.com",
	"password": "secretpassword"

}'
```

### Get All Payments

```sh
curl --request GET \
  --url http://localhost:8000/v1/payments/216d4da9-e59a-4cc6-8df3-3da6e7580b77 \
  --header 'authorization: Bearer $token'
```

### Delete Payment

```sh
curl --request DELETE \
  --url http://localhost:8000/v1/payments/216d4da9-e59a-4cc6-8df3-3da6e7580b77 \
  --header 'authorization: Bearer $token'
```

## Future Improvements

* Add Log  [logrus](https://github.com/sirupsen/logrus)
* Add Documentation [godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)
* Create Api Documentation with [swagger](https://github.com/go-swagger/go-swagger)
* Refactoring The code. (Move Validations from Controllers to Models)
* Deploy in AWS ECS with Terraform