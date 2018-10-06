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
	p *C.protocol
}

func (wp *WrapProtocol) NewWrapProtocol(p *C.protocol) {
	wp.p = p
}

func (wp WrapProtocol) AddListener(event int, channel chan<- interface{}) {
	c := func(event C.int, res *C.char, data unsafe.Pointer) {
		(chan<- interface{})(chan_res) <- res
	}

	wp.p.add_listener(C.int(event), (*C.kuzzle_event_listener)(c))
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- interface{}) {
	// wp.p.remove_listener(C.int(event), (*C.Kuzzle_event_listener)(C.cgo_listener_callback))
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	wp.p.remove_all_listeners(C.int(event))
}
func (wp WrapProtocol) Once(event int, channel chan<- interface{}) {

}

func (wp WrapProtocol) ListenerCount(event int) int {
	return int(wp.p.listener_count(C.int(event)))
}

func (wp WrapProtocol) Connect() (bool, error) {
	if err := wp.p.connect(); err != nil {
		return false, error.New(C.GoString(err))
	}

	return true, nil
}

func (wp WrapProtocol) Send(query []byte, options types.QueryOptions, responseChannel chan<- *types.KuzzleResponse, requestId string) error {
	/*
		if err := wp.p.send(C.CString(query), ); err != nil {
			return error.New(C.GoString(err))
		}

		return nil
	*/
	return nil
}

func (wp WrapProtocol) Close() error {
	if err := wp.p.close(); err != nil {
		return error.New(C.GoString(err))
	}

	return nil
}

func (wp WrapProtocol) State() int {
	return int(wp.p.state())
}

func (wp WrapProtocol) EmitEvent(event int, data interface{}) {
	wp.p.emit_event(C.int(event), unsafe.Pointer(data))
}

func (wp WrapProtocol) RegisterSub(string, string, json.RawMessage, bool, chan<- types.KuzzleNotification, chan<- interface{}) {

}

func (wp WrapProtocol) UnregisterSub(id string) {
	wp.p.unregister_sub(C.CString(id))
}

func (wp WrapProtocol) CancelSubs() {
	wp.p.cancel_subs()
}

func (wp WrapProtocol) RequestHistory() map[string]time.Time {
	return nil
}

func (wp WrapProtocol) StartQueuing() {
	wp.p.start_queuing()
}

func (wp WrapProtocol) StopQueuing() {
	wp.p.stop_queuing()
}

func (wp WrapProtocol) PlayQueue() {
	wp.p.play_queue()
}

func (wp WrapProtocol) ClearQueue() {
	wp.p.clear_queue()
}
