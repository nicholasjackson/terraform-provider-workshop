---
sidebar_position: 12
id: data_source_4
title: Manually Testing the Data Source
---

First you need to build and install your provider, run the following
command in your terminal.

```shell
make build && make install
```

Once that is done change into the directory where you created the example
Terraform config.

```
cd ./examples/data-sources/minecraft_block
```

## Initializing the Config

Like all Terraform configuraion it needs to be initalized so Terraform
fetches the provider from the registy or in your case, the local file system.

```shell
terraform init
```

You should see some output that looks like the following.

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

## Running an Terraform Apply

Let's now run an apply and see the output, you will see the output variable
is set to the Minecraft material representing the block at that location.

```shell
terraform apply
```

```shell
data.minecraft_block.example: Reading...
data.minecraft_block.example: Read complete after 0s [id=LTEyNzMvMjMvMjg4]

Changes to Outputs:
  + example = "minecraft:infested_cracked_stone_bricks"

You can apply this plan to save these new output values to the Terraform state, without changing any
real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes
```

```shell
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

example = "minecraft:infested_cracked_stone_bricks"
```

## Summary

That is all for now, in this workshop you have learned the basics of Terraform
provider development using the new Terraform Plugin SDK. While we have only
implemented a single resource and a single data source, there is no limit to the 
number of resources and data sources that your provider can contain. You only need
to create a stuct that implements the correct interface in exactly the same way
that you have seen today.

If you are feeling curious, why not have a go at adding another resource to
implement the `block` API.

[http://localhost:9090/redoc#tag/Block](http://localhost:9090/redoc#tag/Block)