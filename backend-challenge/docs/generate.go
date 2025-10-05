package docs

// Run: go generate ./docs
// Or from docs folder: go generate

//go:generate go tool oapi-codegen -config config-types.yaml openapi.yaml
//go:generate go tool oapi-codegen -config config-server.yaml openapi.yaml
