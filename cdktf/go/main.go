package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-digitalocean-go/digitalocean/v2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	name := "testing"
	region := "ams"
	stack := cdktf.NewTerraformStack(scope, &id)

	// The code that defines your stack goes here
	digitalocean.NewDigitaloceanProvider(stack, jsii.String("digitalocean"), &digitalocean.DigitaloceanProviderConfig{})

	digitalocean.NewApp(stack, jsii.String("static_site_example"), &digitalocean.AppConfig{
		Spec: &digitalocean.AppSpec{
			Name:   jsii.String(fmt.Sprintf("static-site-%s", name)),
			Region: jsii.String(region),
			StaticSite: []*digitalocean.AppSpecStaticSite{
				{
					Name:      jsii.String(name),
					SourceDir: jsii.String("/src"),

					Github: &digitalocean.AppSpecStaticSiteGithub{
						Repo:         jsii.String("nicholasjackson/mame-wasm"),
						DeployOnPush: jsii.Bool(true),
						Branch:       jsii.String("main"),
					},
				},
			},
		},
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "go")

	app.Synth()
}
