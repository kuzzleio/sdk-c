#ifndef KUZZLE_PROTOCOL_H
#define KUZZLE_PROTOCOL_H

#include "kuzzlesdk.h"
#include "sdk_wrappers_internal.h"

typedef struct {
  int auto_queue;
  int auto_reconnect;
  int auto_resubscribe;
  int auto_replay;
  const char* host;
  offline_queue* kuzzle_offline_queue;
  kuzzle_offline_queue_loader offline_queue_loader;
  int port;
  kuzzle_queue_filter queue_filter;
  int queue_max_size;
  long long queue_ttl;
  long long replay_interval;
  long long reconnection_delay;
  int ssl_connection;


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