package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "protocol.h"

	// Bridges

	static void bridge_protocol_add_listener(void (*f)(int, kuzzle_event_listener*, void*), int event, kuzzle_event_listener* listener, void* data) {
		f(event, listener, data);
	}

	static void bridge_remove_listener(void (*f)(int, kuzzle_event_listener*, void*), int event, kuzzle_event_listener* listener, void* data) {
		f(event, listener, data);
	}

	static void bridge_remove_all_listeners(void (*f)(int, void*), int event, void* data) {
		f(event, data);
	}

	static void bridge_once(void (*f)(int, kuzzle_event_listener*, void*), int event, kuzzle_event_listener* listener, void* data) {
		f(event, listener, data);
	}

	static int bridge_listener_count(int (*f)(int, void*), int event, void* data) {
		return f(event, data);
	}

	static char* bridge_connect(char* (*f)(void*), void* data) {
		return f(data);
	}

	static kuzzle_response* bridge_send(kuzzle_response* (*f)(const char*, query_options*, char*, void*), const char* query, query_options* options, char* request_id, void* data) {
		return f(query, options, request_id, data);
	}

	static char* bridge_close(char* (*f)(void*), void* data) {
		return f(data);
	}

	static int bridge_get_state(int (*f)(void*), void* data) {
		return f(data);
	}

	static void bridge_emit_event(void (*f)(int, void*, void*), int event, void* res, void* data) {
		f(event, res, data);
	}

	static void bridge_unregister_sub(void (*f)(char*, void*), char* id, void* data) {
		f(id, data);
	}

	static void bridge_cancel_subs(void (*f)(void*), void* data) {
		f(data);
	}

	static void bridge_start_queuing(void (*f)(void*), void* data) {
		f(data);
	}

	static void bridge_stop_queuing(void (*f)(void*), void* data) {
		f(data);
	}

	static void bridge_play_queue(void (*f)(void*), void* data) {
		f(data);
	}

	static void bridge_clear_queue(void (*f)(void*), void* data) {
		f(data);
	}

	static bool bridge_queue_filter(kuzzle_queue_filter f, const char* data) {
		return f(data);
	}

	extern void bridge(int, char*, void*);

	static void call_bridge(int event, char* res, void* data) {
		bridge(event, res, data);
	}

	static kuzzle_event_listener get_fptr() {
		return &call_bridge;
	}
*/
import "C"
import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"unsafe"
	"fmt"

	"github.com/kuzzleio/sdk-go/protocol"
	"github.com/kuzzleio/sdk-go/types"
)

type WrapProtocol struct {
	P *C.protocol
}

var proto_instances sync.Map

// register new instance to the instances map
func registerProtocol(instance interface{}, ptr unsafe.Pointer) {
	proto_instances.Store(instance, ptr)
}

//export unregisterProtocol
func unregisterProtocol(p *C.protocol) {
	proto_instances.Delete(p)
}

func NewWrapProtocol(p *C.protocol) *WrapProtocol {
	ptr_proto := unsafe.Pointer(p.instance)
	registerProtocol(p, ptr_proto)

	return &WrapProtocol{p}
}

//export bridge
func bridge(event C.int, res *C.char, channel unsafe.Pointer) {
	fmt.Printf("okkkkkkkkk\n")
}

func (wp WrapProtocol) AddListener(event int, channel chan<- interface{}) {
	fptr := C.get_fptr()
	C.bridge_protocol_add_listener(wp.P.add_listener, C.int(event), &fptr, wp.P.instance)
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- interface{}) {
	C.bridge_remove_listener(wp.P.remove_listener, C.int(event), (*C.kuzzle_event_listener)(unsafe.Pointer(&channel)), wp.P.instance)
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	C.bridge_remove_all_listeners(wp.P.remove_all_listeners, C.int(event), wp.P.instance)
}

func (wp WrapProtocol) Once(event int, channel chan<- interface{}) {
	C.bridge_once(wp.P.add_listener, C.int(event), (*C.kuzzle_event_listener)(unsafe.Pointer(&channel)), wp.P.instance)
}

func (wp WrapProtocol) ListenerCount(event int) int {
	return int(C.bridge_listener_count(wp.P.listener_count, C.int(event), wp.P.instance))
}

func (wp WrapProtocol) Connect() (bool, error) {
	if err := C.bridge_connect(wp.P.connect, wp.P.instance); err != nil {
		return false, errors.New(C.GoString(err))
	}

	return true, nil
}

func (wp WrapProtocol) Send(query []byte, options types.QueryOptions, responseChannel chan<- *types.KuzzleResponse, requestId string) error {
	res := C.bridge_send(wp.P.send, C.CString(string(query[:])), goToCQueryOptions(options), C.CString(requestId), wp.P.instance)
	if res.error != nil {
		return errors.New(C.GoString(res.error))
	}
	if responseChannel != nil {
		responseChannel <- cToGoKuzzleResponse(res)
	}
	return nil
}

func (wp WrapProtocol) Close() error {
	if err := C.bridge_close(wp.P.close, wp.P.instance); err != nil {
		return errors.New(C.GoString(err))
	}

	return nil
}

func (wp WrapProtocol) State() int {
	return int(C.bridge_get_state(wp.P.get_state, wp.P.instance))
}

func (wp WrapProtocol) EmitEvent(event int, data interface{}) {
	C.bridge_emit_event(wp.P.emit_event, C.int(event), nil, wp.P.instance)
}

func (wp WrapProtocol) RegisterSub(string, string, json.RawMessage, bool, chan<- types.KuzzleNotification, chan<- interface{}) {
	//@todo
}

func (wp WrapProtocol) UnregisterSub(id string) {
	C.bridge_unregister_sub(wp.P.unregister_sub, C.CString(id), wp.P.instance)
}

func (wp WrapProtocol) CancelSubs() {
	C.bridge_cancel_subs(wp.P.cancel_subs, wp.P.instance)
}

func (wp WrapProtocol) RequestHistory() map[string]time.Time {
	//@todo
	return nil
}

func (wp WrapProtocol) StartQueuing() {
	C.bridge_start_queuing(wp.P.start_queuing, wp.P.instance)
}

func (wp WrapProtocol) StopQueuing() {
	C.bridge_stop_queuing(wp.P.stop_queuing, wp.P.instance)
}

func (wp WrapProtocol) PlayQueue() {
	C.bridge_play_queue(wp.P.play_queue, wp.P.instance)
}

func (wp WrapProtocol) ClearQueue() {
	C.bridge_clear_queue(wp.P.clear_queue, wp.P.instance)
}

func (wp WrapProtocol) AutoQueue() bool {
	return bool(wp.P.auto_queue)
}

func (wp WrapProtocol) AutoReconnect() bool {
	return bool(wp.P.auto_reconnect)
}

func (wp WrapProtocol) AutoResubscribe() bool {
	return bool(wp.P.auto_resubscribe)
}

func (wp WrapProtocol) AutoReplay() bool {
	return bool(wp.P.auto_replay)
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
	// @todo
	var offline protocol.OfflineQueueLoader
	return offline
}

func (wp WrapProtocol) Port() int {
	return int(wp.P.port)
}

func (wp WrapProtocol) QueueFilter() protocol.QueueFilter {
	return func(data []byte) bool {
		return bool(C.bridge_queue_filter(wp.P.queue_filter, C.CString(string(data))))
	}
}

func (wp WrapProtocol) QueueMaxSize() int {
	return int(wp.P.queue_max_size)
}

func (wp WrapProtocol) QueueTTL() time.Duration {
	return time.Duration(int(wp.P.queue_ttl))
}

func (wp WrapProtocol) ReplayInterval() time.Duration {
	return time.Duration(int(wp.P.replay_interval))
}

func (wp WrapProtocol) ReconnectionDelay() time.Duration {
	return time.Duration(int(wp.P.reconnection_delay))
}

func (wp WrapProtocol) SslConnection() bool {
	return bool(wp.P.ssl_connection)
}

func (wp WrapProtocol) SetAutoQueue(value bool) {
	wp.P.auto_queue = C.bool(value)
}

func (wp WrapProtocol) SetAutoReplay(value bool) {
	wp.P.auto_replay = C.bool(value)
}

func (wp WrapProtocol) SetOfflineQueueLoader(value protocol.OfflineQueueLoader) {
	//@todo
}

func (wp WrapProtocol) SetQueueFilter(value protocol.QueueFilter) {
	wp.P.queue_filter = C.kuzzle_queue_filter(unsafe.Pointer(&value))
}

func (wp WrapProtocol) SetQueueMaxSize(value int) {
	wp.P.queue_max_size = C.ulonglong(value)
}

func (wp WrapProtocol) SetQueueTTL(value time.Duration) {
	wp.P.queue_ttl = C.ulonglong(value)
}

func (wp WrapProtocol) SetReplayInterval(value time.Duration) {
	wp.P.replay_interval = C.ulonglong(value)
}
