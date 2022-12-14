---
sidebar_position: 4
id: schema_2
title: Schema Resource - Read Method
---

The Read method is used by Terraform to check the status of an existing 
resource. It is also used to synchronise computed fields from an API.
An example of this might be a token such as the Kubernetes config used
to connect to a server that has been created with a Terraform resource.

If the token changes for example due to a Time to Live constraint, Terraform
can automatically reload this token through the Read method.

We are going to implement a very basic usecase for `Read`, we are only going
to check if the schema exists on the server. If it does not then we will 
return an error as it is likely that the resource has been deleted outside
of Terraform.

To perform a read, you can use the client SDK method `getSchemaDetails`.
This method takes a single parameter which is the ID of the schema, the same
ID that was retured when creating a schema and you saved into the state
when writing the `Create` method.

If you take a look at the `Read` method you will see that the top of the
method has a very similar approach to the `Create` method.

The main difference is that Read uses the `State` from the request where
Create used the `Plan`. The `State` type is the state that was saved at the
end of the Create method.

```go
var data *SchemaResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

if resp.Diagnostics.HasError() {
	return
}
```

Earlier it was mentioned that to cast between SDK types and Go types you use
methods on the type such as `ValueString` as unlike Go types SDK types are 
references and may be nil. Read is called before any other method in your
provider so it is possible that the Id parameter of your `SchemaResourceModel`
is nil. You need to defensively code against this and can use the `IsNull`
method to check if there has been a value set before attempting to cast it to
a Go type.

Add the following code to the end of the Read method.

```go
if data.Id.IsNull() {
	return
}

_, err := r.minecraftClient.getSchemaDetails(data.Id.ValueString())
if err != nil {
	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read schema, got error: %s", err))
	return
}
```

## Checking the Read Method

Before progressing, Let's just make sure that the code compiles and is all 
ok. Run the following command in your terminal.

```shell
make build
```

```shell
go build -o bin/terraform-provider-example_v0.1.0
```

All being well your code will compile and you can move on to the next method 
for a resource `Modify`.