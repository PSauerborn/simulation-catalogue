locals {
  container_images = [
    "web"
  ]
}

# ECR repositories
resource "aws_ecr_repository" "ecr_repositories" {
  for_each = toset(local.container_images)

  name                 = "${var.app_name}/${each.value}"
  image_tag_mutability = "MUTABLE"
}
