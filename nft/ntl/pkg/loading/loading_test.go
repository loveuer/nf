package loading

import (
	"context"
	"testing"
	"time"
)

func TestLoadingPrint(t *testing.T) {
	ch := make(chan *Loading)

	Print(context.TODO(), ch)
	ch <- &Loading{Content: "处理中(1)..."}

	time.Sleep(3 * time.Second)

	ch <- &Loading{Content: "处理完成(1)", Type: TypeSuccess}

	ch <- &Loading{Content: "处理中(2)..."}

	time.Sleep(4 * time.Second)

	ch <- &Loading{Content: "处理失败(2)", Type: TypeError}

	time.Sleep(2 * time.Second)
	close(ch)
}
