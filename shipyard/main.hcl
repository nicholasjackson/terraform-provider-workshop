variable "disable_vscode" {
  default = false
}

container "vscode" {
  disabled = var.disable_vscode

  network {
      name = "network.local"
  }

  image {
    name = "nicholasjackson/vscodeserver:tfw"
  }

  port {
    local  = 8000
    host   = 8000
    remote = 8000
  }
  
  volume {
    source = "../"
    destination = "/home/src"
  }
}

docs "docs" {
  port = 3000
  
  network {
      name = "network.local"
  }

  path = "./docs"
  
  image {
    name = "shipyardrun/docs:v0.6.1"
  }

  index_title = "Docs"
}

network "local" {
  subnet = "10.0.0.0/16"
}
