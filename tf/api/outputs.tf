output "VPC" {
  description = "The ID of the VPC"
  value       = "${module.vpc.vpc_id}"
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = ["${module.vpc.public_subnets}"]
}

output "iam_role_arn" {
  description = "Iam Role Arn"
  value       = "${aws_iam_role.ecs_task_execution_role.arn}"
}


output "cluster" {
  description = "Cluster"
  value       = "${module.ecs.this_ecs_cluster_id}"
}

output "Listener" {
  description = "Listener"
  value       = "${aws_lb_listener_rule.listener_rule.arn}"
}

output "ContainerSecurityGroup" {
  description = "ContainerSecurityGroup"
  value       = "${aws_security_group.container_security_group.arn}"
}

output "LoadBalancerDNS" {
  description = "LoadBalancerDNS"
  value       = "${aws_lb.load_balancer.dns_name}"
}


output "ApiEndpoint" {
  description = "ApiEndpoint"
  value       = "http://${aws_lb.load_balancer.dns_name}/v1/"
}
