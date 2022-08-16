package pool

import (
	"testing"
	"transaction-matching-engine/models"
)

func BenchmarkKeyPtrSpeed(b *testing.B) {
	m := models.SortKey{}
	for i := 0; i < b.N; i++ {
		convertPtr(&m)
	}
}

func BenchmarkKeyModelSpeed(b *testing.B) {
	m := models.SortKey{}
	for i := 0; i < b.N; i++ {
		convertModel(m)
	}
}
func convertPtr(a interface{}) {
	_ = a.(*models.SortKey)
}
func convertModel(a interface{}) {
	_ = a.(models.SortKey)
}
