# Paths
CATALOG_SPEC=pkg/api/openapi/polaris-catalog-service.yaml
MGMT_SPEC=pkg/api/openapi/polaris-management-service.yaml

CATALOG_CONF=pkg/api/openapi/catalog-config.yaml
MGMT_CONF=pkg/api/openapi/management-config.yaml

.PHONY: generate
generate:
	@echo "Generating Catalog API..."
	docker run --rm -v $(shell pwd):/host -w /host \
		deepmap/oapi-codegen:v2.1.0 \
		--config /host/$(CATALOG_CONF) /host/$(CATALOG_SPEC)

	@echo "Generating Management API..."
	docker run --rm -v $(shell pwd):/host -w /host \
		deepmap/oapi-codegen:v2.1.0 \
		--config /host/$(MGMT_CONF) /host/$(MGMT_SPEC)