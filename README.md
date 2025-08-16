# BC CodeGen

A CLI tool that parses Microsoft Business Central OData API metadata and generates typed schemas for various languages.

## Features

- **API-focused filtering**: Only generates types for entities exposed through the OData API (EntitySets)
- **Navigation properties**: Automatically includes related entities as optional properties using `z.lazy()` for circular references
- **Create/Update types**: Generates separate TypeScript types for POST (Create) and PATCH (Update) operations
- **Smart property filtering**: Excludes read-only fields (id, timestamps, etc.) from Create/Update types
- **Branded types**: Uses Zod branded strings for GUIDs, DateTime, and Date for type safety
- **Complex type handling**: Option to skip complex Business Central-specific types
- **Dependency resolution**: Automatically includes referenced entities and their dependencies
- **Extensible**: Framework ready for additional target languages

## Installation

```bash
go build -o bc-codegen
```

## Usage

```bash
# Generate TypeScript schemas from metadata
./bc-codegen -input metadata.xml -output ./generated

# Generate with specific options
./bc-codegen -input metadata.xml -output ./types -skip-complex=false
```

### CLI Options

- `-input` (required): Path to OData metadata XML file
- `-output`: Output directory for generated files (default: current directory)
- `-lang`: Target language, currently supports "typescript" (default: "typescript")
- `-package`: Package/module name for generated code
- `-skip-complex`: Skip complex types during generation (default: true)

## Generated Output

For TypeScript, the tool generates:

```typescript
import { z } from "zod";

// Branded types for Business Central
const Guid = z.string().brand<"Guid">();
const DateTime = z.string().brand<"DateTime">();
const Date = z.string().brand<"Date">();

export const SalesOrder = z.object({
  id: Guid,
  number: z.string(),
  orderDate: Date,
  customerId: Guid,
  customerName: z.string(),
  totalAmountIncludingTax: z.number(),
  // Navigation properties
  customer: z.lazy(() => Customer).optional(),
  salesOrderLines: z.array(z.lazy(() => SalesOrderLine)).optional(),
  // ... more properties
});

export type SalesOrder = z.infer<typeof SalesOrder>;

// Create type (for POST operations)
export type SalesOrderCreate = {
  orderDate: string;
  customerId: string;
  customerName: string;
  totalAmountIncludingTax: number;
  // ... excludes read-only fields like id, number, lastModifiedDateTime
};

// Update type (for PATCH operations) 
export type SalesOrderUpdate = {
  orderDate?: string;
  customerId?: string;
  customerName?: string;
  totalAmountIncludingTax?: number;
  // ... all fields optional, excludes id and other key fields
};
```

## Type Mappings

| OData Type | TypeScript/Zod |
|------------|-----------------|
| Edm.Guid | Branded string |
| Edm.DateTimeOffset | Branded DateTime string |
| Edm.Date | Branded Date string |
| Edm.String | z.string() |
| Edm.Int32/Int64 | z.number().int() |
| Edm.Decimal/Double | z.number() |
| Edm.Boolean | z.boolean() |
| Collection(T) | z.array(T) |

## Example

With the provided v2-metadata.xml:

```bash
go run . -input v2-metadata.xml -output generated
```

This generates `generated/schemas.ts` with Zod schemas and TypeScript types for all Business Central entities like SalesOrder, Item, Customer, etc.