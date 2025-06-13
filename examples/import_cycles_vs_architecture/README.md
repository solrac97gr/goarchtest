# Import Cycles vs Architectural Violations

This example demonstrates the key difference between what Go's compiler prevents (import cycles) and what GoArchTest prevents (architectural violations).

## The Key Difference

### Go Compiler: Prevents Import Cycles
```go
// Go compiler ERROR: import cycle not allowed
package A
import "B"  // A → B

package B  
import "A"  // B → A (creates cycle A → B → A)
```

### GoArchTest: Prevents Architectural Violations
```go
// Go compiler: ✅ Compiles fine (no cycle)
// Clean Architecture: ❌ Violation (inner depends on outer)
package domain          // Inner layer
import "infrastructure" // Depends on outer layer

package infrastructure  // Outer layer  
import "domain"         // Depends on inner layer (this is fine)
```

## Running the Example

```bash
go test -v
```

This will run tests that demonstrate:

1. **What Go allows but violates architecture** - Examples of code that compiles but breaks architectural principles
2. **How GoArchTest catches violations** - Specific tests that detect architectural problems
3. **Value proposition** - Why you need both Go's import cycle protection AND GoArchTest's architectural enforcement

## Real-World Scenarios

### Scenario 1: Domain Layer Pollution
```go
// Compiles in Go ✅, but violates Clean Architecture ❌
package domain
import "myapp/infrastructure/database"  // Domain knowing about databases

type User struct {
    ID   string
    conn *database.Connection  // Domain directly using infrastructure
}
```

### Scenario 2: Layer Skipping
```go
// Compiles in Go ✅, but violates layered architecture ❌  
package presentation
import "myapp/infrastructure/repository"  // Skipping application layer

func (h *Handler) GetUser(id string) {
    user := repository.FindUser(id)  // Direct infrastructure access
}
```

### Scenario 3: Framework Contamination
```go
// Compiles in Go ✅, but violates business logic isolation ❌
package business
import "github.com/gin-gonic/gin"  // Business logic tied to web framework

func CalculatePrice(c *gin.Context) float64 {  // Business logic coupled to HTTP
    // calculation logic...
}
```

## Why This Matters

- **Go's import cycle protection** ≠ **Architectural compliance**
- **Compiling code** ≠ **Well-architected code**
- **No circular dependencies** ≠ **Proper dependency direction**

GoArchTest fills this gap by enforcing architectural rules that Go's compiler cannot check.

## Key Benefits

1. **Prevent architectural drift** over time
2. **Enforce design patterns** consistently  
3. **Maintain clean boundaries** between layers
4. **Document architecture** as executable tests
5. **Onboard new team members** with clear rules
6. **Scale development** across multiple teams
