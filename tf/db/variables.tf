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

variable "db_identifier" {
  default = "payments-api-database"
}

variable "db_name" {
}

variable "db_username" {

}

variable "db_password" {

}
