# Paths
CATALOG_SPEC=pkg/api/openapi/polaris-catalog-service.yaml
MGMT_SPEC=pkg/api/openapi/polaris-management-service.yaml

CATALOG_CONF=pkg/api/openapi/catalog-config.yaml
MGMT_CONF=pkg/api/openapi/management-config.yaml

# The Tool Version (Easy to update)
CODEGEN_VERSION=v2.4.1

.PHONY: generate
generate:
	@echo "Generating Catalog API..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(CODEGEN_VERSION) \
		--config $(CATALOG_CONF) $(CATALOG_SPEC)

	@echo "Generating Management API..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(CODEGEN_VERSION) \
		--config $(MGMT_CONF) $(MGMT_SPEC)