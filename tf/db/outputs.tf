output "DB_ENDPOINT" {
  description = "Database Endpoint"
  value       = "${aws_db_instance.database.endpoint}"
}
