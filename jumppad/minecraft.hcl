resource "copy" "server_files" {
  source      = "./mc_server"
  destination = data("mc_server")
}

variable "server_source" {
  //default = "./mc_server"
  default = data("mc_server")
}

resource "container" "minecraft" {

  network {
    id = resource.network.main.meta.id
  }

  image {
    name = "hashicraft/minecraft:v1.18.2-fabric"
  }

  volume {
    source      = "${variable.server_source}/mods"
    destination = "/minecraft/mods"
  }

  volume {
    source      = "${variable.server_source}/plugins"
    destination = "/minecraft/plugins"
  }

  volume {
    source      = "${variable.server_source}/world"
    destination = "/minecraft/world"
  }

  volume {
    source      = "${variable.server_source}/worlds"
    destination = "/minecraft/worlds"
  }

  volume {
    source      = "${variable.server_source}/config"
    destination = "/minecraft/config"
  }

  port {
    local  = 25565
    remote = 25565
    host   = 25565
  }

  port {
    local  = 27015
    remote = 27015
    host   = 27015
  }

  port {
    local  = 9090
    remote = 9090
    host   = 9090
  }

  environment = {
    JAVA_MEMORY       = "4G"
    MINECRAFT_MOTD    = "HashiCraft"
    WHITELIST_ENABLED = "false"
    RCON_PASSWORD     = "password"
    RCON_ENABLED      = "true"
    ONLINE_MODE       = "false"
  }
}

resource "container" "minecraft_web" {
  network {
    id = resource.network.main.meta.id
  }

  image {
    name = "hashicraft/minecraft-web:0.4.0"
  }

  port {
    local  = 8080
    remote = 8080
    host   = 8080
  }

  volume {
    source      = "./mc_client/config.json"
    destination = "/app/config.json"
  }
}
