package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "protocol.h"

  // void (*add_listener)(int, kuzzle_event_listener*);
  // void (*remove_listener)(int, kuzzle_event_listener*);
  // void (*once)(int, kuzzle_event_listener*);
  // char* (*send)(const char*, query_options*, kuzzle_response*, char*);
  // void (*emit_event)(int, void*);
  // void (*register_sub)(const char*, const char*, const char*, int, kuzzle_notification_listener*, void*);

	// void (*clear_queue)();

	// Bridges

	static void bridge_add_listener(void (*f)(int, kuzzle_event_listener*), int event, kuzzle_event_listener* listener) {
		f(event, listener);
	}

	static void bridge_remove_all_listeners(void (*f)(int), int event) {
		f(event);
	}

	static int bridge_listener_count(int (*f)(int), int event) {
		return f(event);
	}

	static char* bridge_connect(char* (*f)()) {
		return f();
	}

	static char* bridge_close(char* (*f)()) {
		return f();
	}

	static int bridge_get_state(int (*f)()) {
		return f();
	}

	static void bridge_emit_event(void (*f)(int, void*), int event, void* data) {
		f(event, data);
	}

	static void bridge_unregister_sub(void (*f)(char*), char* id) {
		f(id);
	}

	static void bridge_cancel_subs(void (*f)()) {
		f();
	}

	static void bridge_start_queuing(void (*f)()) {
		f();
	}

	static void bridge_stop_queuing(void (*f)()) {
		f();
	}

	static void bridge_play_queue(void (*f)()) {
		f();
	}

	static void bridge_clear_queue(void (*f)()) {
		f();
	}

	void goAddListenerCallback(char*, void*);

	static void bridge_kuzzle_event_listener(char* res, void* chan) {
		goAddListenerCallback(res, chan);
	}
*/
import "C"
import (
	"encoding/json"
	"errors"
	"time"
	"unsafe"

	"github.com/kuzzleio/sdk-go/protocol"
	"github.com/kuzzleio/sdk-go/types"
)

type WrapProtocol struct {
	P *C.protocol
}

func NewWrapProtocol(p *C.protocol) *WrapProtocol {
	return &WrapProtocol{p}
}

//export goAddListenerCallback
func goAddListenerCallback(res *C.char, channel unsafe.Pointer) {
	// (chan<- interface{})(channel) <- res
}

func (wp WrapProtocol) AddListener(event int, channel chan<- interface{}) {
	// C.bridge_add_listener(wp.P.add_listener, C.int(event), C.bridge_kuzzle_event_listener)
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- interface{}) {
	// wp.P.remove_listener(C.int(event), (*C.Kuzzle_event_listener)(C.cgo_listener_callback))
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	C.bridge_remove_all_listeners(wp.P.remove_all_listeners, C.int(event))
}

func (wp WrapProtocol) Once(event int, channel chan<- interface{}) {

}

func (wp WrapProtocol) ListenerCount(event int) int {
	return int(C.bridge_listener_count(wp.P.listener_count, C.int(event)))
}

func (wp WrapProtocol) Connect() (bool, error) {
	if err := C.bridge_connect(wp.P.connect); err != nil {
		return false, errors.New(C.GoString(err))
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
	if err := C.bridge_close(wp.P.close); err != nil {
		return errors.New(C.GoString(err))
	}

	return nil
}

func (wp WrapProtocol) State() int {
	return int(C.bridge_get_state(wp.P.get_state))
}

func (wp WrapProtocol) EmitEvent(event int, data interface{}) {
	C.bridge_emit_event(wp.P.emit_event, C.int(event), data.(unsafe.Pointer))
}

func (wp WrapProtocol) RegisterSub(string, string, json.RawMessage, bool, chan<- types.KuzzleNotification, chan<- interface{}) {

}

func (wp WrapProtocol) UnregisterSub(id string) {
	C.bridge_unregister_sub(wp.P.unregister_sub, C.CString(id))
}

func (wp WrapProtocol) CancelSubs() {
	C.bridge_cancel_subs(wp.P.cancel_subs)
}

func (wp WrapProtocol) RequestHistory() map[string]time.Time {
	return nil
}

func (wp WrapProtocol) StartQueuing() {
	C.bridge_start_queuing(wp.P.start_queuing)
}

func (wp WrapProtocol) StopQueuing() {
	C.bridge_stop_queuing(wp.P.stop_queuing)
}

func (wp WrapProtocol) PlayQueue() {
	C.bridge_play_queue(wp.P.play_queue)
}

func (wp WrapProtocol) ClearQueue() {
	C.bridge_clear_queue(wp.P.clear_queue)
}

func (wp WrapProtocol) AutoQueue() bool {
	if int(wp.P.auto_queue) == 1 {
		return true
	}
	return false
}

func (wp WrapProtocol) AutoReconnect() bool {
	if int(wp.P.auto_reconnect) == 1 {
		return true
	}
	return false
}

func (wp WrapProtocol) AutoResubscribe() bool {
	if int(wp.P.auto_resubscribe) == 1 {
		return true
	}
	return false
}

func (wp WrapProtocol) AutoReplay() bool {
	if int(wp.P.auto_replay) == 1 {
		return true
	}
	return false
}

func (wp WrapProtocol) Host() string {
	return C.GoString(wp.P.host)
}

func (wp WrapProtocol) OfflineQueue() []*types.QueryObject {
	tmpslice := (*[1<<28 - 1]*C.query_object)(unsafe.Pointer(wp.P.kuzzle_offline_queue.queries))[:wp.P.kuzzle_offline_queue.queries_length]
	goOfflineQueue := make([]*types.QueryObject, 0, int(wp.P.kuzzle_offline_queue.queries_length))

	for _, s := range tmpslice {
		cQuery := cToGoQueryObject(s, nil)
		goOfflineQueue = append(goOfflineQueue, cQuery)
	}

	return goOfflineQueue
}

func (wp WrapProtocol) OfflineQueueLoader() protocol.OfflineQueueLoader {
	var offline protocol.OfflineQueueLoader
	return offline
}

func (wp WrapProtocol) Port() int {
	return -1
}

func (wp WrapProtocol) QueueFilter() protocol.QueueFilter {
	return func([]byte) bool { return false }
}

func (wp WrapProtocol) QueueMaxSize() int {
	return -1
}

func (wp WrapProtocol) QueueTTL() time.Duration {
	return time.Duration(1)
}

func (wp WrapProtocol) ReplayInterval() time.Duration {
	return time.Duration(1)
}

func (wp WrapProtocol) ReconnectionDelay() time.Duration {
	return time.Duration(1)
}

func (wp WrapProtocol) SslConnection() bool {
	return false
}

func (wp WrapProtocol) SetAutoQueue(value bool)                                 {}
func (wp WrapProtocol) SetAutoReplay(value bool)                                {}
func (wp WrapProtocol) SetOfflineQueueLoader(value protocol.OfflineQueueLoader) {}
func (wp WrapProtocol) SetQueueFilter(value protocol.QueueFilter)               {}
func (wp WrapProtocol) SetQueueMaxSize(value int)                               {}
func (wp WrapProtocol) SetQueueTTL(time.Duration)                               {}
func (wp WrapProtocol) SetReplayInterval(time.Duration)                         {}
