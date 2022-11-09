FROM golang:latest as golang

FROM ubuntu

# Install default packages
RUN apt update && apt upgrade -y && apt install -y curl wget zip git jq 

# Install Go
COPY --from=golang /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

# Install nodejs
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
  apt install -y nodejs

# Clean packages
RUN rm -rf /var/lib/apt/lists/*

# Install the Terraform CDK
RUN npm install --global cdktf-cli@latest

# Install Terraform
RUN wget -O terraform.zip https://releases.hashicorp.com/terraform/1.3.4/terraform_1.3.4_linux_amd64.zip && \
  unzip terraform.zip && \
  mv terraform /usr/local/bin

# Install code server
RUN wget -O- https://aka.ms/install-vscode-server/setup.sh | sh

RUN code-server serve-local --accept-server-license-terms --install-extension shipyard.shipyard
RUN code-server serve-local --accept-server-license-terms --install-extension golang.go
RUN code-server serve-local --accept-server-license-terms --install-extension hashicorp.hcl
RUN code-server serve-local --accept-server-license-terms --install-extension hashicorp.terraform
RUN code-server serve-local --accept-server-license-terms --install-extension pkief.material-icon-theme
RUN code-server serve-local --accept-server-license-terms --install-extension whizkydee.material-palenight-theme

ENTRYPOINT ["code-server", "serve-local","--accept-server-license-terms","--without-connection-token","--host","0.0.0.0","--verbose"]
