terraform {
  required_providers {
    minecraft = {
      source  = "local/hashicraft/minecraft"
      version = "0.1.0"
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
  schema = "../../../schemas/car.zip"
}