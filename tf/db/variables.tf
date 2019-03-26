variable "aws_region" {
  default = "eu-west-1"
}

variable "vpc_cidr" {
  default = "11.0.0.0/16"
}

variable "public_subnets" {
  type    = "list"
  default = [ "11.0.0.0/20", "11.0.32.0/20" ]
}

variable "env" {
  default = "dev"
}

variable "name" {
  default = "payments-database"
}
