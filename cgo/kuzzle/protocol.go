package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "protocol.h"

	// Bridges

	static void bridge_protocol_add_listener(void (*f)(int, kuzzle_event_listener*, void*), int event, kuzzle_event_listener* listener, void* data) {
		f(event, listener, data);
	}

	static void bridge_protocol_once(void (*f)(int, kuzzle_event_listener*, void*), int event, kuzzle_event_listener* listener, void* data) {
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

	static void bridge_register_sub(void (*f)(const char*, const char*, const char*, bool, kuzzle_notification_listener*, void*), const char* channel, const char* room_id, const char* filters, bool subscribe_to_self, kuzzle_notification_listener* listener, void* data) {
		f(channel, room_id, filters, subscribe_to_self, listener, data);
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

	static bool bridge_is_auto_queue(bool (*f)(void*), void* data) {
		return f(data);
	}

	static bool bridge_is_auto_reconnect(bool (*f)(void*), void* data) {
		return f(data);
	}

	static bool bridge_is_auto_resubscribe(bool (*f)(void*), void* data) {
		return f(data);
	}

	static const char* bridge_get_host(const char* (*f)(void*), void* data) {
		return f(data);
	}

	static unsigned int bridge_get_port(unsigned int (*f)(void*), void* data) {
		return f(data);
	}

	static unsigned long long bridge_get_reconnection_delay(unsigned long long (*f)(void*), void* data) {
		return f(data);
	}

	static bool bridge_is_ssl_connection(bool (*f)(void*), void* data) {
		return f(data);
	}

	extern void bridge_listener(int, char*, void*);
	extern void bridge_listener_once(int, char*, void*);
	extern void bridge_notification(notification_result*, void*);

	static void call_bridge(int event, char* res, void* data) {
		bridge_listener(event, res, data);
	}

	static void call_bridge_once(int event, char* res, void* data) {
		bridge_listener_once(event, res, data);
	}

	static void call_notification_bridge(notification_result* result, void* data) {
		bridge_notification(result, data);
	}

	static kuzzle_event_listener get_bridge_fptr() {
		return &call_bridge;
	}

	static kuzzle_event_listener get_bridge_once_fptr() {
		return &call_bridge_once;
	}

	static kuzzle_notification_listener get_bridge_notification_listener_fptr() {
		return &call_notification_bridge;
	}
*/
import "C"
import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"unsafe"

	"github.com/kuzzleio/sdk-go/protocol"
	"github.com/kuzzleio/sdk-go/types"
)

type WrapProtocol struct {
	P *C.protocol
}

var proto_instances sync.Map

var _list_listeners map[int]map[chan<- json.RawMessage]bool
var _list_once_listeners map[int]map[chan<- json.RawMessage]bool
var _list_notification_listeners map[string]chan<- types.NotificationResult

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
	_list_listeners = make(map[int]map[chan<- json.RawMessage]bool)
	_list_once_listeners = make(map[int]map[chan<- json.RawMessage]bool)
	_list_notification_listeners = make(map[string]chan<- types.NotificationResult)

	return &WrapProtocol{p}
}

//export bridge_listener
func bridge_listener(event C.int, res *C.char, data unsafe.Pointer) {
	for c := range _list_listeners[int(event)] {
		c <- json.RawMessage(C.GoString(res))
	}
}

//export bridge_listener_once
func bridge_listener_once(event C.int, res *C.char, data unsafe.Pointer) {
	for c := range _list_once_listeners[int(event)] {
		c <- json.RawMessage(C.GoString(res))
	}
}

func (wp WrapProtocol) AddListener(event int, channel chan<- json.RawMessage) {
	fptr := C.get_bridge_fptr()
	if _list_listeners[event] == nil {
		_list_listeners[event] = make(map[chan<- json.RawMessage]bool)
	}
	_list_listeners[event][channel] = true
	C.bridge_protocol_add_listener(wp.P.add_listener, C.int(event), &fptr, wp.P.instance)
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- json.RawMessage) {
	fptr := C.get_bridge_fptr()
	C.bridge_remove_listener(wp.P.remove_listener, C.int(event), &fptr, wp.P.instance)
	delete(_list_listeners[event], channel)
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	C.bridge_remove_all_listeners(wp.P.remove_all_listeners, C.int(event), wp.P.instance)
}

func (wp WrapProtocol) Once(event int, channel chan<- json.RawMessage) {
	fptr := C.get_bridge_once_fptr()
	if _list_once_listeners[event] == nil {
		_list_once_listeners[event] = make(map[chan<- json.RawMessage]bool)
	}
	_list_once_listeners[event][channel] = true
	C.bridge_protocol_once(wp.P.once, C.int(event), &fptr, wp.P.instance)
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

//export bridge_notification
func bridge_notification(res *C.notification_result, data unsafe.Pointer) {
	_list_notification_listeners[C.GoString(res.room_id)] <- *cToGoNotificationResult(res)
}

func (wp WrapProtocol) RegisterSub(channel, roomID string, filters json.RawMessage, subscribeToSelf bool, notifChan chan<- types.NotificationResult, onReconnectChannel chan<- interface{}) {
	fptr := C.get_bridge_notification_listener_fptr()
	_list_notification_listeners[channel] = notifChan
	C.bridge_register_sub(wp.P.register_sub, C.CString(channel), C.CString(roomID), C.CString(string(filters)), C.bool(subscribeToSelf), &fptr, wp.P.instance)
}

func (wp WrapProtocol) UnregisterSub(id string) {
	C.bridge_unregister_sub(wp.P.unregister_sub, C.CString(id), wp.P.instance)
	delete(_list_notification_listeners, id)
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
	return bool(C.bridge_is_auto_queue(wp.P.is_auto_queue, wp.P.instance))
}

func (wp WrapProtocol) AutoReconnect() bool {
	return bool(C.bridge_is_auto_reconnect(wp.P.is_auto_reconnect, wp.P.instance))
}

func (wp WrapProtocol) AutoResubscribe() bool {
	return bool(C.bridge_is_auto_resubscribe(wp.P.is_auto_resubscribe, wp.P.instance))
}

// Deprecated
func (wp WrapProtocol) AutoReplay() bool {
	//@todo will be moved in Kuzzle object
	return false
}

func (wp WrapProtocol) Host() string {
	return C.GoString(C.bridge_get_host(wp.P.get_host, wp.P.instance))
}

// Deprecated
func (wp WrapProtocol) OfflineQueue() []*types.QueryObject {
	//@todo will be moved in Kuzzle object
	return nil
}

// Deprecated
func (wp WrapProtocol) OfflineQueueLoader() protocol.OfflineQueueLoader {
	//@todo will be moved in Kuzzle object
	var offline protocol.OfflineQueueLoader
	return offline
}

func (wp WrapProtocol) Port() int {
	return int(C.bridge_get_port(wp.P.get_port, wp.P.instance))
}

// Deprecated
func (wp WrapProtocol) QueueFilter() protocol.QueueFilter {
	//@todo will be moved in Kuzzle object
	return func(data []byte) bool {
		return false
	}
}

// Deprecated
func (wp WrapProtocol) QueueMaxSize() int {
	//@todo will be moved in Kuzzle object
	return -1
}

// Deprecated
func (wp WrapProtocol) QueueTTL() time.Duration {
	//@todo will be moved in Kuzzle object
	return time.Duration(0)
}

// Deprecated
func (wp WrapProtocol) ReplayInterval() time.Duration {
	//@todo will be moved in Kuzzle object
	return time.Duration(0)
}

func (wp WrapProtocol) ReconnectionDelay() time.Duration {
	return time.Duration(int(C.bridge_get_reconnection_delay(wp.P.get_reconnection_delay, wp.P.instance)))
}

func (wp WrapProtocol) SslConnection() bool {
	return bool(C.bridge_is_ssl_connection(wp.P.is_ssl_connection, wp.P.instance))
}

// Deprecated
func (wp WrapProtocol) SetAutoQueue(value bool) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetAutoReplay(value bool) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetOfflineQueueLoader(value protocol.OfflineQueueLoader) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetQueueFilter(value protocol.QueueFilter) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetQueueMaxSize(value int) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetQueueTTL(value time.Duration) {
	//@todo will be moved in Kuzzle object
}

// Deprecated
func (wp WrapProtocol) SetReplayInterval(value time.Duration) {
	//@todo will be moved in Kuzzle object
}
