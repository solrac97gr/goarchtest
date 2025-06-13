# Domain-Driven Design with Clean Architecture Example

This example demonstrates how to implement and test Domain-Driven Design (DDD) with Clean Architecture using GoArchTest.

## Architecture Overview

```
internal/
├── user/
│   ├── domain/          # User domain entities, value objects, interfaces
│   ├── application/     # User use cases, application services  
│   └── infrastructure/  # User repositories, external services
├── products/
│   ├── domain/          # Product domain entities, value objects, interfaces
│   ├── application/     # Product use cases, application services
│   └── infrastructure/  # Product repositories, external services
└── shared/              # Shared kernel (domain concepts used across domains)

pkg/
└── ...                  # Reusable utilities, libraries, frameworks
```

## Key Principles Enforced

### 1. **Bounded Context Isolation**
- User domain cannot import from products domain
- Products domain cannot import from user domain
- Each domain is completely independent

### 2. **Clean Architecture Within Each Domain**
- Domain layer has no dependencies on application or infrastructure
- Application layer only depends on domain
- Infrastructure layer depends on domain (implements interfaces)

### 3. **Shared Kernel Usage**
- Only domain layers can use the shared kernel
- Application and infrastructure cannot directly use shared

### 4. **Utility Layer (pkg/)**
- Any layer can use pkg/ utilities
- pkg/ should not depend on internal/ code

## Running the Tests

```bash
go test -v
```

This will validate:
- ✅ Cross-domain isolation
- ✅ Clean Architecture within each domain
- ✅ Proper shared kernel usage
- ✅ Architectural boundaries

## Real-World Benefits

- **Team Independence**: Teams can work on different domains without conflicts
- **Scalability**: Easy to split domains into microservices later
- **Maintainability**: Clear boundaries prevent architectural erosion
- **Testing**: Each domain can be tested in isolation
