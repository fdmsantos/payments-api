aws_region = "eu-west-1"
vpc_cidr = "10.0.0.0/16"
public_subnets = [ "10.0.0.0/20", "10.0.32.0/20" ]
env = "dev"
api_name = "payment-api"
container_name  = "payment-api"
container_port = 8000