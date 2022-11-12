terraform {
  required_providers {
    minecraft = {
      source  = "local/nicholasjackson/workshop"
      version = "0.1.1"
    }
  }
}

provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
