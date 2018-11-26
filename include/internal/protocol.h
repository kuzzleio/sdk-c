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

#ifndef KUZZLE_PROTOCOL_H
#define KUZZLE_PROTOCOL_H

#include "kuzzle_structs.h"
#include "sdk_wrappers_internal.h"

typedef struct {
  void* instance;

  void (*add_listener)(int, kuzzle_event_listener*, void*);
  void (*remove_listener)(int, kuzzle_event_listener*, void*);
  void (*remove_all_listeners)(int, void*);
  void (*once)(int, kuzzle_event_listener*, void*);
  int (*listener_count)(int, void*);
  char* (*connect)(void*);
  kuzzle_response* (*send)(const char*, query_options*, char*, void*);
  const char* (*close)(void*);
  int (*get_state)(void*);
  void (*emit_event)(int, void*, void*);
  void (*register_sub)(const char*, const char*, const char*, bool, kuzzle_notification_listener*, void*);
  void (*unregister_sub)(const char*, void*);
  void (*cancel_subs)(void*);
  void (*start_queuing)(void*);
  void (*stop_queuing)(void*);
  void (*play_queue)(void*);
  void (*clear_queue)(void*);

  bool (*is_auto_queue)(void*);
  bool (*is_auto_reconnect)(void*);
  bool (*is_auto_resubscribe)(void*);
  const char* (*get_host)(void*);
  unsigned int (*get_port)(void*);
  unsigned long long (*get_reconnection_delay)(void*);
  bool (*is_ssl_connection)(void*);

} protocol;

#endif