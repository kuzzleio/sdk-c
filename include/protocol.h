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

  void (*add_listener)(int, kuzzle_event_listener*);
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