import assert from "node:assert/strict";
import { describe, test } from "node:test";
import z from "zod";
import * as V2 from "./generated/v2.schema.ts";

describe("SalesOrder Schema Tests", () => {
	const fakeOrder = {
		id: "550e8400-e29b-41d4-a716-446655440000",
		number: "SO-001",
		externalDocumentNumber: "EXT-12345",
		orderDate: "2024-01-15",
		postingDate: "2024-01-15",
		customerId: "550e8400-e29b-41d4-a716-446655440001",
		customerNumber: "CUST-001",
		customerName: "Test Customer Inc.",
		billToName: "Test Customer Inc.",
		billToCustomerId: "550e8400-e29b-41d4-a716-446655440001",
		billToCustomerNumber: "CUST-001",
		shipToName: "Test Customer Inc.",
		shipToContact: "John Doe",
		sellToAddressLine1: "123 Main St",
		sellToAddressLine2: "Suite 100",
		sellToCity: "Seattle",
		sellToCountry: "US",
		sellToState: "WA",
		sellToPostCode: "98101",
		billToAddressLine1: "123 Main St",
		billToAddressLine2: "Suite 100",
		billToCity: "Seattle",
		billToCountry: "US",
		billToState: "WA",
		billToPostCode: "98101",
		shipToAddressLine1: "123 Main St",
		shipToAddressLine2: "Suite 100",
		shipToCity: "Seattle",
		shipToCountry: "US",
		shipToState: "WA",
		shipToPostCode: "98101",
		shortcutDimension1Code: "",
		shortcutDimension2Code: "",
		currencyId: "550e8400-e29b-41d4-a716-446655440002",
		currencyCode: "USD",
		pricesIncludeTax: false,
		paymentTermsId: "550e8400-e29b-41d4-a716-446655440003",
		shipmentMethodId: "550e8400-e29b-41d4-a716-446655440004",
		salesperson: "SP001",
		partialShipping: false,
		requestedDeliveryDate: "2024-01-20",
		discountAmount: 0,
		discountAppliedBeforeTax: false,
		totalAmountExcludingTax: 149.95,
		totalTaxAmount: 12.00,
		totalAmountIncludingTax: 161.95,
		fullyShipped: false,
		lastModifiedDateTime: "2024-01-15T10:30:00Z",
		phoneNumber: "555-123-4567",
		email: "orders@testcustomer.com",
		customer: { id: "550e8400-e29b-41d4-a716-446655440001" },
		salesOrderLines: [
			{ id: "550e8400-e29b-41d4-a716-446655440010" },
			{ id: "550e8400-e29b-41d4-a716-446655440011" }
		],
	};

	const fakeLines = [
		{
			id: "550e8400-e29b-41d4-a716-446655440010",
			documentId: "550e8400-e29b-41d4-a716-446655440000",
			sequence: 1,
			itemId: "550e8400-e29b-41d4-a716-446655440020",
			accountId: "550e8400-e29b-41d4-a716-446655440021",
			lineObjectNumber: "ITEM-001",
			description: "Test Product 1",
			description2: "",
			unitOfMeasureId: "550e8400-e29b-41d4-a716-446655440022",
			unitOfMeasureCode: "PCS",
			quantity: 5,
			unitPrice: 29.99,
			discountAmount: 0,
			discountPercent: 0,
			discountAppliedBeforeTax: false,
			amountExcludingTax: 149.95,
			taxCode: "STANDARD",
			taxPercent: 8.0,
			totalTaxAmount: 12.00,
			amountIncludingTax: 161.95,
			invoiceDiscountAllocation: 0,
			netAmount: 149.95,
			netTaxAmount: 12.00,
			netAmountIncludingTax: 161.95,
			shipmentDate: "2024-01-20",
			shippedQuantity: 0,
			invoicedQuantity: 0,
			invoiceQuantity: 5,
			shipQuantity: 5,
			itemVariantId: "550e8400-e29b-41d4-a716-446655440023",
			locationId: "550e8400-e29b-41d4-a716-446655440024",
			salesOrder: { id: "550e8400-e29b-41d4-a716-446655440000" },
			item: { id: "550e8400-e29b-41d4-a716-446655440020" },
		},
		{
			id: "550e8400-e29b-41d4-a716-446655440011",
			documentId: "550e8400-e29b-41d4-a716-446655440000",
			sequence: 2,
			itemId: "550e8400-e29b-41d4-a716-446655440025",
			accountId: "550e8400-e29b-41d4-a716-446655440021",
			lineObjectNumber: "ITEM-002",
			description: "Test Product 2",
			description2: "",
			unitOfMeasureId: "550e8400-e29b-41d4-a716-446655440022",
			unitOfMeasureCode: "PCS",
			quantity: 3,
			unitPrice: 49.99,
			discountAmount: 7.50,
			discountPercent: 5,
			discountAppliedBeforeTax: true,
			amountExcludingTax: 142.47,
			taxCode: "STANDARD",
			taxPercent: 8.0,
			totalTaxAmount: 11.40,
			amountIncludingTax: 153.87,
			invoiceDiscountAllocation: 0,
			netAmount: 142.47,
			netTaxAmount: 11.40,
			netAmountIncludingTax: 153.87,
			shipmentDate: "2024-01-20",
			shippedQuantity: 0,
			invoicedQuantity: 0,
			invoiceQuantity: 3,
			shipQuantity: 3,
			itemVariantId: "550e8400-e29b-41d4-a716-446655440026",
			locationId: "550e8400-e29b-41d4-a716-446655440024",
			salesOrder: { id: "550e8400-e29b-41d4-a716-446655440000" },
			item: { id: "550e8400-e29b-41d4-a716-446655440025" },
		},
	];

	test("SalesOrder parses correctly", () => {
		const result = V2.SalesOrder.parse(fakeOrder);

		assert.equal(result.id, "550e8400-e29b-41d4-a716-446655440000");
		assert.equal(result.number, "SO-001");
		assert.equal(result.customerName, "Test Customer Inc.");
		assert.equal(result.currencyCode, "USD");
		assert.equal(result.totalAmountIncludingTax, 161.95);
		assert.deepEqual(result.customer, { id: "550e8400-e29b-41d4-a716-446655440001" });
	});

	test("SalesOrderWithLines parses correctly", () => {
		const SalesOrderWithLines = z.object({
			...V2.SalesOrder.shape,
			salesOrderLinesDetailed: z.array(V2.SalesOrderLine).optional(),
		});

		const salesOrderData = {
			...fakeOrder,
			salesOrderLinesDetailed: fakeLines,
		};

		const result = SalesOrderWithLines.parse(salesOrderData);

		assert.equal(result.id, "550e8400-e29b-41d4-a716-446655440000");
		assert.equal(result.customerName, "Test Customer Inc.");
		assert.deepEqual(result.customer, { id: "550e8400-e29b-41d4-a716-446655440001" });
		assert.equal(result.salesOrderLinesDetailed?.length, 2);
		assert(result.salesOrderLinesDetailed, "salesOrderLinesDetailed should exist");
		assert.equal(result.salesOrderLinesDetailed[0]?.id, "550e8400-e29b-41d4-a716-446655440010");
		assert.equal(result.salesOrderLinesDetailed[0]?.description, "Test Product 1");
		assert.equal(result.salesOrderLinesDetailed[1]?.id, "550e8400-e29b-41d4-a716-446655440011");
		assert.equal(result.salesOrderLinesDetailed[1]?.description, "Test Product 2");
	});


});
