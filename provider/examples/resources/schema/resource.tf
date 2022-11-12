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

resource "minecraft_schema" "bus" {
  x = -1278
  y = 24
  z = 288
  rotation = 270
  schema = "../../../example_schemas/car.zip"
}
