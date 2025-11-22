package main

import (
	"fmt"
	"math/rand"
	"time"
)

// StringGenerator interface generates random strings
type StringGenerator interface {
	Generate() string
}

// NameGenerator creates Kubernetes-style names with random suffixes
type NameGenerator struct {
	generator StringGenerator
}

// GenerateName creates a name by appending a random suffix to the base name
func (ng *NameGenerator) GenerateName(baseName string) string {
	suffix := ng.generator.Generate()
	return fmt.Sprintf("%s-%s", baseName, suffix)
}

// RandomStringGenerator generates random alphanumeric strings
type RandomStringGenerator struct {
	rng *rand.Rand
}

// Generate creates a random 5-character alphanumeric string
func (r *RandomStringGenerator) Generate() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const suffixLen = 5

	suffix := make([]byte, suffixLen)
	for i := range suffix {
		suffix[i] = charset[r.rng.Intn(len(charset))]
	}
	return string(suffix)
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	strGen := &RandomStringGenerator{rng: rng}
	nameGen := &NameGenerator{generator: strGen}

	// Generate some k8s-style names
	fmt.Println("Generated names:")
	fmt.Printf("  Pod:        %s\n", nameGen.GenerateName("my-pod"))
	fmt.Printf("  Deployment: %s\n", nameGen.GenerateName("nginx-deployment"))
	fmt.Printf("  Service:    %s\n", nameGen.GenerateName("api-service"))
}
