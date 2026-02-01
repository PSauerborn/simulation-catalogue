resource "helm_release" "api" {
  name      = "${var.app_name}-api"
  chart     = "modules/helm/charts/web-server"
  namespace = var.namespace

  values = [
    # NOTE: static configuration is kept in a separate values file for clarity
    file("modules/helm/values/api.yaml"),
    # NOTE: yamlencode is used here in conjunction with values instead
    # of set to pass environment-specific values. set requires a string value
    # which makes it hard to pass complex structures
    # like maps or lists. Using yamlencode allows us to pass these complex structures
    # directly as YAML, which Helm can then parse correctly.
    yamlencode({
      image = {
        tag = var.image_tags.api

        pullPolicy = "Always"
      }

      imagePullSecrets = [
        {
          name = "aws-ecr-credentials"
        }
      ]

      fullNameOverride = "${var.app_name}-api"

      env = {
        API_VERSION       = "v1"
        API_PORT          = 8080
        POSTGRES_HOST     = var.postgres_config.host
        POSTGRES_PORT     = var.postgres_config.port
        POSTGRES_DATABASE = var.postgres_config.database
      }

      secrets = {
        POSTGRES_USER = {
          secretName = var.postgres_config.credentials_secret_name
          key        = "username"
        }

        POSTGRES_PASSWORD = {
          secretName = var.postgres_config.credentials_secret_name
          key        = "password"
        }
      }

      ingress = {
        enabled = true
        hosts = [
          {
            paths = [
              {
                path     = "/api"
                pathType = "Prefix"
              }
            ]
            host = var.api_subdomain
          }
        ]
      }
    })
  ]
}
