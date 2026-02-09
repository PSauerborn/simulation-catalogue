locals {
  app_name = "simulation-catalogue"

  api_subdomain = "simulation-catalogue.s31-software.com"
  web_subdomain = "simulation-catalogue.s31-software.com"
}

# namespace to hold the k8s resources
resource "kubernetes_namespace_v1" "namespace" {
  metadata {
    name = local.app_name
  }
}

# ecr credentials expire in 6 hours. deploy initial secret
# and deploy rotator cron job to periodically refresh the secret
module "ecr_creds" {
  depends_on = [kubernetes_namespace_v1.namespace]
  source     = "git::ssh://git@github.com/s31-software-co/terraform-modules.git//k8s-ecr-credentials?ref=1.0.0"

  namespace      = kubernetes_namespace_v1.namespace.metadata[0].name
  deploy_rotator = true
  base_name      = local.app_name

  cron_schedule = "0 * * * *"
}

# postgres credentials for the api
resource "kubernetes_secret_v1" "postgres_creds" {
  metadata {
    name      = "postgres-creds"
    namespace = kubernetes_namespace_v1.namespace.metadata[0].name
  }

  data = {
    username = var.postgres_user
    password = var.postgres_password
  }
}

# web application
module "web" {
  depends_on = [module.ecr_creds, kubernetes_namespace_v1.namespace]
  source     = "./modules/web"

  app_name  = local.app_name
  namespace = kubernetes_namespace_v1.namespace.metadata[0].name

  web_subdomain = local.web_subdomain

  image_tags = {
    web = "0.2.0"
  }
}

# api application
module "api" {
  depends_on = [module.ecr_creds, kubernetes_namespace_v1.namespace]
  source     = "./modules/api"

  app_name  = local.app_name
  namespace = kubernetes_namespace_v1.namespace.metadata[0].name

  api_subdomain = local.api_subdomain

  postgres_config = {
    host                    = "postgres-prod-rw.shared-infra"
    port                    = 5432
    database                = "simulation_catalogue"
    credentials_secret_name = kubernetes_secret_v1.postgres_creds.metadata[0].name
  }

  image_tags = {
    api = "0.2.0"
  }
}
