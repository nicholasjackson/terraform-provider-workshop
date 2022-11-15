---
sidebar_position: 10
id: data_source_2
title: Data Source - Implementing the Read Method
---

Where resources have a full CRUD implementation with `Create`, `Read`, `Update`
and `Delete` methods, a data source only has a `Read` method.  Creating 
this method is very similar to how you implemented the `Read` method in the
resource. Hopefully implementing `Read` in the data source should be a 
little bit familliar.

To get the details specified in the data soruce stanza you can deserialzie
these into your model.

```go
var data BlockDataSourceModel

// Read Terraform configuration data into the model
resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

if resp.Diagnostics.HasError() {
	return
}
```

Then like the resource, you need to convert them into native Go types so that
you can use them with the API.

```go
x,_ := data.X.ValueBigFloat().Int64()
y,_ := data.Y.ValueBigFloat().Int64()
z,_ := data.Z.ValueBigFloat().Int64()
```

Next let's call the API using the clients `getBlock` method, this takes
3 parameters which are all integers, you can populate the parameters from
the attributes you deserialized.

`getBlock` will return an error if a block is not found or an internal error
has occured. Let's return this error using the `diagnostics` so that Terraform
can display it to the end user. If there is an error at this point you need
to return immediately.

```go
block, err := d.minecraftClient.getBlock(int(x), int(y), int(z))
if err != nil {
  resp.Diagnostics.AddError( 
    "Unable to retrieve block",
    fmt.Sprintf("Unable to get block, got error: %s", err),
  )
  return
}
```

Finally you can set the block `Id` that is returned from the API and set this
into the state.

```go
data.Id = types.StringValue(block.ID)

// Write logs using the tflog package
// Documentation: https://terraform.io/plugin/log
tflog.Trace(ctx, "read a data source")

// Save data into Terraform state
resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
```

The complete method looks like the following example, replace the `Read` method
in your `block_data_source.go` with this code.

```go
func (d *BlockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BlockDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

  x,_ := data.X.ValueBigFloat().Int64()
  y,_ := data.Y.ValueBigFloat().Int64()
  z,_ := data.Z.ValueBigFloat().Int64()

  block, err := d.minecraftClient.getBlock(int(x), int(y), int(z))
  if err != nil {
    resp.Diagnostics.AddError( 
      "Unable to retrieve block",
      fmt.Sprintf("Unable to get block, got error: %s", err),
    )

    return
  }
 
  data.Material = types.StringValue(block.Material)
  data.Id = types.StringValue(block.ID)
  
  // Write logs using the tflog package
  // Documentation: https://terraform.io/plugin/log
  tflog.Trace(ctx, "read a data source")
  // Save data into Terraform state
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

## Testing the Provider Compiles

Let's test that the provider code compiles and you can write an example
configuration that uses the new data source.