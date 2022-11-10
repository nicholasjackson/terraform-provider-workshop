terraform {
  required_providers {
    hashicraft = {
      source  = "nicholasjackson/workshop"
      version = "0.1.0"
    }
  }
}

provider "hashicraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}

resource "hashicraft_block" "example" {
  x = 12
  y = 23
  z = 43
  material = "minecraft:stone"
}
