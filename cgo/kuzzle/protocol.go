package main

/*
	#cgo CFLAGS: -I../../include
	#include <stdlib.h>
	#include "internal/protocol.h"
	#include "internal/bridges.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"unsafe"
	"github.com/kuzzleio/sdk-go/types"
)

type WrapProtocol struct {
	P *C.protocol
}

var proto_instances sync.Map

var _list_listeners map[int]map[chan<- json.RawMessage]bool = make(map[int]map[chan<- json.RawMessage]bool)
var _list_once_listeners map[int]map[chan<- json.RawMessage]bool = make(map[int]map[chan<- json.RawMessage]bool)
var _list_notification_listeners map[string]chan<- types.NotificationResult = make(map[string]chan<- types.NotificationResult)

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
	if _list_listeners[event] == nil {
		_list_listeners[event] = make(map[chan<- json.RawMessage]bool)
	}
	_list_listeners[event][channel] = true

	C.bridge_protocol_add_listener(
		wp.P.add_listener,
		C.int(event),
		C.get_bridge_fptr(),
		wp.P.instance)
}

func (wp WrapProtocol) RemoveListener(event int, channel chan<- json.RawMessage) {
	C.bridge_remove_listener(
		wp.P.remove_listener,
		C.int(event),
		C.get_bridge_fptr(),
		wp.P.instance)

	delete(_list_listeners[event], channel)
}

func (wp WrapProtocol) RemoveAllListeners(event int) {
	C.bridge_remove_all_listeners(wp.P.remove_all_listeners, C.int(event), wp.P.instance)
}

func (wp WrapProtocol) Once(event int, channel chan<- json.RawMessage) {
	if _list_once_listeners[event] == nil {
		_list_once_listeners[event] = make(map[chan<- json.RawMessage]bool)
	}
	_list_once_listeners[event][channel] = true

	C.bridge_protocol_once(
		wp.P.once,
		C.int(event),
		C.get_bridge_once_fptr(),
		wp.P.instance)
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
	str, err := json.Marshal(data)

	if err != nil {
		str = []byte("")
	}

	p := unsafe.Pointer(C.CString(string(str)))
	defer C.free(p)
	C.bridge_emit_event(wp.P.emit_event, C.int(event), p, wp.P.instance)
}

//export bridge_notification
func bridge_notification(res *C.notification_result, data unsafe.Pointer) {
	goNotification := cToGoNotificationResult(res)
	_list_notification_listeners[C.GoString(res.room_id)] <- *goNotification
}

func (wp WrapProtocol) RegisterSub(channel, roomID string,
	filters json.RawMessage,
	subscribeToSelf bool,
	notifChan chan<- types.NotificationResult,
	onReconnectChannel chan<- interface{}) {

	_list_notification_listeners[channel] = notifChan

	C.bridge_register_sub(wp.P.register_sub, C.CString(channel),
		C.CString(roomID), C.CString(string(filters)), C.bool(subscribeToSelf),
		wp.P.instance)
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

func (wp WrapProtocol) AutoReconnect() bool {
	return bool(C.bridge_is_auto_reconnect(wp.P.is_auto_reconnect, wp.P.instance))
}

func (wp WrapProtocol) AutoResubscribe() bool {
	return bool(C.bridge_is_auto_resubscribe(wp.P.is_auto_resubscribe, wp.P.instance))
}

func (wp WrapProtocol) Port() int {
	return int(C.bridge_get_port(wp.P.get_port, wp.P.instance))
}

func (wp WrapProtocol) ReconnectionDelay() time.Duration {
	return time.Duration(int(C.bridge_get_reconnection_delay(wp.P.get_reconnection_delay, wp.P.instance)))
}

func (wp WrapProtocol) SslConnection() bool {
	return bool(C.bridge_is_ssl_connection(wp.P.is_ssl_connection, wp.P.instance))
}

func (wp WrapProtocol) IsReady() bool {
	return bool(C.bridge_is_ready(wp.P.is_ready, wp.P.instance))
}

func (wp WrapProtocol) Host() string {
	return C.GoString(C.bridge_get_host(wp.P.get_host, wp.P.instance))
}