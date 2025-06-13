# CQRS Architecture Testing Example

This example demonstrates how to test **CQRS (Command Query Responsibility Segregation)** and **Event Sourced CQRS** architectural patterns using GoArchTest.

## What is CQRS?

CQRS is an architectural pattern that separates the models for reading and writing data. Instead of using a single model for both commands (writes) and queries (reads), CQRS uses separate models optimized for their specific purposes.

### Key Principles:

1. **Separation of Concerns**: Commands and queries are handled by different models
2. **Independent Scaling**: Read and write sides can be scaled independently
3. **Optimized Models**: Write models focus on business logic, read models focus on queries
4. **No Cross-Dependencies**: Commands should not depend on queries and vice versa

## CQRS Architecture Components

```
┌─────────────┐    ┌─────────────┐
│  Commands   │    │   Queries   │
│ (Write Side)│    │ (Read Side) │
└─────────────┘    └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐    ┌─────────────┐
│ Write Model │    │ Read Model  │
│ (Aggregates)│    │(Projections)│
└─────────────┘    └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐    ┌─────────────┐
│ Write Store │    │ Read Store  │
│ (Transact.) │    │ (Optimized) │
└─────────────┘    └─────────────┘
```

## Event Sourced CQRS

Event Sourcing can be combined with CQRS where:

- Commands produce **Events** that are stored in an **Event Store**
- **Projections** build read models from events
- Queries use projections, not the event store directly

```
┌─────────────┐    ┌─────────────┐
│  Commands   │    │   Queries   │
└─────────────┘    └─────────────┘
       │                   │
       ▼                   │
┌─────────────┐            │
│   Events    │◄───────────┘
└─────────────┘    (via projections)
       │
       ▼
┌─────────────┐    ┌─────────────┐
│ Event Store │    │ Projections │
│             │───►│ (Read Models)│
└─────────────┘    └─────────────┘
```

## Running the Tests

```bash
# Navigate to the CQRS example directory
cd examples/cqrs

# Run the CQRS architecture tests
go test -v

# Run specific test functions
go test -v -run TestCQRSArchitecture
go test -v -run TestIndividualCQRSRules
go test -v -run TestEventSourcingRules
```

## Architecture Rules Tested

### Basic CQRS Rules:

1. **Commands ⚠️ Queries**: Commands should not depend on queries
2. **Queries ⚠️ Commands**: Queries should not depend on commands  
3. **Write Models ⚠️ Read Models**: Write models should not depend on read models
4. **Read Models ⚠️ Write Models**: Read models should not depend on write models
5. **Commands → Write Models**: Commands should primarily use write models
6. **Queries → Read Models**: Queries should primarily use read models

### Event Sourced CQRS Rules:

7. **Commands → Events**: Commands should depend on events (to produce them)
8. **Commands → Event Store**: Commands should interact with event store
9. **Queries ⚠️ Event Store**: Queries should not access event store directly
10. **Projections → Events**: Projections should depend on events
11. **Queries → Projections**: Queries should use projections for read models

## Expected Project Structure

For these tests to be meaningful, your project should follow this structure:

```
your-project/
├── commands/           # Command handlers
│   ├── create_user.go
│   ├── update_user.go
│   └── delete_user.go
├── queries/            # Query handlers  
│   ├── get_user.go
│   ├── list_users.go
│   └── search_users.go
├── writemodel/         # Write-optimized models
│   ├── user_aggregate.go
│   └── user_repository.go
├── readmodel/          # Read-optimized models
│   ├── user_view.go
│   └── user_projection.go
├── events/             # Domain events (for Event Sourcing)
│   ├── user_created.go
│   ├── user_updated.go
│   └── user_deleted.go
├── eventstore/         # Event store implementation
│   └── event_store.go
├── projections/        # Event projections
│   └── user_projection.go
└── domain/             # Shared domain models
    └── user.go
```

## Integration with Your Project

To integrate these tests into your own project:

1. **Adjust namespaces** to match your project structure
2. **Run tests** to identify violations
3. **Refactor code** to fix violations
4. **Add to CI/CD** pipeline for continuous validation

### Example Integration:

```go
func TestMyProjectCQRS(t *testing.T) {
    projectPath, _ := filepath.Abs("./")
    types := goarchtest.InPath(projectPath)
    
    // Adjust these namespaces to match your project
    cqrsPattern := goarchtest.CQRSArchitecture(
        "pkg/commands",     // Your command namespace
        "pkg/queries",      // Your query namespace  
        "pkg/domain",       // Your domain namespace
        "pkg/aggregates",   // Your write model namespace
        "pkg/views",        // Your read model namespace
    )
    
    validationResults := cqrsPattern.Validate(types)
    
    // Handle results...
}
```

## Common CQRS Violations

### ❌ Command-Query Cross-Dependencies
```go
// BAD: Command depending on query
package commands
import "myapp/queries" // ❌ Violation

// BAD: Query depending on command  
package queries
import "myapp/commands" // ❌ Violation
```

### ❌ Wrong Model Usage
```go
// BAD: Command using read model
package commands
import "myapp/readmodel" // ❌ Should use write model

// BAD: Query using write model
package queries  
import "myapp/writemodel" // ❌ Should use read model
```

### ❌ Direct Event Store Access
```go
// BAD: Query directly accessing event store
package queries
import "myapp/eventstore" // ❌ Should use projections
```

## Benefits of Testing CQRS Architecture

1. **Enforces Separation**: Prevents mixing of read/write concerns
2. **Maintains Performance**: Ensures optimized models for their purpose
3. **Enables Scaling**: Read and write sides can evolve independently
4. **Prevents Violations**: Catches architectural drift early
5. **Documents Intent**: Tests serve as living documentation

## Related Patterns

- **Clean Architecture**: Domain should be independent
- **Hexagonal Architecture**: Commands/queries are primary ports
- **Event Driven Architecture**: Events drive state changes
- **Domain-Driven Design**: Bounded contexts with CQRS

## Further Reading

- [Martin Fowler - CQRS](https://martinfowler.com/bliki/CQRS.html)
- [Event Sourcing Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/event-sourcing)
- [CQRS and Event Sourcing](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs)

## This Example's Test Results

This example includes **intentional architectural violations** to demonstrate how goarchtest detects CQRS violations.

### Expected Test Output

When you run `go test -v`, you should see:

```
❌ CQRS Rule #1 failed: Commands should not depend on queries
❌ CQRS Rule #2 failed: Queries should not depend on commands  
❌ CQRS Rule #3 failed: Write models should not depend on read models
❌ CQRS Rule #4 failed: Read models should not depend on write models
❌ CQRS Rule #5 failed: Commands should not depend on read models
❌ CQRS Rule #6 failed: Queries should not depend on write models
```

### Intentional Violations Included

#### ✅ Correct CQRS Implementation
- `commands/create_user_command.go` - Proper command using only write models
- `queries/get_user_query.go` - Proper query using only read models
- `writemodel/user_write_repository.go` - Write-optimized repository
- `readmodel/user_read_repository.go` - Read-optimized repository

#### ❌ Intentional Violations
- `commands/bad_create_user_command.go` - Violates CQRS by importing:
  - `queries` package (commands shouldn't depend on queries)
  - `readmodel` package (commands should use write models)
- `queries/bad_get_user_query.go` - Violates CQRS by importing:
  - `commands` package (queries shouldn't depend on commands)  
  - `writemodel` package (queries should use read models)

### Import Cycles

**Note**: The violations create import cycles, which is expected and demonstrates how architectural violations can lead to circular dependencies that prevent compilation.

### Key Takeaways

1. **Architecture Testing Works**: goarchtest successfully detects all CQRS violations
2. **Clear Feedback**: Each violation is reported with specific file and package information
3. **Preventive**: These tests help prevent architectural degradation over time
4. **Educational**: Shows developers what proper CQRS separation looks like vs. violations

This example proves that the CQRS architectural patterns in goarchtest are working correctly and can effectively enforce CQRS principles in real Go projects.
