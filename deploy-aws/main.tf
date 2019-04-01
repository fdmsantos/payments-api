provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "available" {}

data "aws_ecr_repository" "repository" {
  name = "${var.repository}"
}

// Get My Public IP To access to database
data "http" "mypublicip" {
  url = "http://ipv4.icanhazip.com"
}

// Create VPC
// Module Source: https://github.com/terraform-aws-modules/terraform-aws-vpc
module "vpc" {
  source               = "terraform-aws-modules/vpc/aws"
  name                 = "${var.api_name}-${var.env}-vpc"
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

// Create AWS RDS Postgres DB
// Module Source: https://github.com/terraform-aws-modules/terraform-aws-rds
resource "aws_security_group" "db_security_group" {
  name        = "${var.api_name}-${var.env}-db-sg"
  description = "Database Security Group"
  vpc_id      = "${module.vpc.vpc_id}"

  ingress {
    from_port       = "${var.db_port}"
    protocol        = "tcp"
    to_port         = "${var.db_port}"
    security_groups = [ "${aws_security_group.container_security_group.id}" ] // Make database acessible by Containers
  }

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_security_group" "db_security_group_myaccess" {
  name        = "${var.api_name}-${var.env}-db-sg-myaccess"
  description = "Database Security Group to give access to my PC"
  vpc_id      = "${module.vpc.vpc_id}"

  ingress {
    from_port       = "${var.db_port}"
    protocol        = "tcp"
    to_port         = "${var.db_port}"
    cidr_blocks = ["${chomp(data.http.mypublicip.body)}/32"] // Make database acessible by Mr
  }

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

module "db" {
  source                  = "terraform-aws-modules/rds/aws"
  identifier              = "${var.api_name}-${var.env}-db"
  engine                  = "${var.db_engine}"
  engine_version          = "${var.db_engine_version}"
  instance_class          = "${var.db_instance_class}"
  allocated_storage       = "${var.db_storage}"
  storage_encrypted       = false
  name                    = "${var.db_name}"
  storage_type            = "gp2"
  family                  = "${var.db_engine}10"

  username                = "${var.db_username}"
  password                = "${var.db_password}"
  port                    = "${var.db_port}"

  maintenance_window      = "Mon:00:00-Mon:03:00"
  backup_window           = "03:00-06:00"
  backup_retention_period = 0
  skip_final_snapshot     = true
  deletion_protection     = false

  subnet_ids              = ["${module.vpc.public_subnets}"]
  vpc_security_group_ids  = [ "${aws_security_group.db_security_group.id}", "${aws_security_group.db_security_group_myaccess.id}" ]

  multi_az                = false
  publicly_accessible     = true


  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

// Create Aws Iam Role to use by Containers
resource "aws_iam_role" "ecs_task_execution_role" {
  name               = "${var.api_name}-${var.env}-ecs-IamExecutionRole"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "ecs-tasks.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  EOF
  path = "/"
  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_iam_role_policy" "ecs_task_execution_role_policy" {
  name   = "${var.api_name}-${var.env}-ecs-IamExecutionRolePolicy"
  role   = "${aws_iam_role.ecs_task_execution_role.id}"
  policy = <<EOF
{
    "Statement": [
        {
            "Action": [
                "ecr:GetAuthorizationToken",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "*",
            "Effect": "Allow"
        }
    ]
}
EOF
}

// Create Cloudwatch Log Group
resource "aws_cloudwatch_log_group" "cloudwatch_log_group" {
  name              = "${var.api_name}-${var.env}-cloudwatch-log-group"
  retention_in_days = 1

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

// Create Load Balancer
resource "aws_security_group" "load_balancer_security_group" {
  name_prefix = "${var.api_name}-${var.env}-lb-sg"
  description = "Security group for loadbalancer to services on ECS"
  vpc_id = "${module.vpc.vpc_id}"

  ingress {
    from_port = 0
    protocol = "-1"
    to_port = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    protocol = "-1"
    to_port = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_lb" "load_balancer" {
  name               = "${var.api_name}-${var.env}-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = ["${aws_security_group.load_balancer_security_group.id}"]
  subnets            = ["${module.vpc.public_subnets[0]}", "${module.vpc.public_subnets[1]}"]

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_lb_target_group" "load_balancer_default_target_group" {
  name = "${var.api_name}-${var.env}-lb-tg-default"
  vpc_id = "${module.vpc.vpc_id}"
  protocol = "HTTP"
  port = 80
  tags = {
    Terraform  = "true"
    Environment = "${var.env}"
  }
}

resource "aws_lb_target_group" "load_balancer_target_group" {
  name     = "${var.api_name}-${var.env}-lb-tg"
  vpc_id   = "${module.vpc.vpc_id}"
  protocol = "HTTP"
  port     = "80"

  health_check {
    interval          = 10
    path              = "/v1/health"  // Endpoint to check api healthy
    protocol          = "HTTP"
    timeout           = "5"
    healthy_threshold = 10
    matcher           = "200"
  }

  target_type = "ip"

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_lb_listener" "load_balancer_listener" {
  load_balancer_arn = "${aws_lb.load_balancer.arn}"
  port              = 80

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.load_balancer_default_target_group.arn}"
  }

}

resource "aws_lb_listener_rule" "listener_rule" {
  listener_arn = "${aws_lb_listener.load_balancer_listener.arn}"
  priority     = 2

  action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.load_balancer_target_group.arn}"
  }

  condition {
    field  = "path-pattern"
    values = ["/v1/*"]
  }
}

// Create ECS Cluster
// Module Source: https://github.com/terraform-aws-modules/terraform-aws-ecs
module "ecs" {
  source = "terraform-aws-modules/ecs/aws"
  name   = "${var.api_name}-${var.env}-ecs-cluster"

  tags = {
      Terraform   = "true"
      Environment = "${var.env}"
  }
}

resource "aws_security_group" "container_security_group" {
  name_prefix  = "${var.api_name}-${var.env}-container-sg"
  description = "For ecs containers"
  vpc_id      = "${module.vpc.vpc_id}"

  ingress {
    from_port       = 0
    protocol        = "-1"
    to_port         = 0
    security_groups = ["${aws_security_group.load_balancer_security_group.id}"]

  }

  egress {
    from_port = 0
    protocol = "-1"
    to_port = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Terraform  = "true"
    Environment = "${var.env}"
  }
}

resource "template_file" "task_definition" {
  template = "${file("task-definition.json.tmpl")}"
  vars {
    name                 = "${var.api_name}"
    image                = "${data.aws_ecr_repository.repository.repository_url}:${var.image_tag}"
    cpu                  = 256
    memory               = 512
    containerPort        = "${var.container_port}"
    protocol             = "tcp"
    awsRegion            = "${var.aws_region}"
    dbUser               = "${var.db_username}"
    dbPass               = "${var.db_pass}"
    dbName               = "${var.db_name}"
    dbHost               = "${element(split(":", module.db.this_db_instance_endpoint), 0)}"
    dbPort               = "${element(split(":", module.db.this_db_instance_endpoint), 1)}"
  }
}

resource "aws_ecs_task_definition" "ecs_task" {
  family                   = "${var.api_name}"
  container_definitions    = "${template_file.task_definition.rendered}"
  cpu                      = "256"
  memory                   = "512"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = "${aws_iam_role.ecs_task_execution_role.arn}"

  tags = {
    Terraform   = "true"
    Environment = "${var.env}"
  }
}

resource "aws_ecs_service" "ecs_service" {
  name                               = "${var.api_name}-${var.env}-ecs-service"
  depends_on                         = ["aws_lb_listener_rule.listener_rule"]
  cluster                            = "${module.ecs.this_ecs_cluster_id}"
  task_definition                    = "${aws_ecs_task_definition.ecs_task.arn}"
  launch_type                        = "FARGATE"
  desired_count                      = 2
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 70
  network_configuration {
    subnets          = ["${module.vpc.public_subnets[0]}", "${module.vpc.public_subnets[1]}"]
    security_groups  = ["${aws_security_group.container_security_group.id}"]
    assign_public_ip = true
  }

  load_balancer {
    container_name   = "${var.api_name}"
    container_port   = "${var.container_port}"
    target_group_arn = "${aws_lb_target_group.load_balancer_target_group.arn}"
  }

}




