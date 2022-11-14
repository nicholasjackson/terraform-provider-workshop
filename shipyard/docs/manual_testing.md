---
sidebar_position: 8
id: manual_testing
title: Manually Testing the Provider
---

Before testing the provider you need to build it and install it to your local
folder. Run the following commands in your terminal.

```shell
make build
make install
```

```shell
go build -o bin/terraform-provider-example_v0.1.0
rm -rf examples/resources/block/.terraform*
rm -rf examples/resources/block/terraform.tfstate*
mkdir -p ~/.terraform.d/plugins/local/hashicraft/example/0.1.0/darwin_amd64
mv bin/terraform-provider-example_v0.1.0 ~/.terraform.d/plugins/local/hashicraft/example/0.1.0/darwin_amd64/
```

## Creating the Example

In the `examples` folder there is a folder called `scaffolding_example` 
change the name of this folder to `minecraft_schema`. Then in the `resource.tf`
file delete any existing contents.

You first need to define a required providers stanza that tells Terraform
where to find your provider.

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

Next you can configure the provider, the address `minecraft.container.shipyard.run` will resolve to the local Minecraft server that the environment is running.

```javascript
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
```

Finally let's specify the resource

```javascript
resource "minecraft_schema" "bus" {
  x = -1278
  y = 24
  z = 288
  rotation = 270
  schema = "../../../schemas/car.zip"
}
```

Your final file should look like the following.

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

resource "minecraft_schema" "bus" {
  x = -1278
  y = 24
  z = 288
  rotation = 270
  schema = "../../../schemas/car.zip"
}
```

## Init the Terraform Config

To use a provider with Terraform you first need to run `terraform init`,
Terraform will locate any required providers and copy them to the local
folder.

In your terminal change directory to the `examples/minecraft_schema` folder
and run the following command.

```shell
terraform init
```

You will see some output that looks like the following

```shell
Initializing the backend...

Initializing provider plugins...
- Finding local/hashicraft/minecraft versions matching "0.1.0"...
- Installing local/hashicraft/minecraft v0.1.0...
- Installed local/hashicraft/minecraft v0.1.0 (unauthenticated)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

## Running a Plan

Let's now check the provider by running a `terraform plan`. Run the following 
command in your terminal.

```shell
terraform plan
```

You should see output like the following:

```shell
Terraform used the selected providers to generate the following execution plan. Resource
actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # minecraft_schema.bus will be created
  + resource "minecraft_schema" "bus" {
      + id          = (known after apply)
      + rotation    = 270
      + schema      = "../../../schemas/car.zip"
      + schema_hash = (known after apply)
      + x           = -1278
      + y           = 24
      + z           = 288
    }

Plan: 1 to add, 0 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to
take exactly these actions if you run "terraform apply" now.
```

Let's now create the resources.

## Running Apply

Let's now run `terraform apply` to create the resources.