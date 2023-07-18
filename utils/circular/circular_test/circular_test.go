package circular_test

import (
	"testing"

	"github.com/cwloo/gonet/utils/circular"
)

func TestMain(m *testing.M) {
	m.Run()
}

func circular_test(t *testing.T) {
	circular.Test001()
	circular.Test002()
	circular.Test003()
	circular.Test004()
}

func Test(t *testing.T) {
	t.Run("circular.Test001", circular_test)
}
