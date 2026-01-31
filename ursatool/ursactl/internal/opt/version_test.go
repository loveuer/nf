package opt

import (
	"fmt"
	"testing"
)

func TestBanner(t *testing.T) {
	fmt.Printf("%s\nursactl: %s\n\n", Banner, Version)
}
