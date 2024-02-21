resource "network" "main" {
  subnet = "10.0.0.0/16"
}

variable "docs_url" {
  description = "The URL for the documentation site"
  default     = "http://localhost"
}

variable "prismarine_url" {
  description = "The URL for prismarine"
  default     = "http://localhost:8080"
}

variable "minecraft_url" {
  description = "The URL for the Minecraft server"
  default     = "minecraft.container.jumppad.dev"
}

variable "api_url" {
  description = "The URL for the Minecraft API"
  default     = "http://localhost:9090"
}

variable "vscode_token" {
  default = "token"
}

resource "copy" "source_files" {
  source      = "../provider"
  destination = data("provider_source")
}

resource "template" "vscode_jumppad" {
  source = <<-EOF
  {
  "tabs": [
    {
      "name": "Docs",
      "uri": "${variable.docs_url}",
      "type": "browser",
      "active": true
    },
    {
      "name": "Terminal",
      "location": "editor",
      "type": "terminal"
    }
  ]
  }
  EOF

  destination = "${data("vscode")}/workspace.json"
}

resource "template" "vscode_settings" {
  source = <<-EOF
  {
      "workbench.colorTheme": "Palenight Theme",
      "editor.fontSize": 16,
      "workbench.iconTheme": "material-icon-theme",
      "terminal.integrated.fontSize": 16
  }
  EOF

  destination = "${data("vscode")}/settings.json"
}

resource "container" "vscode" {
  network {
    id = resource.network.main.meta.id
  }

  image {
    name = "nicholasjackson/terraform-provider-workshop:v0.2.0"
  }

  volume {
    source      = resource.copy.source_files.destination
    destination = "/provider"
  }

  volume {
    source      = resource.template.vscode_jumppad.destination
    destination = "/provider/.vscode/workspace.json"
  }

  volume {
    source      = resource.template.vscode_settings.destination
    destination = "/provider/.vscode/settings.json"
  }

  environment = {
    CONNECTION_TOKEN = variable.vscode_token
    DEFAULT_FOLDER   = "/provider"
  }

  port {
    local  = 8000
    remote = 8000
    host   = 8000
  }

  health_check {
    timeout = "100s"

    //http {
    //  address       = "http://${resource.docs.docs.fqdn}/docs/provider/introduction/what_is_terraform"
    //  success_codes = [200]
    //}

    http {
      address       = "http://localhost:8000/"
      success_codes = [200, 302, 403]
    }
  }
}

module "workshop" {
  source = "./workshop"

  variables = {
    working_directory = "/provider"
    docs_url          = variable.docs_url
    minecraft_url     = variable.minecraft_url
    prismarine_url    = variable.prismarine_url
    redoc_url         = variable.redoc_url
  }
}

resource "docs" "docs" {
  network {
    id = resource.network.main.meta.id
  }

  /* 
  have docs support multiple paths that get combined into docs?
  grabs all the books from the library and generates navigation
  mounts the library to a volume
  */

  // logo {
  //   url = "https://companieslogo.com/img/orig/HCP.D-be08ca6f.png"
  //   width = 32
  //   height = 32
  // }

  content = [
    module.workshop.output.book
  ]

  assets = "./workshop/images"
}