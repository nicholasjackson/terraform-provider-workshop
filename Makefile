REPO=nicholasjackson/terraform-provider-workshop
VERSION=v0.3.1

build_codeserver:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx create --name vscode || true
	docker buildx use vscode
	docker buildx inspect --bootstrap
	docker buildx build --platform linux/arm64,linux/amd64 \
		-t ${REPO}:${VERSION} \
	  -f ./jumppad/dockerfiles/codeserver/Dockerfile \
		--no-cache \
	  ./jumppad \
		--push
	
functional_test:
	dagger -m ./dagger call functional-test --src . --working-directory jumppad --runtime docker

build_packer:
	cd ./jumppad/packer && packer build -var-file=./main.pkrvars.hcl ./main.pkr.hcl