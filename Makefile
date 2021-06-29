HOSTNAME=github.com
NAMESPACE=Ouest-France
NAME=ldap
BINARY=terraform-provider-${NAME}
VERSION=0.8.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

generate-docs:
	tfplugindocs
