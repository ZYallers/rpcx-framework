package safe

import (
	"fmt"
	"runtime/debug"

	"github.com/ZYallers/rpcx-framework/helper/sender"
)

func Defer() {
	r := recover()
	if r == nil {
		return
	}
	msg := fmt.Sprintf("recovery from panic:\n%s\n%s", fmt.Sprintf("%v", r), string(debug.Stack()))
	sender.Graceful(msg, true, "error")
}
