FROM golang:1.22.4 as golang

FROM ghcr.io/jumppad-labs/vscode:vscode-v0.5.1

ARG TARGETARCH

# Install default packages
RUN apt-get update && apt-get upgrade -y && apt-get install -y curl wget zip git jq build-essential

# Clean packages
RUN rm -rf /var/lib/apt/lists/*

# Install Go
COPY --from=golang /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

RUN wget -O terraform.zip https://releases.hashicorp.com/terraform/1.8.5/terraform_1.8.5_linux_${TARGETARCH}.zip && \
  unzip terraform.zip && \
  mv terraform /usr/local/bin

# Install VSCode extensions
RUN code-server --install-extension hashicorp.terraform & \
  code-server --install-extension hashicorp.hcl & \
  code-server --install-extension golang.go

# Install Go tools
RUN go install honnef.co/go/tools/cmd/staticcheck@latest & \
  go install golang.org/x/tools/gopls@latest & \
  go install github.com/cweill/gotests/gotests@v1.6.0 & \
  go install github.com/fatih/gomodifytags@v1.16.0 & \
  go install github.com/josharian/impl@v1.1.0 & \
  go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 & \
  go install github.com/go-delve/delve/cmd/dlv@latest