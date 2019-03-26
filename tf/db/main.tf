//resource "aws_db_instance" "database" {
//  allocated_storage    = 20
//  engine               = "postgres"
//  engine_version       = "10.6"
//  instance_class       = "db.t2.micro"
//  identifier           = "payments-api-database"
//  name                 = "payments"
//  username             = "api"
//  password             = "test123456"
//  multi_az             = false
//  publicly_accessible  = true
//}
//
//


provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "available" {}

module "vpc" {
  source               = "terraform-aws-modules/vpc/aws"
  name                 = "${var.name}-${var.env}-vpc"
  cidr                 = "${var.vpc_cidr}"
  public_subnets       = "${var.public_subnets}"
  azs                  = ["${data.aws_availability_zones.available.names[0]}", "${data.aws_availability_zones.available.names[1]}"]
  enable_dns_hostnames = true
  enable_dns_support   = true
  instance_tenancy     = "default"
  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}