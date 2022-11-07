container "vscode" {
  network {
      name = "network.local"
  }

  image {
    name = "nicholasjackson/codeserver:tfd"
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
    name = "shipyardrun/docs:v0.6.0"
  }

  index_title = "Docs"
}

network "local" {
  subnet = "10.0.0.0/16"
}
