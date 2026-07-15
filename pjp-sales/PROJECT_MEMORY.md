# Sales Scylla X - Project Memory

## Project Overview

**Sales Scylla X** is a comprehensive sales management system built with Go (Golang) and Fiber framework. It's a microservice that handles sales operations including orders, invoices, promotions, discounts, and various business processes for a distribution/sales company.

## Architecture

### Technology Stack
- **Language**: Go 1.23.0
- **Framework**: Fiber v2 (Fast HTTP framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Cloud Storage**: Huawei Cloud OBS (Object Storage Service)
- **Message Queue**: RabbitMQ
- **Containerization**: Docker
- **Validation**: Go Playground Validator

### Project Structure

```
sales/
├── adapter/           # External service adapters (OBS Huawei)
├── controller/        # HTTP controllers (API endpoints)
├── entity/           # Request/Response DTOs and business entities
├── model/            # Database models (GORM)
├── repository/       # Data access layer
├── service/          # Business logic layer
├── pkg/              # Shared packages and utilities
│   ├── config/       # Configuration management
│   ├── middleware/   # HTTP middleware
│   ├── validation/   # Input validation
│   └── ...
├── migration/        # Database migrations
└── main.go          # Application entry point
```

## Core Business Domains

### 1. Order Management
- **Orders (RO - Receipt Orders)**: Core sales orders with status tracking
- **Order Details**: Line items with products, quantities, pricing
- **Order Rewards**: Promotional rewards and bonuses
- **Order Approval**: Multi-level approval workflows

### 2. Invoice Management
- **Invoices**: Billing and invoicing system
- **Invoice Details**: Line items for invoicing
- **Payment Types**: Cash on delivery, credit, etc.

### 3. Promotion System (V2)
- **Promotions**: Advanced promotional campaigns
- **Promotion Types**: 
  - **Slab**: Tiered discount structures
  - **Strata**: Sequential reward structures
- **Promotion Criteria**: Product and outlet targeting
- **Reward Products**: Free products as rewards

### 4. Discount Management
- **Discount Criteria**: Rule-based discount calculations
- **Discount Groups**: Categorized discount structures
- **Minimum Price Management**: Price floor controls

### 5. Inventory & Stock
- **Stock Management**: Warehouse inventory tracking
- **Product Conversion**: Unit conversion (smallest, middle, largest units)
- **Stock Validation**: Availability checks

### 6. Sales Operations
- **Sales Orders (SO)**: Sales order processing
- **Returns**: Return order management
- **Consignment**: Consignment inventory
- **TLS (Transfer)**: Inter-warehouse transfers

### 7. Reporting & Analytics
- **Reports**: Business intelligence and reporting
- **Gamification**: Sales performance tracking
- **Hierarchy Approval**: Multi-level approval system

## API Structure

### Authentication
- JWT-based authentication with middleware protection
- All endpoints require valid JWT tokens

### API Versioning
- **v1**: Legacy API endpoints
- **v2**: Enhanced API endpoints (especially for promotions)

### Key Endpoints

#### Orders (`/v1/orders`)
- `POST /` - Create new order
- `GET /` - List orders with filtering
- `GET /:ro_no` - Get order details
- `PATCH /:ro_no` - Update order
- `DELETE /:ro_no` - Delete order
- `POST /conversion` - Product unit conversion
- `GET /discount` - Get discount information

#### Promotions (`/v1/promotions`, `/v2/promotions`)
- `POST /` - Create promotion
- `GET /` - List promotions
- `GET /:promo_id` - Get promotion details
- `PATCH /:promo_id` - Update promotion
- `DELETE /:promo_id` - Delete promotion
- `POST /consult` - Consult promotion eligibility

#### Invoices (`/v1/invoices`)
- `POST /` - Create invoice
- `GET /` - List invoices
- `GET /:ro_no` - Get invoice details
- `PATCH /:ro_no` - Update invoice

## Data Models

### Core Entities

#### Order Entity
```go
type Order struct {
    CustID            string     // Customer ID
    RoNo              string     // Receipt Order Number
    SalesmanId        *int64     // Sales representative
    WhId              *int64     // Warehouse ID
    OutletID          *int64     // Outlet ID
    SubTotal          *float64   // Subtotal amount
    Total             *float64   // Total amount
    DataStatus        *int64     // Order status
    // ... many more fields
}
```

#### Promotion V2 Entity
```go
type PromotionV2 struct {
    PromoID           string           // Promotion ID
    PromoDesc         string           // Description
    PromoType         PromotionType    // slab | strata
    PromoStatus       PromotionV2Status // draft | active | closed
    EffectiveFrom     string           // Start date
    EffectiveTo       string           // End date
    // ... detailed promotion configuration
}
```

### Status Enums

#### Order Status
- `NEED_REVIEW` (1) - Pending review
- `PROCESSED` (2) - Processed
- `ON_DELIVERY` (3) - In transit
- `RECEIVED` (4) - Delivered
- `COMPLETED` (7) - Completed
- `CANCELLED` (9) - Cancelled

#### Payment Types
- `PAY_TYPE_CASH_ON_DELIVERY` (1)
- `PAY_TYPE_CASH_BEFORE_DELIVERY` (2)
- `PAY_TYPE_CREDIT` (3)

## Database Schema

### Key Tables
- `sls.order` - Main orders table
- `sls.order_detail` - Order line items
- `sls.invoice` - Invoice records
- `promo.promotions` - Promotion campaigns (V2)
- `sls.stock` - Inventory management
- `sls.return` - Return orders

### Database Features
- **Schema Separation**: Different schemas for different domains (`sls`, `promo`, `mst`)
- **Soft Deletes**: Using GORM's soft delete functionality
- **Audit Trails**: Created/Updated timestamps and user tracking
- **Enums**: PostgreSQL enums for status and type fields

## Business Logic

### Order Processing Flow
1. **Order Creation**: Validate customer, products, stock
2. **Promotion Consultation**: Check eligible promotions
3. **Discount Calculation**: Apply applicable discounts
4. **Stock Validation**: Ensure product availability
5. **Credit Limit Check**: Validate customer credit
6. **Approval Workflow**: Multi-level approval if required
7. **Invoice Generation**: Create invoice after approval

### Promotion System (V2)
- **Slab Promotions**: Tiered discounts based on quantity/value thresholds
- **Strata Promotions**: Sequential rewards with multiple tiers
- **Budget Control**: Limited or unlimited budget tracking
- **Coverage Control**: National or distributor-specific coverage
- **Claim Management**: Full or partial claim processing

### Validation System
- **Stock Validation**: Real-time inventory checks
- **Credit Limit Validation**: Customer credit limit enforcement
- **Overdue Validation**: Payment overdue checks
- **Outstanding Validation**: Outstanding balance checks

## Configuration

### Environment Variables
- Database connection settings
- JWT secret keys
- OBS (Object Storage) credentials
- Server configuration (timeouts, ports)
- Application metadata (name, version, status)

### Database Configuration
- Connection pooling with configurable limits
- Debug mode for development
- Connection timeouts and lifecycle management

## External Integrations

### Huawei Cloud OBS
- File storage for documents and reports
- Image and document management
- Cloud-based file operations

### RabbitMQ
- Message queuing for asynchronous operations
- Event-driven architecture support

## Development Features

### Code Organization
- **Clean Architecture**: Clear separation of concerns
- **Dependency Injection**: Service and repository injection
- **Interface-based Design**: Testable and maintainable code
- **Validation**: Comprehensive input validation
- **Error Handling**: Structured error responses

### Middleware Stack
- **CORS**: Cross-origin resource sharing
- **Logging**: Request/response logging
- **Recovery**: Panic recovery
- **Request ID**: Unique request tracking
- **JWT Protection**: Authentication middleware

## Deployment

### Docker Support
- Multi-stage Docker build
- Alpine Linux base image
- Optimized for production deployment
- Port 9004 exposed

### Database Migrations
- SQL migration files in `migration/` directory
- Schema creation and updates
- Enum type definitions

## Key Features

### Sales Management
- Complete order-to-invoice lifecycle
- Multi-level approval workflows
- Real-time stock validation
- Credit limit management

### Promotion Engine
- Advanced promotion system with V2 architecture
- Flexible discount structures (slab/strata)
- Budget and coverage controls
- Claim management system

### Reporting & Analytics
- Comprehensive reporting system
- Sales performance tracking
- Business intelligence features

### Integration Capabilities
- Cloud storage integration
- Message queue support
- External service adapters

## Current Development Status

The project is currently on the `promotion-v2` branch with active development on:
- Promotion V2 system enhancements
- Database schema updates
- Service layer improvements
- Repository pattern implementations

This is a mature, production-ready sales management system with comprehensive business logic, robust architecture, and extensive feature set for managing complex sales operations.
