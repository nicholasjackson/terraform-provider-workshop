default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

arch=$(shell /bin/bash ./arch.sh)
name = minecraft
organization = hashicraft
version = 0.1.0
log_level = info

build:
	@echo "\033[1mBuilding the provider\033[0m"
	go build -o bin/terraform-provider-$(name)_v$(version) -buildvcs=false

	@echo ""

install: build clean
	@echo "\033[1mInstalling the provider\033[0m"
	mkdir -p ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)
	mv bin/terraform-provider-$(name)_v$(version) ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)/

	@echo ""

clean:
	@echo "\033[1mCleaning up\033[0m"
	rm -rf examples/resources/minecraft_schema/.terraform*
	rm -rf examples/resources/minecraft_schema/terraform.tfstate*
	rm -rf examples/data-sources/minecraft_block/.terraform*
	rm -rf examples/data-sources/minecraft_block/terraform.tfstate*

	@echo ""

init:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/schema init

plan:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/schema plan

apply:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/schema apply -auto-approve

destroy:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/schema destroy -auto-approve

default: testacc

prepare-release:
	rm -rf ./release
	mkdir -p ./release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-$(name)_v$(version) -buildvcs=false
	cd bin && zip ../release/terraform-provider-$(name)_$(version)_linux_amd64.zip terraform-provider-$(name)_v$(version)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/terraform-provider-$(name)_v$(version) -buildvcs=false
	cd bin && zip ../release/terraform-provider-$(name)_$(version)_linux_arm64.zip terraform-provider-$(name)_v$(version)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-$(name)_v$(version) -buildvcs=false
	cd bin && zip ../release/terraform-provider-$(name)_$(version)_darwin_amd64.zip terraform-provider-$(name)_v$(version)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/terraform-provider-$(name)_v$(version) -buildvcs=false
	cd bin && zip ../release/terraform-provider-$(name)_$(version)_darwin_arm64.zip terraform-provider-$(name)_v$(version)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-$(name)_v$(version).exe -buildvcs=false
	cd bin && zip ../release/terraform-provider-$(name)_$(version)_windows_amd64.zip terraform-provider-$(name)_v$(version).exe
	cp ./terraform-registry-manifest.json ./release/terraform-provider-$(name)_$(version)_manifest.json
	cd release && shasum -a 256 *.zip > terraform-provider-$(name)_$(version)_SHA256SUMS
	cd release && shasum -a 256 *.json >> terraform-provider-$(name)_$(version)_SHA256SUMS
	cd release && gpg --detach-sign terraform-provider-$(name)_$(version)_SHA256SUMS