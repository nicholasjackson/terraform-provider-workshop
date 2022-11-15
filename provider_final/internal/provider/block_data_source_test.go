package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.minecraft_block.example", "x", "-1272"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
data "minecraft_block" "example" {
  x = -1272
  y = 23
  z = 288
}
`
