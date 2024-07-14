package version

import (
	"github.com/loveuer/nf/nft/log"
	"testing"
)

func TestUpgradePrint(t *testing.T) {
	log.SetLogLevel(log.LogLevelDebug)
	UpgradePrint("v24.07.14-r5")
}

func TestCheck(t *testing.T) {
	log.SetLogLevel(log.LogLevelDebug)
	v := Check(true, true, 1)
	t.Logf("got version: %s", v)
}
