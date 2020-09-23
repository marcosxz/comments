// @Package
package comments

import (
	"os"
	"testing"
	"time"
)

// @Comment1 1
// @Comment2 2
// @Comment3 3
type test struct {
	// Field1
	Field1 string
	// Field2
	Field2 int
	// Field3
	Field3 bool
	// Field4
	Field4 float64
	// Field5
	Field5 time.Time
}

// @FuncName FuncA
func (test) FuncA(i int) error {
	return nil
}

func TestAstFind(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err = AstFind(dir, true, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
