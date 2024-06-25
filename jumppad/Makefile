REPO=nicholasjackson/terraform-provider-workshop
VERSION=v0.3.0

build_codeserver:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx create --name vscode || true
	docker buildx use vscode
	docker buildx inspect --bootstrap
	docker buildx build --platform linux/arm64,linux/amd64 \
		-t ${REPO}:${VERSION} \
	  -f ./dockerfiles/codeserver/Dockerfile \
		--no-cache \
	  . \
		--push
	
