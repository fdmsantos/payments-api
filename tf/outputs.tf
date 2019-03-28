output "ApiEndpoint" {
  description = "Api Endpoint"
  value       = "http://${aws_lb.load_balancer.dns_name}/v1/"
}

output "DatabaseEndpoint" {
  description = "Database Endpoint"
  value = "${module.db.this_db_instance_endpoint}"
}