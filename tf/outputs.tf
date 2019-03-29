output "ApiEndpoint" {
  description = "Api Endpoint"
  value       = "http://${aws_lb.load_balancer.dns_name}/v1/"
}

output "DatabaseEndpoint" {
  description = "Database Endpoint"
  value       = "${module.db.this_db_instance_endpoint}"
}

output "MyPublicIP" {
  description = "Public IP with access to Database"
  value       = "${chomp(data.http.mypublicip.body)}"
}

output "CloudWathLogGroup" {
  description = "Cloudwath Log Group to see API Logs"
  value       = "${aws_cloudwatch_log_group.cloudwatch_log_group.name}"
}