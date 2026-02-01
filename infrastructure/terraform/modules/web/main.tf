resource "helm_release" "web" {
  name      = "${var.app_name}-web"
  chart     = "modules/helm/charts/web-server"
  namespace = var.namespace

  values = [
    # NOTE: static configuration is kept in a separate values file for clarity
    file("modules/helm/values/web.yaml"),
    # NOTE: yamlencode is used here in conjunction with values instead
    # of set to pass environment-specific values. set requires a string value
    # which makes it hard to pass complex structures
    # like maps or lists. Using yamlencode allows us to pass these complex structures
    # directly as YAML, which Helm can then parse correctly.
    yamlencode({
      image = {
        tag = var.image_tags.web

        pullPolicy = "Always"
      }

      imagePullSecrets = [
        {
          name = "aws-ecr-credentials"
        }
      ]

      fullNameOverride = "${var.app_name}-web"

      ingress = {
        enabled = true
        hosts = [
          {
            paths = [
              {
                path     = "/"
                pathType = "Prefix"
              }
            ]
            host = var.web_subdomain
          }
        ]
      }
    })
  ]
}
