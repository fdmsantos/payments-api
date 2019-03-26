
provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "available" {}

module "vpc" {
  source                                 = "terraform-aws-modules/vpc/aws"
  name                                   = "${var.name}-${var.env}-vpc"
  cidr                                   = "${var.vpc_cidr}"
  public_subnets                         = "${var.public_subnets}"
  azs                                    = ["${data.aws_availability_zones.available.names[0]}", "${data.aws_availability_zones.available.names[1]}"]
  enable_dns_hostnames                   = true
  enable_dns_support                     = true
  instance_tenancy                       = "default"
  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_db_subnet_group" "default" {
  name       = "${var.name}-${var.env}-db-subnet-group"
  subnet_ids = ["${module.vpc.public_subnets[0]}", "${module.vpc.public_subnets[1]}"]

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_security_group" "db_security_group" {
  name        = "${var.name}-${var.env}-db-security-group"
  description = "Database Security Group"
  vpc_id      = "${module.vpc.vpc_id}"

  ingress {
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_db_instance" "database" {
  allocated_storage       = "${var.db_storage}"
  storage_type            = "gp2"
  engine                  = "${var.db_engine}"
  engine_version          = "${var.db_engine_version}"
  instance_class          = "${var.db_instance_class}"
  identifier              = "${var.db_identifier}"
  name                    = "${var.db_name}"
  username                = "${var.db_username}"
  password                = "${var.db_password}"
  multi_az                = false
  publicly_accessible     = true
  db_subnet_group_name    = "${aws_db_subnet_group.default.name}"
  backup_retention_period = 0
  skip_final_snapshot     = true
  vpc_security_group_ids  = ["${aws_security_group.db_security_group.id}"]
}
