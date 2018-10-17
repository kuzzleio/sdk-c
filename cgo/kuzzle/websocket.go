package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "kuzzlesdk.h"
	#include "protocol.h"

	static void bridge_trigger_event_listener(kuzzle_event_listener listener, int event, char* res, void* data) {
		listener(event, res, data);
	}
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

//export kuzzle_websocket_new_web_socket
func kuzzle_websocket_new_web_socket(ws *C.web_socket, host *C.char, options *C.options) {
	inst := websocket.NewWebSocket(C.GoString(host), SetOptions(options))

	ptr := unsafe.Pointer(&inst)
	ws.instance = ptr

	registerWebSocket(inst, ptr)
}

//export kuzzle_websocket_add_listener
func kuzzle_websocket_add_listener(ws *C.web_socket, event C.int, listener C.kuzzle_event_listener, data unsafe.Pointer) {
	channel := make(chan interface{})
	go func() {
		for {
			res, ok := <-channel
			if ok == false {
				break
			}

			C.bridge_trigger_event_listener(listener, event, C.CString(res.(string)), data)
		}
	}()
	(*websocket.WebSocket)(ws.instance).AddListener(int(event), channel)
}