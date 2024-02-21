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
  default     = "localhost"
}

variable "api_url" {
  description = "The URL for the Minecraft API"
  default     = "http://localhost:9090"
}

resource "chapter" "introduction" {
  title = "Introduction"

  page "introduction" {
    content = template_file("docs/introduction/intro.mdx", {
      docs_url       = variable.docs_url
      prismarine_url = variable.prismarine_url
      minecraft_url  = variable.minecraft_url
      api_url        = variable.api_url
    })
  }
}

resource "chapter" "resources" {
  title = "Resources"

  page "overview" {
    content = file("docs/resources/overview.mdx")
  }

  page "schema_create" {
    content = template_file("docs/resources/schema_resource_1.mdx", {
      docs_url       = variable.docs_url
      prismarine_url = variable.prismarine_url
      minecraft_url  = variable.minecraft_url
      api_url        = variable.api_url
    })
  }

  page "schema_custom" {
    content = file("docs/resources/schema_resource_2.mdx")
  }

  page "schema_read" {
    content = file("docs/resources/schema_resource_3.mdx")
  }

  page "schema_update" {
    content = file("docs/resources/schema_resource_4.mdx")
  }

  page "schema_delete" {
    content = file("docs/resources/schema_resource_5.mdx")
  }

  page "provider_configure" {
    content = file("docs/resources/provider_configure.mdx")
  }

  page "manual_testing" {
    content = file("docs/resources/manual_testing.mdx")
  }
}

resource "chapter" "data_sources" {
  title = "Data Sources"

  page "creating" {
    content = file("docs/data_sources/data_source_1.mdx")
  }

  page "read" {
    content = file("docs/data_sources/data_source_2.mdx")
  }

  page "config" {
    content = file("docs/data_sources/data_source_3.mdx")
  }

  page "manual_testing" {
    content = file("docs/data_sources/manual_testing.mdx")
  }
}

resource "book" "terraform_provider" {
  title = "Building a Terraform Provider"

  chapters = [
    resource.chapter.introduction,
    resource.chapter.resources,
    resource.chapter.data_sources,
  ]
}

output "book" {
  value = resource.book.terraform_provider
}