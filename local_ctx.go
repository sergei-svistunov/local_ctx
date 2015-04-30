package local_ctx

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

const addrPrefix = "github.com/sergei-svistunov/local_ctx.(*lCtx).call("

type lCtx struct {
	parent *lCtx
	data   interface{}
}

func (c *lCtx) call(f func()) {
	f()
}

func Call(ctx interface{}, f func()) {
	lctx := &lCtx{data: ctx}
	lctx.call(f)
}

func Go(f func()) {
	pctx, err := getCtx()
	if err != nil {
		panic(err)
	}
	go func(pctx *lCtx) {
		lctx := &lCtx{parent: pctx}
		lctx.call(f)
	}(pctx)
}

func Data() (interface{}, error) {
	ctx, err := getCtx()
	if err != nil {
		return nil, err
	}

	return ctx.data, nil
}

func getCtx() (*lCtx, error) {
	buf := make([]byte, 128)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, len(buf)*2)
	}

	strStack := string(buf)
	addrPos := strings.Index(strStack, addrPrefix) + len(addrPrefix)
	if addrPos < 0 {
		return nil, fmt.Errorf("No local context")
	}

	addr, err := strconv.ParseUint(strStack[addrPos:addrPos+12], 0, 64)
	if err != nil {
		panic(fmt.Sprintf("Cannot get address: %s", err))
	}

	ptr := unsafe.Pointer(uintptr(addr))
	ctx := reflect.NewAt(reflect.TypeOf(lCtx{}), ptr).Interface().(*lCtx)
	
	for ctx.parent != nil {
		ctx = ctx.parent
	}
	
	return ctx, nil
}
