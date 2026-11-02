## OpenAPI

- Code generation for CMS server

```bash
oapi-codegen -generate types,std-http -package schema -o internal/apps/cms/adapters/openapi/schema/schema.go ../api/openapi/cms.openapi.json
```