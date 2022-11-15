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

data "minecraft_block" "example" {
  x = -1273
  y = 23
  z = 288
}

output "example" {
  value = data.minecraft_block.example.material
}