# Polaris CLI Implementation Status

This file tracks the implementation status of CLI commands mapping to the Polaris Management and Catalog APIs.

## Management API

### Catalogs
Files: `cmd/catalogs.go`
- [x] List (`catalogs list`) - `ListCatalogs`
- [x] Create (`catalogs create`) - `CreateCatalog`
- [x] Describe (`catalogs describe`) - `GetCatalog`
- [x] Delete (`catalogs delete`) - `DeleteCatalog`
- [ ] Update (`catalogs update`) - `UpdateCatalog`

### Catalog Roles
Files: `cmd/catalog_roles.go` (Proposed)
- [ ] List - `ListCatalogRoles`
- [ ] Create - `CreateCatalogRole`
- [ ] Delete - `DeleteCatalogRole`
- [ ] Describe - `GetCatalogRole`
- [ ] Update - `UpdateCatalogRole`
- [ ] List Grants - `ListGrantsForCatalogRole`
- [ ] Revoke Grant - `RevokeGrantFromCatalogRole`
- [ ] Add Grant - `AddGrantToCatalogRole`
- [ ] List Assignee Principal Roles - `ListAssigneePrincipalRolesForCatalogRole`

### Principals
Files: `cmd/principals.go` (Proposed)
- [ ] List - `ListPrincipals`
- [ ] Create - `CreatePrincipal`
- [ ] Delete - `DeletePrincipal`
- [ ] Describe - `GetPrincipal`
- [ ] Update - `UpdatePrincipal`
- [ ] Rotate Credentials - `RotateCredentials`
- [ ] Reset Credentials - `ResetCredentials`
- [ ] List Assigned Roles - `ListPrincipalRolesAssigned`
- [ ] Assign Role - `AssignPrincipalRole`
- [ ] Revoke Role - `RevokePrincipalRole`

### Principal Roles
Files: `cmd/principal_roles.go` (Proposed)
- [ ] List - `ListPrincipalRoles`
- [ ] Create - `CreatePrincipalRole`
- [ ] Delete - `DeletePrincipalRole`
- [ ] Describe - `GetPrincipalRole`
- [ ] Update - `UpdatePrincipalRole`
- [ ] List Catalog Roles - `ListCatalogRolesForPrincipalRole`
- [ ] Assign Catalog Role - `AssignCatalogRoleToPrincipalRole`
- [ ] Revoke Catalog Role - `RevokeCatalogRoleFromPrincipalRole`
- [ ] List Assignee Principals - `ListAssigneePrincipalsForPrincipalRole`

### Configuration (API)
- [ ] Get Config - `GetConfig`

## Catalog API

### Namespaces
Files: `cmd/catalog_namespaces.go`
- [x] List (`catalog namespaces list`) - `ListNamespaces`
- [x] Create (`catalog namespaces create`) - `CreateNamespace`
- [ ] Drop (`catalog namespaces drop`) - `DropNamespace`
- [ ] Get (`catalog namespaces get`) - `LoadNamespaceMetadata`
- [ ] Exists (`catalog namespaces exists`) - `NamespaceExists`
- [ ] Update (`catalog namespaces update`) - `UpdateProperties`

### Tables
Files: `cmd/catalog_tables.go`
- [x] List (`catalog tables list`) - `ListTables`
- [ ] Create (`catalog tables create`) - `CreateTable`
- [ ] Drop (`catalog tables drop`) - `DropTable`
- [ ] Get (`catalog tables get`) - `LoadTable`
- [ ] Exists (`catalog tables exists`) - `TableExists`
- [ ] Update (`catalog tables update`) - `UpdateTable`
- [ ] Register (`catalog tables register`) - `RegisterTable`
- [ ] Rename (`catalog tables rename`) - `RenameTable`
- [ ] Load Credentials - `LoadCredentials`
- [ ] Report Metrics - `ReportMetrics`
- [ ] Send Notification - `SendNotification`

### Views
Files: `cmd/catalog_views.go` (Proposed)
- [ ] List (`catalog views list`) - `ListViews`
- [ ] Create (`catalog views create`) - `CreateView`
- [ ] Drop (`catalog views drop`) - `DropView`
- [ ] Get (`catalog views get`) - `LoadView`
- [ ] Exists (`catalog views exists`) - `ViewExists`
- [ ] Update (`catalog views update`) - `ReplaceView`
- [ ] Rename (`catalog views rename`) - `RenameView`

### Other
- [ ] Commit Transaction - `CommitTransaction`
