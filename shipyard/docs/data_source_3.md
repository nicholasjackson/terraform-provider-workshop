---
sidebar_position: 11
id: data_source_3
title: Data Source - Creating the Config Example
---

Like in the reosource example, you need to define a `required_providers` 
example that tells terraform where to find your installed provider.

```javascript
terraform {
  required_providers {
    minecraft = {
      source  = "local/hashicraft/minecraft"
      version = "0.1.0"
    }
  }
}
```

Then you can configure the provider with the endpoint and the api_key.

```javascript
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
```

And finally the data source, to test that things are working, let's also define 
an output variable that returns the material for the block data source.

```javascript
data "minecraft_block" "example" {
  x = -1273
  y = 23
  z = 288
}

output "example" {
  value = data.minecraft_block.example.material
}
```

The full output will look like the following block, add this to your `data-source.tf`
file.

```javascript
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
```

With all that done it is now time to test the provider.