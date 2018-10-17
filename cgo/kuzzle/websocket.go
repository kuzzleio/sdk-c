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
	"encoding/json"
)

// map which stores instances to keep references in case the gc passes
var webSocketInstances sync.Map

// map which stores channel and function's pointers adresses for listeners
var websocket_listeners_list map[uintptr]chan<- interface{}

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
	c := make(chan interface{})
	websocket_listeners_list[uintptr(unsafe.Pointer(listener))] = c
	(*websocket.WebSocket)(ws.instance).AddListener(int(event), c)
	go func() {
		for {
			res, ok := <-c
			if ok == false {
				break
			}
			r, _ := json.Marshal(res)

			C.bridge_trigger_event_listener(listener, event, C.CString(string(r)), data)
		}
	}()
}

//export kuzzle_websocket_emit_event
func kuzzle_websocket_emit_event(ws *C.web_socket, event C.int, body *C.char) {
	(*websocket.WebSocket)(ws.instance).EmitEvent(int(event), C.GoString(body))
}