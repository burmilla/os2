.DEFAULT_GOAL := package
REPO?=rancher/os2
TAG?=dev
IMAGE=${REPO}:${TAG}

.dapper:
	@echo Downloading dapper
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-$$(uname -s)-$$(uname -m) > .dapper.tmp
	@@chmod +x .dapper.tmp
	@./.dapper.tmp -v
	@mv .dapper.tmp .dapper

.PHONY: ci
ci: .dapper
	./.dapper ci

.PHONY: package
package: .dapper
	./.dapper package

.PHONY: clean
clean:
	rm -rf build dist

.PHONY: build-framework
build-framework:
	docker build \
		--build-arg CACHEBUST=${CACHEBUST} \
		--build-arg IMAGE_TAG=${TAG} \
		--build-arg IMAGE_REPO=${REPO}-framework \
		-t ${REPO}-framework:${TAG} .

.PHONY: build
build:
	docker build \
		--build-arg CACHEBUST=${CACHEBUST} \
		--build-arg IMAGE_TAG=${TAG} \
		--build-arg IMAGE_REPO=${REPO} \
		-t ${IMAGE} .

.PHONY: push
push:
	docker push ${IMAGE}

.PHONY: push
push-framework: build-framework
	docker push ${REPO}-framework:${TAG}

.PHONY: iso
iso:
	./ros-image-build ${IMAGE} iso
	@echo "INFO: ISO available at build/output.iso"

.PHONY: qcow
qcow:
	./ros-image-build ${IMAGE} qcow
	@echo "INFO: QCOW image available at build/output.qcow.gz"

.PHONY: ami-%
ami-%:
	AWS_DEFAULT_REGION=$* ./ros-image-build ${IMAGE} ami

.PHONY: ami
ami:
	./ros-image-build ${IMAGE} ami

.PHONY: run
run:
	./scripts/run

.PHONY: run
pxe:
	./scripts/run pxe

serve-docs: mkdocs
	docker run -p 8000:8000 --rm -it -v $${PWD}:/docs mkdocs serve -a 0.0.0.0:8000

mkdocs:
	docker build -t mkdocs -f Dockerfile.docs .

all-amis: \
	ami-us-west-1 \
	ami-us-west-2
	#ami-ap-east-1 \
	#ami-ap-northeast-1 \
	#ami-ap-northeast-2 \
	#ami-ap-northeast-3 \
	#ami-ap-southeast-1 \
	#ami-ap-southeast-2 \
	#ami-ca-central-1 \
	#ami-eu-central-1 \
	#ami-eu-south-1 \
	#ami-eu-west-1 \
	#ami-eu-west-2 \
	#ami-eu-west-3 \
	#ami-me-south-1 \
	#ami-sa-east-1 \
	#ami-us-east-1 \
	#ami-us-east-2 \
