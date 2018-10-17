// Copyright 2015-2018 Kuzzle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"encoding/json"
	"sync"
	"unsafe"

	"github.com/kuzzleio/sdk-go/protocol/websocket"
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
func kuzzle_websocket_new_web_socket(ws *C.web_socket, host *C.char, options *C.options, cppInstance unsafe.Pointer) {
	if websocket_listeners_list == nil {
		websocket_listeners_list = make(map[uintptr]chan<- interface{})
	}

	inst := websocket.NewWebSocket(C.GoString(host), SetOptions(options))

	ws.instance = cppInstance

	registerWebSocket(inst, cppInstance)
}

//export kuzzle_websocket_add_listener
func kuzzle_websocket_add_listener(ws *C.web_socket, event C.int, listener C.kuzzle_event_listener) {
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

			C.bridge_trigger_event_listener(listener, event, C.CString(string(r)), ws.instance)
		}
	}()
}

//export kuzzle_websocket_emit_event
func kuzzle_websocket_emit_event(ws *C.web_socket, event C.int, body *C.char) {
	(*websocket.WebSocket)(ws.instance).EmitEvent(int(event), C.GoString(body))
}
