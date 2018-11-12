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
	#include "sdk_wrappers_internal.h"

	static void bridge_trigger_kuzzle_notification_result(kuzzle_notification_listener f, notification_result* res, void* data) {
    f(res, data);
	}

	static void bridge_trigger_kuzzle_response(kuzzle_response_listener f, kuzzle_response* res, void* data) {
    f(res, data);
	}
*/
import "C"
import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/kuzzleio/sdk-go/collection"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/security"
	"github.com/kuzzleio/sdk-go/types"
)

// convert a C char** to a go array of string
func cToGoStrings(arr **C.char, len C.size_t) []string {
	if len == 0 {
		return nil
	}

	tmpslice := (*[1 << 27]*C.char)(unsafe.Pointer(arr))[:len:len]
	goStrings := make([]string, 0, len)

	for _, s := range tmpslice {
		goStrings = append(goStrings, C.GoString(s))
	}

	return goStrings
}

func cToGoShards(cShards *C.shards) *types.Shards {
	return &types.Shards{
		Total:      int(cShards.total),
		Successful: int(cShards.successful),
		Failed:     int(cShards.failed),
	}
}

func cToGoKuzzleMeta(cMeta *C.meta) *types.Meta {
	return &types.Meta{
		Author:    C.GoString(cMeta.author),
		CreatedAt: int(cMeta.created_at),
		UpdatedAt: int(cMeta.updated_at),
		Updater:   C.GoString(cMeta.updater),
		Active:    bool(cMeta.active),
		DeletedAt: int(cMeta.deleted_at),
	}
}

func cToGoCollection(c *C.collection) *collection.Collection {
	return collection.NewCollection((*kuzzle.Kuzzle)(c.k.instance))
}

func cToGoPolicyRestriction(r *C.policy_restriction) *types.PolicyRestriction {
	restriction := &types.PolicyRestriction{
		Index: C.GoString(r.index),
	}

	if r.collections == nil {
		return restriction
	}

	slice := (*[1<<28 - 1]*C.char)(unsafe.Pointer(r.collections))[:r.collections_length:r.collections_length]
	for _, col := range slice {
		restriction.Collections = append(restriction.Collections, C.GoString(col))
	}

	return restriction
}

func cToGoPolicy(p *C.policy) *types.Policy {
	policy := &types.Policy{
		RoleId: C.GoString(p.role_id),
	}

	if p.restricted_to == nil {
		return policy
	}

	slice := (*[1<<28 - 1]*C.policy_restriction)(unsafe.Pointer(p.restricted_to))[:p.restricted_to_length]
	for _, crestriction := range slice {
		policy.RestrictedTo = append(policy.RestrictedTo, cToGoPolicyRestriction(crestriction))
	}

	return policy
}

func cToGoProfile(p *C.profile) *security.Profile {
	profile := &security.Profile{
		Id: C.GoString(p.id),
	}

	if p.policies == nil {
		return profile
	}

	slice := (*[1<<28 - 1]*C.policy)(unsafe.Pointer(p.policies))[:p.policies_length]
	for _, cpolicy := range slice {
		profile.Policies = append(profile.Policies, cToGoPolicy(cpolicy))
	}

	return profile
}

func cToGoUser(u *C.kuzzle_user) *security.User {
	if u == nil {
		return nil
	}

	user := security.NewUser(C.GoString(u.id), nil)

	return user
}

func cToGoUserRigh(r *C.user_right) *types.UserRights {
	right := &types.UserRights{
		Controller: C.GoString(r.controller),
		Action:     C.GoString(r.action),
		Index:      C.GoString(r.index),
		Value:      C.GoString(r.value),
	}

	return right
}

func cToGoKuzzleRequest(r *C.kuzzle_request) *types.KuzzleRequest {
	request := &types.KuzzleRequest{
		RequestId:    C.GoString(r.request_id),
		Controller:   C.GoString(r.controller),
		Action:       C.GoString(r.action),
		Index:        C.GoString(r.index),
		Collection:   C.GoString(r.collection),
		Id:           C.GoString(r.id),
		From:         int(r.from),
		Size:         int(r.size),
		Scroll:       C.GoString(r.scroll),
		ScrollId:     C.GoString(r.scroll_id),
		Strategy:     C.GoString(r.strategy),
		ExpiresIn:    int(r.expires_in),
		Scope:        C.GoString(r.scope),
		Users:        C.GoString(r.users),
		Start:        int(r.start),
		Stop:         int(r.stop),
		End:          int(r.end),
		Bit:          int(r.bit),
		Member:       C.GoString(r.member),
		Member1:      C.GoString(r.member1),
		Member2:      C.GoString(r.member2),
		Lon:          float64(r.lon),
		Lat:          float64(r.lat),
		Distance:     float64(r.distance),
		Unit:         C.GoString(r.unit),
		Cursor:       int(r.cursor),
		Offset:       int(r.offset),
		Field:        C.GoString(r.field),
		Subcommand:   C.GoString(r.subcommand),
		Pattern:      C.GoString(r.pattern),
		Idx:          int(r.idx),
		Min:          C.GoString(r.min),
		Max:          C.GoString(r.max),
		Limit:        C.GoString(r.limit),
		Count:        int(r.count),
		Match:        C.GoString(r.match),
		Reset:        bool(r.reset),
		IncludeTrash: bool(r.include_trash),
	}

	if r.body != nil {
		request.Body = json.RawMessage(C.GoString(r.body))
	}
	if r.volatiles != nil {
		request.Volatile = json.RawMessage(C.GoString(r.volatiles))
	}

	request.Members = cToGoStrings(r.members, r.members_length)
	request.Keys = cToGoStrings(r.keys, r.keys_length)
	request.Fields = cToGoStrings(r.fields, r.fields_length)
	options := cToGoStrings(r.options, r.options_length)
	for _, option := range options {
		request.Options = append(request.Options, option)
	}

	return request
}

func cToGoKuzzleResponse(r *C.kuzzle_response) *types.KuzzleResponse {
	response := &types.KuzzleResponse{
		RequestId:  C.GoString(r.request_id),
		Index:      C.GoString(r.index),
		Collection: C.GoString(r.collection),
		Controller: C.GoString(r.controller),
		Action:     C.GoString(r.action),
		RoomId:     C.GoString(r.room_id),
		Channel:    C.GoString(r.channel),
		Status:     int(r.status),
	}

	if r.error != nil {
		response.Error = types.NewError(C.GoString(r.error), int(r.status))
		response.Error.Stack = C.GoString(r.stack)
	}

	response.Result = json.RawMessage(C.GoString(r.result))
	response.Volatile, _ = json.Marshal(r.volatiles)

	return response
}

func cToGoSearchResult(s *C.search_result) *types.SearchResult {
	options := types.NewQueryOptions()

	options.SetSize(int(s.options.size))
	options.SetFrom(int(s.options.from))
	options.SetScroll(C.GoString(s.options.scroll))
	options.SetScrollId(C.GoString(s.options.scroll_id))

	kuzzle := (*kuzzle.Kuzzle)(s.k.instance)
	scrollAction := C.GoString(s.scroll_action)
	request := cToGoKuzzleRequest(s.request)
	response := cToGoKuzzleResponse(s.response)

	sr, err := types.NewSearchResult(kuzzle, scrollAction, request, options, response)
	if err != nil {
		panic(err)
	}

	return sr
}

func cToGoKuzzleNotificationChannel(c *C.kuzzle_notification_listener) chan<- types.KuzzleNotification {
	return make(chan<- types.KuzzleNotification)
}

func cToGoQueryObject(cqo *C.query_object, data unsafe.Pointer) *types.QueryObject {
	notifChan := make(chan *types.KuzzleNotification)
	go func() {
		for {
			res, ok := <-notifChan
			if ok == false {
				break
			}

			C.bridge_trigger_kuzzle_notification_result(cqo.listener, goToCNotificationResult(res), data)
		}
	}()

	resChan := make(chan *types.KuzzleResponse)
	go func() {
		for {
			res, ok := <-resChan
			if ok == false {
				break
			}

			C.bridge_trigger_kuzzle_response(cqo.response_listener, goToCKuzzleResponse(res), data)
		}
	}()

	goqo := &types.QueryObject{
		Query:     []byte(C.GoString(cqo.query)),
		Options:   SetQueryOptions(&cqo.options),
		NotifChan: notifChan,
		ResChan:   resChan,
		Timestamp: time.Unix(int64(cqo.timestamp), 0),
		RequestId: C.GoString(cqo.request_id),
	}

	return goqo
}

