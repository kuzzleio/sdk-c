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

#ifndef __SDK_WRAPPERS_INTERNAL
#define __SDK_WRAPPERS_INTERNAL

#include "internal/kuzzle_structs.h"

#ifdef __cplusplus
namespace kuzzleio {
#endif

typedef char *char_ptr;
typedef policy *policy_ptr;
typedef policy_restriction *policy_restriction_ptr;
typedef query_object *query_object_ptr;

// used by memory_storage.geopos
typedef double geopos_arr[2];

void set_errno(int err);
void kuzzle_notify(kuzzle_notification_listener f, notification_result* res,
                   void* data);
void kuzzle_trigger_event(int event, kuzzle_event_listener f, char* res,
                          void* data);
void room_on_subscribe(kuzzle_subscribe_listener f, room_result* res,
                       void* data);
bool kuzzle_filter_query(kuzzle_queue_filter f, const char *rq);
void free_char_array(char **arr, size_t length);
void assign_geopos(double (*ptr)[2], int idx, double lon, double lat);

#ifdef __cplusplus // end of namespace kuzzleio
}
#endif

#endif
