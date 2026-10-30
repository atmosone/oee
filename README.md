## OpenAPI

- Code generation for machines server
```bash
oapi-codegen -generate types,std-http,spec -package openapi -o internal/modules/machines/openapi/openapi.go api/openapi/machines.openapi.json
```

- Code generation for workshifts server
```bash
oapi-codegen -generate types,std-http,spec -package openapi -o internal/modules/workshifts/openapi/openapi.go api/openapi/workshifts.openapi.json
```

