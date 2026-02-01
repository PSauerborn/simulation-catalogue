terraform {
  backend "s3" {
    key          = "state/app=simulation_catalogue/state.tfstate"
    use_lockfile = true
    bucket       = "tfstate.s31-software.com"

    assume_role = {
      role_arn = "arn:aws:iam::450876623734:role/terraform-ci"
    }
  }
}

provider "aws" {
  region = "us-east-1"

  assume_role {
    role_arn = "arn:aws:iam::450876623734:role/terraform-ci"
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes = {
    config_path = "~/.kube/config"
  }
}
