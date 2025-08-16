package generate

import "github.com/erlorenz/bc-codegen/metadata"

// Generator generates schemas and types from a metadata.EdmxModel.
type Generator interface {
	Generate(metadata.Model) error
}
