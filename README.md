# BC CodeGen

A CLI tool that parses Microsoft Business Central OData API metadata and generates typed schemas for TypeScript.

## Features

- **TypeScript/Zod Schema Generation**: Generates Zod schemas and TypeScript types from Business Central OData metadata
- **API-focused filtering**: Only generates types for entities exposed through the OData API (EntitySets)
- **Navigation properties**: Automatically includes related entities as optional properties for circular references
- **Create/Update types**: Generates separate TypeScript types for POST (Create) and PATCH (Update) operations
- **Smart property filtering**: Excludes read-only fields (id, timestamps, etc.) from Create/Update types
- **Branded types**: Uses Zod branded strings for GUIDs, DateTime, and Date for type safety
- **Bound Actions**: Generates union types for available bound actions (e.g., `SalesOrderActions = 'post' | 'ship'`)
- **Complex type handling**: Skips complex Business Central-specific types by default
- **Dependency resolution**: Automatically includes referenced entities and their dependencies
- **Extensible**: Framework ready for additional target languages in the future

## Installation

```bash
go build -o bc-codegen
```

## Usage

```bash
# Generate TypeScript schemas from metadata
./bc-codegen -out generated/schema.ts metadata.xml

# Generate with short flag
./bc-codegen -o types/v2.schema.ts metadata.xml
```

### CLI Options

- `-out`: Output file path (default: "schema.ts")
- `-o`: Output file path (short form)
- `-lang`: Target language, currently supports "typescript" (default: "typescript")

## Generated Output

For TypeScript, the tool generates:

```typescript
import { z } from "zod";

// Branded types for Business Central
const Guid = z.string().brand<"Guid">();
const DateTime = z.string().brand<"DateTime">();
const DateOnly = z.string().brand<"DateOnly">();

export const SalesOrder = z.object({
  id: Guid,
  number: z.string(),
  orderDate: DateOnly,
  customerId: Guid,
  customerName: z.string(),
  totalAmountIncludingTax: z.number(),
  // Navigation properties
  customer: RefOne.optional(),
  salesOrderLines: RefMany.optional(),
  // ... more properties
});

export type SalesOrder = z.infer<typeof SalesOrder>;

// Create type (for POST operations)
export type SalesOrderCreate = {
  orderDate?: string;
  customerId?: string;
  customerName?: string;
  totalAmountIncludingTax?: number;
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

// Bound Action Types
export type SalesOrderActions = 'shipAndInvoice';
export type SalesInvoiceActions = 'cancel' | 'cancelAndSend' | 'makeCorrectiveCreditMemo' | 'post' | 'postAndSend' | 'send';
```

## Type Mappings

| OData Type | Output Type |
|------------|-------------|
| Edm.Guid | Branded GUID string |
| Edm.DateTimeOffset | Branded DateTime string |
| Edm.Date | Branded DateOnly string |
| Edm.String | string |
| Edm.Int32/Int64 | number |
| Edm.Decimal/Double | number |
| Edm.Boolean | boolean |
| Collection(T) | Array<T> |

## Example

With Business Central OData metadata:

```bash
./bc-codegen -out generated/v2.schema.ts metadata.xml
```

This generates `generated/v2.schema.ts` with:
- Zod schemas and TypeScript types for all Business Central entities (SalesOrder, Item, Customer, etc.)
- Create and Update types for API operations
- Union types for bound actions available on each entity