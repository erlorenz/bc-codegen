package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/erlorenz/bc-codegen/metadata"
)

func NewTypeScript(outDir, name string) (*TypeScriptGenerator, error) {
	outputFile := filepath.Join(outDir, name+".ts")

	generator := &TypeScriptGenerator{
		outputPath: outputFile,
		builder:    &strings.Builder{},
	}

	return generator, nil
}

type TypeScriptGenerator struct {
	outputPath string
	builder    *strings.Builder
}

func (g *TypeScriptGenerator) Generate(model metadata.Model) error {
	g.writeHeader()

	// Build a map of API-accessible entities from EntitySets (excluding specific entities)
	apiEntities := make(map[string]bool)
	excludedEntities := map[string]bool{
		"company":           true,
		"entityMetadata":    true,
		"apicategoryroutes": true,
	}

	for _, entitySet := range model.DataServices.Schema.EntityContainer.EntitySets {
		entityTypeName := g.extractEntityTypeName(entitySet.EntityType)
		if !excludedEntities[entityTypeName] {
			apiEntities[entityTypeName] = true
		}
	}

	// Collect all referenced types through navigation properties
	allReferencedTypes := make(map[string]bool)
	for entityName := range apiEntities {
		g.collectReferencedTypes(entityName, &model, allReferencedTypes)
	}

	// Build dependency graph and sort topologically
	entitiesToGenerate := make([]metadata.EntityType, 0)
	for _, entityType := range model.DataServices.Schema.EntityTypes {
		if (apiEntities[entityType.Name] || allReferencedTypes[entityType.Name]) && !excludedEntities[entityType.Name] {
			entitiesToGenerate = append(entitiesToGenerate, entityType)
		}
	}

	// Generate main entity schemas
	for _, entityType := range entitiesToGenerate {
		g.writeEntityType(entityType)
	}

	// Generate Create and Update types
	for _, entityType := range entitiesToGenerate {
		g.writeCreateType(entityType)
		g.writeUpdateType(entityType)
	}

	// Write the content to file
	return g.writeToFile()
}

func (g *TypeScriptGenerator) writeToFile() error {
	file, err := os.Create(g.outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.WriteString(g.builder.String())
	return err
}

func (g *TypeScriptGenerator) collectReferencedTypes(entityName string, model *metadata.Model, referencedTypes map[string]bool) {
	if referencedTypes[entityName] {
		return // Already processed
	}
	referencedTypes[entityName] = true

	// Find the entity and collect its navigation property references
	for _, entityType := range model.DataServices.Schema.EntityTypes {
		if entityType.Name == entityName {
			for _, navProp := range entityType.NavigationProperties {
				targetType := g.extractNavigationTargetType(navProp.Type)
				excludedEntities := map[string]bool{
					"company":           true,
					"entityMetadata":    true,
					"apicategoryroutes": true,
				}
				if targetType != entityName && !excludedEntities[targetType] { // Avoid self-reference and excluded entities
					g.collectReferencedTypes(targetType, model, referencedTypes)
				}
			}
			break
		}
	}
}

func (g *TypeScriptGenerator) extractEntityTypeName(fullTypeName string) string {
	parts := strings.Split(fullTypeName, ".")
	return parts[len(parts)-1]
}

func (g *TypeScriptGenerator) extractNavigationTargetType(navType string) string {
	// Handle Collection(Microsoft.NAV.entityName) or Microsoft.NAV.entityName
	if strings.HasPrefix(navType, "Collection(") {
		navType = strings.TrimSuffix(strings.TrimPrefix(navType, "Collection("), ")")
	}
	return g.extractEntityTypeName(navType)
}

func (g *TypeScriptGenerator) writeHeader() {
	g.builder.WriteString(`import { z } from "zod";` + "\n")
	g.builder.WriteString("\n")

	g.builder.WriteString(`// Branded types for Business Central` + "\n")
	g.builder.WriteString(`const Guid = z.string().brand<"Guid">();` + "\n")
	g.builder.WriteString(`const DateTime = z.string().brand<"DateTime">();` + "\n")
	g.builder.WriteString(`const DateOnly = z.string().brand<"DateOnly">();` + "\n")
	g.builder.WriteString("\n")

	g.builder.WriteString(`// Generic reference types` + "\n")
	g.builder.WriteString(`const RefOne = z.object({ id: Guid });` + "\n")
	g.builder.WriteString(`const RefMany = z.array(RefOne);` + "\n")
	g.builder.WriteString("\n")
	g.builder.WriteString(`export type RefOne = z.infer<typeof RefOne>;` + "\n")
	g.builder.WriteString(`export type RefMany = z.infer<typeof RefMany>;` + "\n")
	g.builder.WriteString("\n")

}

func (g *TypeScriptGenerator) writeEntityType(entity metadata.EntityType) {
	schemaName := toPascalCase(entity.Name)

	g.builder.WriteString(fmt.Sprintf("export const %s = z.object({\n", schemaName))

	for _, prop := range entity.Properties {
		if g.isComplexProperty(prop) {
			continue
		}
		g.writeProperty(prop)
	}

	// Add navigation properties
	for _, navProp := range entity.NavigationProperties {
		g.writeNavigationProperty(navProp)
	}

	g.builder.WriteString("});\n")
	g.builder.WriteString("\n")
	g.builder.WriteString(fmt.Sprintf("export type %s = z.infer<typeof %s>;\n", schemaName, schemaName))
	g.builder.WriteString("\n")
}

func (g *TypeScriptGenerator) writeCreateType(entity metadata.EntityType) {
	typeName := toPascalCase(entity.Name) + "Create"

	g.builder.WriteString(fmt.Sprintf("export type %s = {\n", typeName))

	for _, prop := range entity.Properties {
		if g.isComplexProperty(prop) || g.isReadOnlyProperty(prop) {
			continue
		}
		g.writeOptionalTypeProperty(prop)
	}

	g.builder.WriteString("};\n")
	g.builder.WriteString("\n")
}

func (g *TypeScriptGenerator) writeUpdateType(entity metadata.EntityType) {
	typeName := toPascalCase(entity.Name) + "Update"

	g.builder.WriteString(fmt.Sprintf("export type %s = {\n", typeName))

	for _, prop := range entity.Properties {
		if g.isComplexProperty(prop) || g.isReadOnlyForUpdateProperty(prop) || g.isKeyProperty(entity, prop) {
			continue
		}
		g.writeOptionalTypeProperty(prop)
	}

	g.builder.WriteString("};\n")
	g.builder.WriteString("\n")
}

func (g *TypeScriptGenerator) writeNavigationProperty(navProp metadata.NavigationProperty) {
	propName := toCamelCase(navProp.Name)

	// Use generic reference types - just id for now, can be expanded later
	if strings.HasPrefix(navProp.Type, "Collection(") {
		g.builder.WriteString(fmt.Sprintf("  %s: RefMany.optional(),\n", propName))
	} else {
		g.builder.WriteString(fmt.Sprintf("  %s: RefOne.optional(),\n", propName))
	}
}

func (g *TypeScriptGenerator) writeOptionalTypeProperty(prop metadata.Property) {
	propType := g.mapODataTypeToTypeScript(prop.Type)
	g.builder.WriteString(fmt.Sprintf("  %s?: %s;\n", toCamelCase(prop.Name), propType))
}

func (g *TypeScriptGenerator) isReadOnlyProperty(prop metadata.Property) bool {
	// Common read-only properties in Business Central (excluding ID which can be provided for creates)
	readOnlyProps := map[string]bool{
		"systemVersion":        true,
		"timestamp":            true,
		"systemCreatedAt":      true,
		"systemCreatedBy":      true,
		"systemModifiedAt":     true,
		"systemModifiedBy":     true,
		"lastModifiedDateTime": true,
		"entryNumber":          true,
		"number":               true, // Often auto-generated
	}
	return readOnlyProps[prop.Name]
}

func (g *TypeScriptGenerator) isReadOnlyForUpdateProperty(prop metadata.Property) bool {
	// Properties that are read-only for updates (including ID)
	readOnlyProps := map[string]bool{
		"id":                   true, // ID cannot be changed in updates
		"systemVersion":        true,
		"timestamp":            true,
		"systemCreatedAt":      true,
		"systemCreatedBy":      true,
		"systemModifiedAt":     true,
		"systemModifiedBy":     true,
		"lastModifiedDateTime": true,
		"entryNumber":          true,
		"number":               true, // Often auto-generated
	}
	return readOnlyProps[prop.Name]
}

func (g *TypeScriptGenerator) isKeyProperty(entity metadata.EntityType, prop metadata.Property) bool {
	for _, keyProp := range entity.Key.PropertyRefs {
		if keyProp.Name == prop.Name {
			return true
		}
	}
	return false
}

func (g *TypeScriptGenerator) writeProperty(prop metadata.Property) {
	zodType := g.mapODataTypeToZod(prop.Type)
	isOptional := prop.Nullable == "true"

	if isOptional {
		g.builder.WriteString(fmt.Sprintf("    %s: %s.optional(),\n", toCamelCase(prop.Name), zodType))
	} else {
		g.builder.WriteString(fmt.Sprintf("    %s: %s,\n", toCamelCase(prop.Name), zodType))
	}
}

func (g *TypeScriptGenerator) mapODataTypeToZod(odataType string) string {
	switch {
	case strings.Contains(odataType, "Guid"):
		return "Guid"
	case strings.Contains(odataType, "DateTime"):
		return "DateTime"
	case strings.Contains(odataType, "Date"):
		return "DateOnly"
	case odataType == "Edm.String":
		return "z.string()"
	case odataType == "Edm.Int32" || odataType == "Edm.Int64":
		return "z.number().int()"
	case odataType == "Edm.Decimal" || odataType == "Edm.Double":
		return "z.number()"
	case odataType == "Edm.Boolean":
		return "z.boolean()"
	case strings.HasPrefix(odataType, "Collection("):
		innerType := strings.TrimSuffix(strings.TrimPrefix(odataType, "Collection("), ")")
		return fmt.Sprintf("z.array(%s)", g.mapODataTypeToZod(innerType))
	default:
		return "z.unknown()"
	}
}

func (g *TypeScriptGenerator) mapODataTypeToTypeScript(odataType string) string {
	switch {
	case strings.Contains(odataType, "Guid"):
		return "string"
	case strings.Contains(odataType, "DateTime"):
		return "string"
	case strings.Contains(odataType, "Date"):
		return "string"
	case odataType == "Edm.String":
		return "string"
	case odataType == "Edm.Int32" || odataType == "Edm.Int64":
		return "number"
	case odataType == "Edm.Decimal" || odataType == "Edm.Double":
		return "number"
	case odataType == "Edm.Boolean":
		return "boolean"
	case strings.HasPrefix(odataType, "Collection("):
		innerType := strings.TrimSuffix(strings.TrimPrefix(odataType, "Collection("), ")")
		return fmt.Sprintf("%s[]", g.mapODataTypeToTypeScript(innerType))
	default:
		return "unknown"
	}
}

func (g *TypeScriptGenerator) isComplexProperty(prop metadata.Property) bool {
	return strings.Contains(prop.Type, "Microsoft.NAV") ||
		strings.HasPrefix(prop.Type, "Collection(Microsoft.NAV")
}

func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
