package opt

import (
	"fmt"
	"testing"
)

func TestBanner(t *testing.T) {
	fmt.Printf("%s\nnfctl: %s\n\n", Banner, Version)
}
