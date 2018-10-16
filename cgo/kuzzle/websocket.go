package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "kuzzlesdk.h"
	#include "protocol.h"
*/
import "C"
import (
	"unsafe"
	"sync"
	"github.com/kuzzleio/sdk-go/protocol/websocket"
)

// map which stores instances to keep references in case the gc passes
var webSocketInstances sync.Map

// register new instance to the instances map
func registerWebSocket(instance interface{}, ptr unsafe.Pointer) {
	webSocketInstances.Store(instance, ptr)
}

// export kuzzle_new_websocket
func kuzzle_new_websocket(ws *C.web_socket, host *C.char, options *C.options) {
	inst := websocket.NewWebSocket(C.GoString(host), SetOptions(options))

	ptr := unsafe.Pointer(&inst)
	ws.instance = ptr

	registerWebSocket(inst, ptr)
}

