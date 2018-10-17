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

#include "kuzzlesdk.h"
#include "sdk_wrappers_internal.h"

typedef struct {
  bool auto_queue;
  bool auto_reconnect;
  bool auto_resubscribe;
  bool auto_replay;
  const char* host;
  offline_queue* kuzzle_offline_queue;
  kuzzle_offline_queue_loader offline_queue_loader;
  int port;
  kuzzle_queue_filter queue_filter;
  unsigned long long queue_max_size;
  unsigned long long queue_ttl;
  unsigned long long replay_interval;
  unsigned long long reconnection_delay;
  bool ssl_connection;

  void (*add_listener)(int, kuzzle_event_listener*, void*);
  void (*remove_listener)(int, kuzzle_event_listener*);
  void (*remove_all_listeners)(int);
  void (*once)(int, kuzzle_event_listener*);
  int (*listener_count)(int);
  char* (*connect)();
  char* (*send)(const char*, query_options*, kuzzle_response*, char*);
  char* (*close)();
  int (*get_state)();
  void (*emit_event)(int, void*);
  void (*register_sub)(const char*, const char*, const char*, int, kuzzle_notification_listener*, void*);
  void (*unregister_sub)(const char*);
  void (*cancel_subs)();
  void (*start_queuing)();
  void (*stop_queuing)();
  void (*play_queue)();
  void (*clear_queue)();
} protocol;

#endif