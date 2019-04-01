variable "aws_region" {
  default = "eu-west-1"
}

variable "vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "public_subnets" {
  type    = "list"
  default = [ "10.0.0.0/20", "10.0.32.0/20" ]
}

variable "env" {
  default = "dev"
}

variable "api_name" {
  default = "payment-api"
}

variable "container_port" {
  default = 8000
}

variable "db_pass" {
}

variable "repository" {
}

variable "image_tag" {
  default = "v1"
}


variable "db_storage" {
  default = 20
}

variable "db_engine" {
  default = "postgres"
}

variable "db_engine_version" {
  default = "10.6"
}

variable "db_instance_class" {
  default = "db.t2.micro"
}

variable "db_name" {
}

variable "db_username" {

}

variable "db_password" {

}

variable "db_port" {
  default = "5432"
}
