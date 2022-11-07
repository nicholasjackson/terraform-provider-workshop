FROM golang:latest as golang

FROM ubuntu

RUN apt update && apt upgrade -y && apt install -y curl wget zip git jq && rm -rf /var/lib/apt/lists/*

RUN wget -O- https://aka.ms/install-vscode-server/setup.sh | sh

RUN code-server serve-local --accept-server-license-terms --install-extension shipyard.shipyard

RUN wget -O terraform.zip https://releases.hashicorp.com/terraform/1.3.4/terraform_1.3.4_linux_amd64.zip && \
  unzip terraform.zip && \
  mv terraform /usr/local/bin

COPY --from=golang /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

ENTRYPOINT ["code-server", "serve-local","--accept-server-license-terms","--without-connection-token","--host","0.0.0.0","--verbose"]
