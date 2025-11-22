---
author: Gopher
date: MMMM dd, YYYY
paging: Slide %d / %d
---

# Implicit Interfacing in Go

## Testing Without the Pain

*A practical guide to mocking dependencies*

---

## The Problem

Testing code that depends on:
- **Random generators** - non-deterministic, hard to verify
- **Time/sleep** - tests that wait are slow tests
- **HTTP APIs** - slow, unreliable, requires network

**Solution**: Use Go's implicit interfaces to inject test doubles

---

## What Makes Go Special?

```go
// Define interface where you USE it, not where it's implemented
type StringGenerator interface {
    Generate() string
}

// Any type with Generate() string satisfies this - no "implements" keyword!
```

**Key insight**: Accept interfaces, return structs

---

## Demo 1: Random String Generator

**The Problem**: How do you test something that returns different values every time?

Press **→** to see the code WITHOUT interfaces

---

## Random Generator: ❌ Without Interfaces

```go
type NameGenerator struct{}

// PROBLEM: Directly uses rand.Intn - can't control in tests!
func (n *NameGenerator) GenerateName(baseName string) string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, 5)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]  // Hard-coded!
    }
    suffix := string(b)
    return fmt.Sprintf("%s-%s", baseName, suffix)
}
```

**How do you test this?** Every test run generates different output!

You can't easily verify the exact output without complicated workarounds.

---

## Random Generator: ✅ With Interfaces

```go
type StringGenerator interface {
    Generate() string
}

type NameGenerator struct {
    generator StringGenerator  // Injected dependency!
}

// Now testable - we control the generator
func (n *NameGenerator) GenerateName(baseName string) string {
    suffix := n.generator.Generate()
    return fmt.Sprintf("%s-%s", baseName, suffix)
}

// Production: real randomness
type RandomStringGenerator struct{}

func (r RandomStringGenerator) Generate() string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, 5)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
```

---

## Random Generator: Run Production

**In production**, names are truly random

```bash
cd demos/rng && go run main.go
```

*Press* **Ctrl+E** *to see random names generated!*

---

## Random Generator: Test Code

```go
type mockStringGenerator struct {
    value string
}

func (m mockStringGenerator) Generate() string {
    return m.value  // Always returns the same value!
}

func TestGenerateName(t *testing.T) {
    tests := map[string]struct {
        baseName string
        suffix   string
        expected string
    }{
        "pod name": {
            baseName: "my-pod",
            suffix:   "abc12",
            expected: "my-pod-abc12",
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            gen := NewNameGenerator(mockStringGenerator{tc.suffix})
            result := gen.GenerateName(tc.baseName)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

---

## Random Generator: Run Tests

**In tests**, we control the "random" value for deterministic results

```bash
cd demos/rng && go test -v
```

*Press* **Ctrl+E** *to see deterministic tests pass!*

---

## Demo 2: Time/Sleep Mocking

**The Problem**: Tests that sleep are slow tests

A 2-second sleep = 2-second test... *unless we mock it*

Press **→** to see the code WITHOUT interfaces

---

## Time/Sleep: ❌ Without Interfaces

```go
type Worker struct{}

// PROBLEM: Directly calls time.Sleep - tests must wait!
func (w *Worker) DoWork() string {
    fmt.Println("Starting work...")
    time.Sleep(2 * time.Second)  // Hard-coded!
    fmt.Println("Work complete!")
    return "done"
}
```

**How do you test this?** Every test takes 2+ seconds!

Your test suite gets slower and slower as you add more tests.

---

## Time/Sleep: ✅ With Interfaces

```go
type Sleeper interface {
    Sleep(duration time.Duration)
}

type Worker struct {
    sleeper Sleeper  // Injected dependency!
}

// Now testable - we control the sleeper
func (w *Worker) DoWork() string {
    fmt.Println("Starting work...")
    w.sleeper.Sleep(2 * time.Second)
    fmt.Println("Work complete!")
    return "done"
}

// Production: actually sleeps
type realSleeper struct{}

func (r realSleeper) Sleep(d time.Duration) {
    time.Sleep(d)
}
```

---

## Time/Sleep: Run Production

**In production**, it really sleeps for 2 seconds

```bash
cd demos/sleeper && go run main.go
```

*Press* **Ctrl+E** *and watch the 2-second pause with timestamps!*

---

## Time/Sleep: Test Code

```go
type mockSleeper struct {
    duration time.Duration
}

func (m *mockSleeper) Sleep(d time.Duration) {
    m.duration = d  // Record but DON'T actually sleep!
}

func TestDoWork(t *testing.T) {
    mock := &mockSleeper{}
    worker := NewWorker(mock)

    result := worker.DoWork()

    assert.Equal(t, "done", result)
    assert.Equal(t, 2*time.Second, mock.duration)
}
```

**The magic**: Test verifies the sleep duration without waiting!

---

## Time/Sleep: Run Tests

**In tests**, we verify sleep was called without actually sleeping

```bash
cd demos/sleeper && go test -v
```

*Press* **Ctrl+E** *- Notice how fast it completes (milliseconds, not seconds)!*

---

## Demo 3: HTTP Client Mocking

**The Problem**: Testing API calls without hitting real servers

Press **→** to see the code WITHOUT interfaces

---

## HTTP Client: ❌ Without Interfaces

```go
type User struct {
    Name string `json:"name"`
}

// PROBLEM: Directly uses http.DefaultClient - hits real servers!
func GetUser(id int) (*User, error) {
    resp, err := http.DefaultClient.Get(  // Hard-coded!
        fmt.Sprintf("https://api.example.com/users/%d", id))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}
```

**How do you test this?** You need a real server running!

Tests are slow, flaky, and require network infrastructure.

---

## HTTP Client: ✅ With Interfaces

```go
type HTTPClient interface {
    Get(url string) (*http.Response, error)
}

type User struct {
    Name string `json:"name"`
}

// Now testable - we control the HTTP client
func GetUser(client HTTPClient, id int) (*User, error) {
    resp, err := client.Get(fmt.Sprintf(
        "https://api.example.com/users/%d", id))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }

    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

## HTTP Client: Run Production

**In production**, you pass a real `*http.Client`

```bash
cd demos/api-call && go run main.go
```

*Press* **Ctrl+E** *to run the production code!*

---

## HTTP Client: Test Code

```go
type mockHTTPClient struct {
    statusCode int
    body       string
    err        error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
    if m.err != nil {
        return nil, m.err
    }
    return &http.Response{
        StatusCode: m.statusCode,
        Body:       io.NopCloser(strings.NewReader(m.body)),
    }, nil
}

func TestGetUser(t *testing.T) {
    tests := map[string]struct {
        mock           *mockHTTPClient
        expectedUser   *User
        expectedErrMsg string
    }{
        "success": {
            mock:         &mockHTTPClient{200, `{"name":"Alice"}`, nil},
            expectedUser: &User{Name: "Alice"},
        },
        "not found": {
            mock:           &mockHTTPClient{404, "", nil},
            expectedErrMsg: "unexpected status: 404",
        },
    }
    // ... test each case with mock
}
```

---

## HTTP Client: Run Tests

**In tests**, we inject a mock that never hits the network

```bash
cd demos/api-call && go test -v
```

*Press* **Ctrl+E** *to run the tests - instant, no network needed!*

---

## The Transformation

For each demo, we saw the same pattern:

1. **❌ Before**: Hard-coded dependencies → Hard to test
2. **✅ After**: Interface injection → Easy to test

**The key change**: Accept an interface parameter instead of using globals

---

## Key Principles

✅ **Define narrow interfaces** - Only include methods you need

✅ **Accept interfaces, return structs** - Maximum flexibility

✅ **Interfaces belong to the consumer** - Define where used, not where implemented

✅ **Implicit satisfaction** - No `implements` keyword needed

✅ **Small interfaces are better** - Easier to mock and maintain

---

## Beyond These Examples

This pattern works for:
- **File system operations** - `io.Reader`, `io.Writer`
- **Database operations** - Repository patterns
- **Command execution** - Wrap `exec.Command`
- **Environment variables** - Wrap `os.Getenv`
- **Logging** - Inject logger interfaces

---

## Why This Matters

**Without implicit interfaces**:
- Slow tests (waiting for network, time, etc.)
- Brittle tests (dependent on external services)
- Hard to test edge cases

**With implicit interfaces**:
- ✅ Fast, deterministic tests
- ✅ Easy to test error conditions
- ✅ No external dependencies needed

---

## Live Demo Summary

1. **Random Generator** - `demos/rng` - Deterministic randomness
2. **Time/Sleep** - `demos/sleeper` - Instant time travel
3. **HTTP Client** - `demos/api-call` - Mock API calls

**All tests run in milliseconds!**

---

## Questions?

### Navigation Tips:
- **Space** / **→** - Next slide
- **←** - Previous slide
- **gg** - First slide
- **G** - Last slide
- **/** - Search
- **Ctrl+E** - Execute code block

### Code:
All examples available in `demos/` directory

---

# Thank You!

*Go forth and mock responsibly*
