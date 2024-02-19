resource "chapter" "introduction" {
  title = "Introduction"

  page "introduction" {
    content = file("docs/introduction/intro.mdx")
  }
}

resource "chapter" "resources" {
  title = "Resources"

  page "overview" {
    content = file("docs/resources/overview.mdx")
  }

  page "schema_create" {
    content = file("docs/resources/schema_resource_1.mdx")
  }

  page "schema_read" {
    content = file("docs/resources/schema_resource_2.mdx")
  }

  page "schema_update" {
    content = file("docs/resources/schema_resource_3.mdx")
  }

  page "schema_delete" {
    content = file("docs/resources/schema_resource_4.mdx")
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