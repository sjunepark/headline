package benchmark

import (
	"fmt"
	"github.com/cockroachdb/errors"
	cdbErrors "github.com/cockroachdb/errors"
	"testing"
)

func BenchmarkStandardErrors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errors.New("test error")
	}
}

func BenchmarkCockroachdbErrors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cdbErrors.New("test error")
	}
}

func BenchmarkStandardErrors_Wrap(b *testing.B) {
	original := errors.New("original error")
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("wrapped: %w", original)
	}
}

func BenchmarkCockroachdbErrors_Wrap(b *testing.B) {
	original := cdbErrors.New("original error")
	for i := 0; i < b.N; i++ {
		_ = cdbErrors.Wrap(original, "wrapped")
	}
}
