terraform {
  required_providers {
    minecraft = {
      source  = "nicholasjackson/mc"
      version = "0.1.3"
    }
  }
}

provider "minecraft" {
  endpoint = "http://workshop.hashicraft.com:9090"
  api_key = "supertopsecret"
}

data "minecraft_block" "existing" {
  x = -1278
  y = 23
  z = 152
}

module "tower" {
  source = "./module"
  height = 10
}

resource "minecraft_schema" "bus" {
  x = -1278
  y = 24
  z = 129
  rotation = 90
  schema = "./car.zip"
}

output "existing" {
  value = data.minecraft_block.existing.id
}
