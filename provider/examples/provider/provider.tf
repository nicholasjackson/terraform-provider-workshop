terraform {
  required_providers {
    minecraft = {
      source  = "local/nicholasjackson/minecraft"
      version = "0.1.1"
    }
  }
}

provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
