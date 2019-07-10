package crypto-hug
import (
	"testing"
)

func TestHello(t *testing.T) {
	str := Hello()
	assert.Equal(t, "hello world", str)
}