# Hexagonal Architecture Example

This example demonstrates how to validate a **Hexagonal Architecture** (also known as **Ports and Adapters**) pattern using `goarchtest`.

## Architecture Overview

The Hexagonal Architecture pattern isolates the core business logic from external concerns by defining clear boundaries:

```
┌─────────────────────────────────────────────────────────────┐
│                    Primary Adapters                        │
│                   (Driving Adapters)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              HTTP Handlers                          │   │
│  │         (adapters/primary/http)                     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Core Domain                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 Ports                               │   │
│  │            (core/ports)                             │   │
│  │  ┌─────────────────────────────────────────────┐   │   │
│  │  │            Domain                           │   │   │
│  │  │        (core/domain)                        │   │   │
│  │  └─────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Secondary Adapters                        │
│                   (Driven Adapters)                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │     Database    │    Messaging    │   External APIs │   │
│  │  (adapters/secondary/database)   │   Services      │   │
│  │         (adapters/secondary/messaging)              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
custom-architecture/
├── core/
│   ├── domain/          # Business entities and rules
│   │   ├── order.go     # Order entity with business logic
│   │   └── errors.go    # Domain-specific errors
│   └── ports/           # Interfaces (contracts)
│       ├── repositories.go  # Repository interfaces
│       └── services.go      # Service interfaces and use cases
├── adapters/
│   ├── primary/         # Driving adapters (inbound)
│   │   └── http/        # HTTP handlers
│   │       └── order_handler.go
│   └── secondary/       # Driven adapters (outbound)
│       ├── database/    # Database implementations
│       │   └── order_repository.go
│       └── messaging/   # External service implementations
│           └── services.go
├── go.mod
├── main_test.go         # Architecture validation tests
└── README.md
```

## Architecture Rules

The test validates the following rules:

### 1. **Domain Layer Isolation**
- Domain should not depend on ports
- Domain should not depend on adapters
- Domain contains pure business logic

### 2. **Ports Layer Dependencies**
- Ports should depend on domain (for entities)
- Ports should not depend on adapters
- Ports define interfaces that adapters implement

### 3. **Primary Adapter Rules**
- Primary adapters should depend on ports
- Primary adapters should not depend on secondary adapters
- Primary adapters handle incoming requests

### 4. **Secondary Adapter Rules**
- Secondary adapters should depend on ports (implement interfaces)
- Secondary adapters should depend on domain (use entities)
- Secondary adapters should not depend on primary adapters
- Secondary adapters handle outgoing requests

## Business Domain

This example implements an **Order Management System** with:

- **Domain Entities**: `Order`, `OrderItem`, `OrderStatus`
- **Business Rules**: Order validation, total calculation
- **Use Cases**: Create order, get order, update status, cancel order
- **External Dependencies**: Database, payment processing, notifications

## Running the Tests

```bash
cd test/custom-architecture
go test -v
```

## Key Benefits of This Architecture

1. **Testability**: Core business logic is isolated and easily testable
2. **Flexibility**: Easy to swap implementations (database, messaging, etc.)
3. **Independence**: Business logic doesn't depend on frameworks or external tools
4. **Maintainability**: Clear separation of concerns
5. **Scalability**: Can evolve different parts independently

## Architecture Validation

The `main_test.go` file contains comprehensive tests that ensure:
- Dependency directions are correct
- Layers don't have forbidden dependencies
- Interfaces are properly placed
- Business logic remains pure

This validates that the implementation follows the Hexagonal Architecture principles correctly.
