# Implicit Interfacing in Go - Testing Examples

## Live Demos

### 1. API Calls
- Mock HTTP clients for testing without real network calls
- Demonstrate using interfaces to inject test doubles
- Show how to verify request/response handling

### 2. Time Dependencies (time.Sleep)
- Make time-dependent code testable without waiting
- Mock time operations for fast, deterministic tests
- Control time flow in test scenarios

### 3. Random Number Generation
- Make non-deterministic code deterministic in tests
- Mock random number sources for predictable outcomes
- Verify behavior with specific random sequences

## Additional Use Cases

### File System Operations
- Mock file reading/writing operations
- Test without touching the actual filesystem
- Use `io.Reader`/`io.Writer` interfaces for flexibility

### Database Operations
- Create repository patterns with interface abstractions
- Swap real database connections with test mocks
- Test business logic without database setup

### External Command Execution
- Mock `exec.Command` for testing system interactions
- Test command execution logic without running actual commands
- Verify command arguments and environment setup

### Environment Variables
- Wrap `os.Getenv` behind testable interfaces
- Test different configurations without modifying environment
- Control configuration in isolated tests

### Logger/Output Operations
- Inject different loggers for testing
- Capture and verify log output in tests
- Mock `io.Writer` to validate logging behavior

## Key Principles

- **Define narrow interfaces** - Only include methods you actually need
- **Accept interfaces, return structs** - Maximum flexibility for consumers
- **Interfaces belong to the consumer** - Define them where they're used
- **Implicit satisfaction** - No explicit `implements` keyword needed
