variable "postgres_user" {
  type        = string
  description = "Username for the postgres database."
  sensitive   = true
}

variable "postgres_password" {
  type        = string
  description = "Password for the postgres database."
  sensitive   = true
}
