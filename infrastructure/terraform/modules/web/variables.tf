variable "app_name" {
  type        = string
  description = "Name of the application"
}

variable "web_subdomain" {
  type        = string
  description = "Subdomain for the web application"
}

variable "namespace" {
  type        = string
  description = "Namespace for the web application"
}

variable "image_tags" {
  type = object({
    web = string
  })
}
