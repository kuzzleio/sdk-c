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
	#include "internal/kuzzle_structs.h"
*/
import "C"
import (
	"github.com/kuzzleio/sdk-go/types"
	"net/http"
	"time"
	"unsafe"
)

//export kuzzle_set_default_query_options
func kuzzle_set_default_query_options(copts *C.query_options) {
	opts := types.NewQueryOptions()

	copts.queuable = C.bool(opts.Queuable())
	copts.withdist = C.bool(opts.Withdist())
	copts.withcoord = C.bool(opts.Withcoord())
	copts.from = C.long(opts.From())
	copts.size = C.long(opts.Size())
	copts.scroll = C.CString(opts.Scroll())
	copts.scroll_id = C.CString(opts.ScrollId())
	copts.refresh = C.CString(opts.Refresh())
	copts.if_exist = C.CString(opts.IfExist())
	copts.retry_on_conflict = C.int(opts.RetryOnConflict())
	copts.volatiles = C.CString(string(opts.Volatile()))
}

//export kuzzle_set_default_room_options
func kuzzle_set_default_room_options(copts *C.room_options) {
	opts := types.NewRoomOptions()

	copts.scope = C.CString(opts.Scope())
	copts.state = C.CString(opts.State())
	copts.users = C.CString(opts.Users())
	copts.volatiles = C.CString(string(opts.Volatile()))
	copts.subscribe_to_self = C.bool(opts.SubscribeToSelf())
}

//export kuzzle_set_default_options
func kuzzle_set_default_options(copts *C.options) {
	opts := types.NewOptions()

	copts.port = C.uint(opts.Port())
	copts.queue_ttl = C.uint(opts.QueueTTL())
	copts.queue_max_size = C.ulong(opts.QueueMaxSize())
	copts.auto_queue = C.bool(opts.AutoQueue())
	copts.auto_reconnect = C.bool(opts.AutoReconnect())
	copts.auto_replay = C.bool(opts.AutoReplay())
	copts.auto_resubscribe = C.bool(opts.AutoResubscribe())
	copts.reconnection_delay = C.ulong(opts.ReconnectionDelay())
	copts.replay_interval = C.ulong(opts.ReplayInterval())
	copts.header_size = C.size_t(0)
	copts.header_names = nil
	copts.header_values = nil

	refresh := opts.Refresh()
	if len(refresh) > 0 {
		copts.refresh = C.CString(refresh)
	} else {
		copts.refresh = nil
	}
}

func SetQueryOptions(options *C.query_options) (opts types.QueryOptions) {
	if options == nil {
		return
	}

	opts = types.NewQueryOptions()

	opts.SetQueuable(bool(options.queuable))
	opts.SetWithdist(bool(options.withdist))
	opts.SetWithcoord(bool(options.withcoord))
	opts.SetFrom(int(options.from))
	opts.SetSize(int(options.size))
	opts.SetScroll(C.GoString(options.scroll))
	opts.SetScrollId(C.GoString(options.scroll_id))
	opts.SetRefresh(C.GoString(options.refresh))
	opts.SetIfExist(C.GoString(options.if_exist))
	opts.SetRetryOnConflict(int(options.retry_on_conflict))
	volatiles := types.VolatileData(C.GoString(options.volatiles))
	if volatiles != nil {
		opts.SetVolatile(volatiles)
	}

	return
}

func SetOptions(options *C.options) (opts types.Options) {
	if options == nil {
		return
	}

	opts = types.NewOptions()

	opts.SetPort(int(options.port))
	opts.SetQueueTTL(time.Duration(uint16(options.queue_ttl)))
	opts.SetQueueMaxSize(int(options.queue_max_size))

	opts.SetAutoQueue(bool(options.auto_queue))
	opts.SetAutoReconnect(bool(options.auto_reconnect))
	opts.SetAutoReplay(bool(options.auto_replay))
	opts.SetAutoResubscribe(bool(options.auto_resubscribe))
	opts.SetReconnectionDelay(time.Duration(int(options.reconnection_delay)))
	opts.SetReplayInterval(time.Duration(int(options.replay_interval)))
	if options.refresh != nil {
		opts.SetRefresh(C.GoString(options.refresh))
	}

	if options.header_size > 0 {
		httpHeaders := &http.Header{}

		hnames := (*[1<<28 - 1]*C.char)(unsafe.Pointer(options.header_names))[:options.header_size:options.header_size]
		hvals := (*[1<<28 - 1]*C.char)(unsafe.Pointer(options.header_values))[:options.header_size:options.header_size]

		for i := 0; i < int(options.header_size); i++ {
			httpHeaders.Add(
				C.GoString(hnames[i]),
				C.GoString(hvals[i]))
		}

		opts.SetHeaders(httpHeaders)
	}

	return
}

func SetRoomOptions(options *C.room_options) (opts types.RoomOptions) {
	opts = types.NewRoomOptions()

	if options != nil {
		opts.SetScope(C.GoString(options.scope))
		opts.SetState(C.GoString(options.state))
		opts.SetUsers(C.GoString(options.users))

		opts.SetSubscribeToSelf(bool(options.subscribe_to_self))

		if options.volatiles != nil {
			opts.SetVolatile(types.VolatileData(C.GoString(options.volatiles)))
		}
	}
	return
}
