# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= 0.0.1
API_VERSION ?= v1alpha1

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "candidate,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=candidate,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="candidate,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for bundle and catalog images.
#
# For example, running 'make bundle-build bundle-push catalog-build catalog-push' will build and push both
# flux-framework.org/operator-bundle:$VERSION and flux-framework.org/operator-catalog:$VERSION.
IMAGE_TAG_BASE ?= ghcr.io/flux-framework/flux-operator
KIND_VERSION=v0.11.1
# This kubectl version supports -k for kustomization, taken from mpi
KUBECTL_VERSION=v1.21.4
BUILDENVVAR=CGO_CFLAGS="-I/usr/include" CGO_LDFLAGS="-L/usr/lib -lstdc++ -lczmq -lzmq"

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)
IMG_BUILDER=docker

# BUNDLE_GEN_FLAGS are the flags passed to the operator-sdk generate bundle command
BUNDLE_GEN_FLAGS ?= -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)

# USE_IMAGE_DIGESTS defines if images are resolved via tags or digests
# You can enable this value if you would like to use SHA Based Digests
# To enable set flag to true
USE_IMAGE_DIGESTS ?= false
ifeq ($(USE_IMAGE_DIGESTS), true)
	BUNDLE_GEN_FLAGS += --use-image-digests
endif

# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/flux-framework/flux-operator

# Testing image (for development mostly)
DEVIMG ?= ghcr.io/flux-framework/flux-operator:test
ARMIMG ?= ghcr.io/flux-framework/flux-operator:arm

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.24.1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

HELMIFY ?= $(LOCALBIN)/helmify


# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen openapi-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	${OPENAPI_GEN} --logtostderr=true -i ./api/${API_VERSION}/ -o "" -O zz_generated.openapi -p ./api/${API_VERSION}/ -h ./hack/boilerplate.go.txt -r "-"

.PHONY: api
api: generate api
	go run hack/python-sdk/main.go ${API_VERSION} > ${SWAGGER_API_JSON}
	rm -rf ./sdk/python/${API_VERSION}/fluxoperator/model/*
	rm -rf ./sdk/python/${API_VERSION}/fluxoperator/test/test_*.py
	java -jar ${SWAGGER_JAR} generate -i ${SWAGGER_API_JSON} -g python-legacy -o ./sdk/python/${API_VERSION} -c ./hack/python-sdk/swagger_config.json --git-repo-id flux-operator --git-user-id flux-framework
	cp ./hack/python-sdk/template/* ./sdk/python/${API_VERSION}/

# These were needed for the python (not python-legacy)
# cp ./hack/python-sdk/fluxoperator/* ./sdk/python/${API_VERSION}/fluxoperator/model/

.PHONY: helmify
helmify: $(HELMIFY) ## Download helmify locally if necessary.
$(HELMIFY): $(LOCALBIN)
	test -s $(LOCALBIN)/helmify || GOBIN=$(LOCALBIN) go install github.com/arttor/helmify/cmd/helmify@latest
    
helm: manifests kustomize helmify
	$(KUSTOMIZE) build config/default | $(HELMIFY)

.PHONY: pre-push
pre-push: generate api build-config-arm build-config helm
	git status

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out

.PHONY: list
list:
	kubectl get -n flux-operator pods

.PHONY: reset
reset:
	minikube stop
	minikube delete
	minikube start
	kubectl create namespace flux-operator
	make install
	make redo

.PHONY: clean
clean:
	kubectl delete -n flux-operator svc --all --grace-period=0 --force
	# kubectl delete -n flux-operator secret --all --grace-period=0 --force
	kubectl delete -n flux-operator cm --all --grace-period=0 --force
	# pods, pvc, and pv need to be deleted in this order
	kubectl delete -n flux-operator pods --all --grace-period=0 --force
	kubectl delete -n flux-operator pvc --all --grace-period=0 --force
	kubectl delete -n flux-operator pv --all --grace-period=0 --force
	kubectl delete -n flux-operator jobs --all --grace-period=0 --force
	kubectl delete -n flux-operator MiniCluster --all --grace-period=0 --force

# This applies the basic minicluster (and not extended examples)
apply:
	kubectl apply -f examples/flux-restful/minicluster-lammps.yaml

applyui:
	kubectl apply -f examples/flux-restful/minicluster-$(name).yaml

applytest:
	kubectl apply -f examples/tests/${name}/minicluster-$(name).yaml

example:
	kubectl apply -f examples/flux-restful/minicluster-$(name).yaml

# Clean, apply and run, and apply the job
redo: clean apply run
redo_example: clean example run
redo_test: clean applytest run

log:
	kubectl logs -n flux-operator job.batch/flux-sample $@

##@ Test
# NOTE these are not fully developed yet

bin/kubectl:
	mkdir -p bin
	curl -L -o bin/kubectl https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
	chmod +x bin/kubectl

.PHONY: test_e2e
test_e2e: export TEST_FLUX_OPERATOR_IMAGE = ${IMAGE_TAG_BASE}:latest
test_e2e: bin/kubectl kind images dev_manifest
	go test -tags e2e ./tests/e2e/...

.PHONY: dev_manifest
dev_manifest:
	# Use `~` instead of `/` because image name might contain `/`.
	sed -e "s~%IMAGE_NAME%~${IMAGE_TAG_BASE}~g" -e "s~%IMAGE_TAG%~${VERSION}~g" config/manifests/overlays/dev/kustomization.yaml.template > config/manifests/overlays/dev/kustomization.yaml

.PHONY: kind
kind:
	go install sigs.k8s.io/kind@${KIND_VERSION}


# TODO add build arg for version
.PHONY: images
images:
	@echo "VERSION: ${VERSION}"
	${IMG_BUILDER} build -t ${IMAGE_TAG_BASE}:local .

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	mv ./controllers/flux/keygen.go ./controllers/flux/keygen.go.backup
	cp ./controllers/flux/keygen.go.template ./controllers/flux/keygen.go
	$(BUILDENVVAR) go build -o bin/manager main.go
	mv ./controllers/flux/keygen.go.backup ./controllers/flux/keygen.go

.PHONY: build-container
build-container: generate fmt vet
	cp ./controllers/flux/keygen.go.template ./controllers/flux/keygen.go
	$(BUILDENVVAR) go build -a -o ./manager main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: arm-build
arm-build: test ## Build docker image with the manager.
	docker buildx build --platform linux/arm64 -t ${ARMIMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: build-config
build-config: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > examples/dist/flux-operator.yaml

.PHONY: build-config-arm
build-config-arm: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${ARMIMG}
	$(KUSTOMIZE) build config/default > examples/dist/flux-operator-arm.yaml

# Build a test image, push to the registry at test, and apply the build-config
.PHONY: test-deploy
test-deploy: manifests kustomize
	docker build --no-cache -t ${DEVIMG} .
	docker push ${DEVIMG}
	cd config/manager && $(KUSTOMIZE) edit set image controller=${DEVIMG}
	$(KUSTOMIZE) build config/default > examples/dist/flux-operator-dev.yaml

.PHONY: arm-deploy
arm-deploy: manifests kustomize
	docker buildx build --platform linux/arm64 --push -t ${ARMIMG} .
	cd config/manager && $(KUSTOMIZE) edit set image controller=${ARMIMG}
	$(KUSTOMIZE) build config/default > examples/dist/flux-operator-arm.yaml

# Build a local test image, load into minikube or kind and apply the build-config
.PHONY: deploy-local
deploy-local: manifests kustomize build
	kubectl delete -f examples/dist/flux-operator-local.yaml || true
	docker build -t ${DEVIMG} .
	cd config/manager && $(KUSTOMIZE) edit set image controller=${DEVIMG}
	$(KUSTOMIZE) build config/default > examples/dist/flux-operator-local.yaml
	sed -i 's/        imagePullPolicy: Always/        imagePullPolicy: Never/' examples/dist/flux-operator-local.yaml

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
INSTALL_KUSTOMIZE ?= $(LOCALBIN)/install_kustomize.sh
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
OPENAPI_GEN ?= $(LOCALBIN)/openapi-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
SWAGGER_JAR ?= ${LOCALBIN}/openapi-generator-cli.jar
SWAGGER_API_JSON ?= ./api/${API_VERSION}/swagger.json

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.9.0

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	bash $(INSTALL_KUSTOMIZE) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

# Build the latest openapi-gen from source
.PHONY: openapi-gen
openapi-gen: $(OPENAPI_GEN) ## Download controller-gen locally if necessary.
$(OPENAPI_GEN): $(LOCALBIN)
	which ${OPENAPI_GEN} > /dev/null || (git clone --depth 1 https://github.com/kubernetes/kube-openapi /tmp/kube-openapi && cd /tmp/kube-openapi && go build -o ${OPENAPI_GEN} ./cmd/openapi-gen)

.PHONY: swagger-jar
swagger-jar: $(SWAGGER_JAR) ## Download controller-gen locally if necessary.
$(SWAGGER_JAR): $(LOCALBIN)
	wget -qO ${SWAGGER_JAR} "https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/5.1.0/openapi-generator-cli-5.1.0.jar"

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: bundle
bundle: manifests kustomize ## Generate bundle manifests and metadata, then validate generated files.
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle $(BUNDLE_GEN_FLAGS)
	operator-sdk bundle validate ./bundle

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push: ## Push the bundle image.
	$(MAKE) docker-push IMG=$(BUNDLE_IMG)

.PHONY: opm
OPM = ./bin/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (,$(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.23.0/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif

# A comma-separated list of bundle images (e.g. make catalog-build BUNDLE_IMGS=example.com/operator-bundle:v0.1.0,example.com/operator-bundle:v0.2.0).
# These images MUST exist in a registry and be pull-able.
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# The image tag given to the resulting catalog image (e.g. make catalog-build CATALOG_IMG=example.com/operator-catalog:v0.2.0).
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:v$(VERSION)

# Set CATALOG_BASE_IMG to an existing catalog image tag to add $BUNDLE_IMGS to that image.
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Build a catalog image by adding bundle images to an empty catalog using the operator package manager tool, 'opm'.
# This recipe invokes 'opm' in 'semver' bundle add mode. For more information on add modes, see:
# https://github.com/operator-framework/community-operators/blob/7f1438c/docs/packaging-operator.md#updating-your-existing-operator
.PHONY: catalog-build
catalog-build: opm ## Build a catalog image.
	$(OPM) index add --container-tool docker --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

# Push the catalog image.
.PHONY: catalog-push
catalog-push: ## Push a catalog image.
	$(MAKE) docker-push IMG=$(CATALOG_IMG)
