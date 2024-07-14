package version

import (
	"github.com/loveuer/nf/nft/log"
	"testing"
)

func TestUpgradePrint(t *testing.T) {
	UpgradePrint("v24.07.14-r5")
}

func TestCheck(t *testing.T) {
	log.SetLogLevel(log.LogLevelDebug)
	v := Check(15)
	t.Logf("got version: %s", v)
}
