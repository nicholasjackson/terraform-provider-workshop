default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

arch=$(shell /bin/bash ./arch.sh)
name = example
organization = hashicraft
version = 0.1.0
log_level = info

build:
	go build -o bin/terraform-provider-$(name)_v$(version)

install: build clean
	mkdir -p ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)
	mv bin/terraform-provider-$(name)_v$(version) ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)/

clean:
	rm -rf examples/resources/block/.terraform*
	rm -rf examples/resources/block/terraform.tfstate*

init:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/block init

plan:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/block plan

apply:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/block apply -auto-approve

destroy:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/block destroy -auto-approve

default: testacc