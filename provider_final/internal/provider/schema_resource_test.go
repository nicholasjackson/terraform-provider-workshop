package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSchemaResourceConfig(1, 2, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("minecraft_schema.car", "x", "1"),
					resource.TestCheckResourceAttr("minecraft_schema.car", "y", "2"),
					resource.TestCheckResourceAttr("minecraft_schema.car", "z", "3"),
				),
			},
		},
	})
}

func testAccSchemaResourceConfig(x, y, z int) string {
	return fmt.Sprintf(`
  resource "minecraft_schema" "car" {
	  x = %d
	  y = %d
	  z = %d
	  rotation = 270
	  schema = "../../schemas/car.zip"
	}
  `, x, y, z)
}
