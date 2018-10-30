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
	#include "sdk_wrappers_internal.h"

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
	"github.com/kuzzleio/sdk-go/types"
)

// map which stores instances to keep references in case the gc passes
var webSocketInstances sync.Map

// register new instance to the instances map
func registerWebSocket(instance interface{}, ptr unsafe.Pointer) {
	webSocketInstances.Store(instance, ptr)
}

//export kuzzle_websocket_new_web_socket
func kuzzle_websocket_new_web_socket(ws *C.web_socket, host *C.char, options *C.options, cppInstance unsafe.Pointer) {
	inst := websocket.NewWebSocket(C.GoString(host), SetOptions(options))

	ws.cpp_instance = cppInstance
	ptr := unsafe.Pointer(inst)
	ws.instance = ptr

	registerWebSocket(inst, ptr)
}

//export kuzzle_websocket_add_listener
func kuzzle_websocket_add_listener(ws *C.web_socket, event C.int, listener C.kuzzle_event_listener) {
	c := make(chan interface{})
	go func() {
		for {
			res, ok := <-c
			if ok == false {
				break
			}
			r, _ := json.Marshal(res)

			C.bridge_trigger_event_listener(listener, event, C.CString(string(r)), ws.cpp_instance)
		}
	}()
	(*websocket.WebSocket)(ws.instance).AddListener(int(event), c)
}

//export kuzzle_websocket_emit_event
func kuzzle_websocket_emit_event(ws *C.web_socket, event C.int, body *C.char) {
	(*websocket.WebSocket)(ws.instance).EmitEvent(int(event), nil)
}

//export kuzzle_websocket_connect
func kuzzle_websocket_connect(ws *C.web_socket) *C.char {
	_, err := (*websocket.WebSocket)(ws.instance).Connect()
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("")
}

//export kuzzle_websocket_send
func kuzzle_websocket_send(ws *C.web_socket, query *C.char, options *C.query_options, request_id *C.char) *C.kuzzle_response {
	c := make(chan *types.KuzzleResponse)

	(*websocket.WebSocket)(ws.instance).Send(json.RawMessage(C.GoString(query)), SetQueryOptions(options), c, C.GoString(request_id))

	return goToCKuzzleResponse(<-c)
}

//export kuzzle_websocket_get_state
func kuzzle_websocket_get_state(ws *C.web_socket) C.int {
	return C.int((*websocket.WebSocket)(ws.instance).State())
}

//export kuzzle_websocket_listener_count
func kuzzle_websocket_listener_count(ws *C.web_socket, event int) C.int {
	return C.int((*websocket.WebSocket)(ws.instance).ListenerCount(event))
}

//export kuzzle_websocket_close
func kuzzle_websocket_close(ws *C.web_socket) *C.char {
	err := (*websocket.WebSocket)(ws.instance).Close()
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export kuzzle_websocket_unregister_sub
func kuzzle_websocket_unregister_sub(ws *C.web_socket, id *C.char) {
	(*websocket.WebSocket)(ws.instance).UnregisterSub(C.GoString(id))
}

//export kuzzle_websocket_cancel_subs
func kuzzle_websocket_cancel_subs(ws *C.web_socket) {
	(*websocket.WebSocket)(ws.instance).CancelSubs()
}

//export kuzzle_websocket_start_queuing
func kuzzle_websocket_start_queuing(ws *C.web_socket) {
	(*websocket.WebSocket)(ws.instance).StartQueuing()
}

//export kuzzle_websocket_stop_queuing
func kuzzle_websocket_stop_queuing(ws *C.web_socket) {
	(*websocket.WebSocket)(ws.instance).StopQueuing()
}

//export kuzzle_websocket_play_queue
func kuzzle_websocket_play_queue(ws *C.web_socket) {
	(*websocket.WebSocket)(ws.instance).PlayQueue()
}

//export kuzzle_websocket_clear_queue
func kuzzle_websocket_clear_queue(ws *C.web_socket) {
	(*websocket.WebSocket)(ws.instance).ClearQueue()
}

//export kuzzle_websocket_remove_all_listeners
func kuzzle_websocket_remove_all_listeners(ws *C.web_socket, event C.int) {
	(*websocket.WebSocket)(ws.instance).RemoveAllListeners(int(event))
}