package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSchemaResourceConfig("car"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("minecraft_schema.car", "x", "-1278"),
					resource.TestCheckResourceAttr("minecraft_schema.car", "y", "24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSchemaResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}

resource "minecraft_schema" "%s" {
  x = -1278
  y = 24
  z = 288
  rotation = 270
  schema = "../../example_schemas/car.zip"
}
`, configurableAttribute)
}
