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

resource "minecraft_block" "example" {
  x = -1273
  y = 24
  z = 288
  material = "minecraft:stone"
}

resource "minecraft_block" "example2" {
  x = -1273
  y = 25
  z = 288
  material = "minecraft:stone"
}
