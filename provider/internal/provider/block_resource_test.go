package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBlockResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("minecraft_block.one", "x", "-1273"),
					resource.TestCheckResourceAttr("minecraft_block.one", "y", "25"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBlockResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}

resource "minecraft_block" "%s" {
  x = -1273
  y = 25
  z = 288
  material = "minecraft:stone"
}
`, configurableAttribute)
}
