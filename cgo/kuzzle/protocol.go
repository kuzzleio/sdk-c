package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "kuzzlesdk.h"
*/
import "C"
import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/kuzzleio/sdk-go/types"
)

type WrapProtocol struct {
	P *C.protocol
}

func NewWrapProtocol(p *C.protocol) *WrapProtocol {
	return &WrapProtocol{p}
}

func (wp WrapProtocol) AddListener(event int, channel chan<- interface{}) {
	c := func(event C.int, res *C.char, data unsafe.Pointer) {
		(chan<- interface{})(channel) <- res
	}

	wp.P.add_listener(C.int(event), (*C.kuzzle_event_listener)(c))
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- interface{}) {
	// wp.P.remove_listener(C.int(event), (*C.Kuzzle_event_listener)(C.cgo_listener_callback))
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	wp.P.remove_all_listeners(C.int(event))
}
func (wp WrapProtocol) Once(event int, channel chan<- interface{}) {

}

func (wp WrapProtocol) ListenerCount(event int) int {
	return int(wp.P.listener_count(C.int(event)))
}

func (wp WrapProtocol) Connect() (bool, error) {
	if err := wp.P.connect(); err != nil {
		return false, error.New(C.GoString(err))
	}

	return true, nil
}

func (wp WrapProtocol) Send(query []byte, options types.QueryOptions, responseChannel chan<- *types.KuzzleResponse, requestId string) error {
	/*
		if err := wp.P.send(C.CString(query), ); err != nil {
			return error.New(C.GoString(err))
		}

		return nil
	*/
	return nil
}

func (wp WrapProtocol) Close() error {
	if err := wp.P.close(); err != nil {
		return error.New(C.GoString(err))
	}

	return nil
}

func (wp WrapProtocol) State() int {
	return int(wp.P.state())
}

func (wp WrapProtocol) EmitEvent(event int, data interface{}) {
	wp.P.emit_event(C.int(event), unsafe.Pointer(data))
}

func (wp WrapProtocol) RegisterSub(string, string, json.RawMessage, bool, chan<- types.KuzzleNotification, chan<- interface{}) {

}

func (wp WrapProtocol) UnregisterSub(id string) {
	wp.P.unregister_sub(C.CString(id))
}

func (wp WrapProtocol) CancelSubs() {
	wp.P.cancel_subs()
}

func (wp WrapProtocol) RequestHistory() map[string]time.Time {
	return nil
}

func (wp WrapProtocol) StartQueuing() {
	wp.P.start_queuing()
}

func (wp WrapProtocol) StopQueuing() {
	wp.P.stop_queuing()
}

func (wp WrapProtocol) PlayQueue() {
	wp.P.play_queue()
}

func (wp WrapProtocol) ClearQueue() {
	wp.P.clear_queue()
}
