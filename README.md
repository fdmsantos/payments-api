# REST API Payments

REST API built in Golang with a PostgreSQL database. 

Api Design: [click me](Form3-Payments-api.pdf)

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
See the [examples](#Examples)

### Aws with terraform

**Requirements**

- Docker
- Terraform
- Aws Cli with login configured

**Steps**

- Clone Repository

```sh
git clone https://github.com/fdmsantos/payments-api
cd payments-api
```

- Use AWS CLI to create AWS ECR repository

```sh
aws ecr create-repository --repository-name payments-api
aws ecr get-login --no-include-email | sh
IMAGE_REPO=$(aws ecr describe-repositories --repository-names payments-api --query 'repositories[0].repositoryUri' --output text)
```

-  Build docker image and push to AWS ECR

```sh
docker build -t payments-api .
docker tag payments-api:latest $IMAGE_REPO:v1
docker push $IMAGE_REPO:v1
```

- Deploy to AWS with terraform [Aws Infrastructure Diagram](tf/AWSDiagram.pdf)

```sh
cd tf
terraform init
terraform apply
```

- Use ApiEndpoint Terraform Output to access to api. Should take a few seconds until the service is up.

See the [examples](#Examples). Replace http://localhost:8000/v1/ by ApiEndpoint


```sh
curl --request POST \
  --url http://payment-api-dev-lb-634579817.eu-west-1.elb.amazonaws.com/v1/user \
  --header 'content-type: application/json' \
  --data '{
	"email": "fabiosantos@gmail.com",
	"password": "secretpassword"

}'
```

- Use DatabaseEndpoint Terraform Output to connect to database. It was created a security group to give access from your PC to Database (Using your Public IP)

```sh
# Necessary have postgres client installed
# Check terraform.tfvars file to get database credentials
# psql -h payment-api-dev-db.culfnfuxbney.eu-west-1.rds.amazonaws.com -U api -d payments
Password for user api: 
psql (11.2, server 10.6)
SSL connection (protocol: TLSv1.2, cipher: ECDHE-RSA-AES256-GCM-SHA384, bits: 256, compression: off)
Type "help" for help.

payments=> \dt
               List of relations
 Schema |         Name         | Type  | Owner 
--------+----------------------+-------+-------
 public | accounts             | table | api
 public | attributes           | table | api
 public | beneficiary_parties  | table | api
 public | charges              | table | api
 public | charges_informations | table | api
 public | debtor_parties       | table | api
 public | fxes                 | table | api
 public | payments             | table | api
 public | sponsor_parties      | table | api
(9 rows)

payments=> 

```

- Navigate to AWS Cloudwatch and Use CloudWathLogGroup Terraform Output to discover what is the api cloud watch group and check the logs

- Destroy

```sh
terraform destroy
```

## Examples

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

## Running the tests

```sh
git clone https://github.com/fdmsantos/payments-api
cd payments-api
docker run -it -p 5432:5432 -e POSTGRES_PASSWORD="api" -e POSTGRES_DB="api" -e POSTGRES_USER="api" --name postgresDB -d  postgres
export DB_HOST=localhost
export DB_NAME=api
export DB_USER=api
export DB_PASS=api
export DB_PORT=5432
go test ./...
```