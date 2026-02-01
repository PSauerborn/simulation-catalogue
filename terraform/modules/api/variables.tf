variable "app_name" {
  type        = string
  description = "Name of the application"
}

variable "api_subdomain" {
  type        = string
  description = "Subdomain for the api application"
}

variable "namespace" {
  type        = string
  description = "Namespace for the api application"
}

variable "image_tags" {
  type = object({
    api = string
  })
}

variable "postgres_config" {
  type = object({
    host                    = string
    port                    = number
    credentials_secret_name = string
    database                = string
  })
}
