#include <stdlib.h>
#include "internal/protocol.h"
#include "_cgo_export.h"

void call_bridge(int event, char* res, void* data) {
  bridge_listener(event, res, data);
}

void call_bridge_once(int event, char* res, void* data) {
  bridge_listener_once(event, res, data);
}

void call_notification_bridge(notification_result* result, void* data) {
  bridge_notification(result, data);
}

void bridge_protocol_add_listener(void (*f)(int, kuzzle_event_listener, void*),
                                  int event, kuzzle_event_listener listener,
                                  void* data) {
  f(event, listener, data);
}

void bridge_protocol_once(void (*f)(int, kuzzle_event_listener, void*),
                          int event, kuzzle_event_listener listener,
                          void* data) {
  f(event, listener, data);
}

void bridge_remove_listener(void (*f)(int, kuzzle_event_listener, void*),
                            int event, kuzzle_event_listener listener,
                            void* data) {
  f(event, listener, data);
}

void bridge_remove_all_listeners(void (*f)(int, void*), int event, void* data) {
  f(event, data);
}

void bridge_once(void (*f)(int, kuzzle_event_listener, void*), int event,
                 kuzzle_event_listener listener, void* data) {
  f(event, listener, data);
}

int bridge_listener_count(int (*f)(int, void*), int event, void* data) {
  return f(event, data);
}

char* bridge_connect(char* (*f)(void*), void* data) {
  return f(data);
}

kuzzle_response* bridge_send(
    kuzzle_response* (*f)(const char*, query_options*, char*, void*),
    const char* query,
    query_options* options,
    char* request_id,
    void* data) {
  return f(query, options, request_id, data);
}

char* bridge_close(char* (*f)(void*), void* data) {
  return f(data);
}

int bridge_get_state(int (*f)(void*), void* data) {
  return f(data);
}

void bridge_emit_event(void (*f)(int, void*, void*), int event,
                       void* res, void* data) {
  f(event, res, data);
}

void bridge_register_sub(
    void (*f)(const char*, const char*, const char*, bool,
              kuzzle_notification_listener, void*),
    const char* channel,
    const char* room_id,
    const char* filters,
    bool subscribe_to_self,
    void* data) {
  f(channel, room_id, filters, subscribe_to_self, call_notification_bridge,
    data);
}

void bridge_unregister_sub(void (*f)(char*, void*), char* id, void* data) {
  f(id, data);
}

void bridge_cancel_subs(void (*f)(void*), void* data) {
  f(data);
}

void bridge_start_queuing(void (*f)(void*), void* data) {
  f(data);
}

void bridge_stop_queuing(void (*f)(void*), void* data) {
  f(data);
}

void bridge_play_queue(void (*f)(void*), void* data) {
  f(data);
}

void bridge_clear_queue(void (*f)(void*), void* data) {
  f(data);
}

bool bridge_queue_filter(kuzzle_queue_filter f, const char* data) {
  return f(data);
}

bool bridge_is_auto_queue(bool (*f)(void*), void* data) {
  return f(data);
}

bool bridge_is_auto_reconnect(bool (*f)(void*), void* data) {
  return f(data);
}

bool bridge_is_auto_resubscribe(bool (*f)(void*), void* data) {
  return f(data);
}

const char* bridge_get_host(const char* (*f)(void*), void* data) {
  return f(data);
}

unsigned int bridge_get_port(unsigned int (*f)(void*), void* data) {
  return f(data);
}

unsigned long long bridge_get_reconnection_delay(unsigned long long (*f)(void*),
                                                 void* data) {
  return f(data);
}

bool bridge_is_ssl_connection(bool (*f)(void*), void* data) {
  return f(data);
}

kuzzle_event_listener get_bridge_fptr() {
  return &call_bridge;
}

kuzzle_event_listener get_bridge_once_fptr() {
  return &call_bridge_once;
}

kuzzle_notification_listener get_bridge_notification_listener_fptr() {
  return &call_notification_bridge;
}

bool bridge_is_ready(bool (*f)(void*), void* data) {
  return f(data);
}

void bridge_trigger_event_listener(kuzzle_event_listener listener, int event,
                                   char* res, void* data) {
  listener(event, res, data);
}

void bridge_trigger_notification_listener(kuzzle_notification_listener listener,
                                          notification_result* result,
                                          void* data) {
  listener(result, data);
}

void bridge_trigger_kuzzle_notification_result(kuzzle_notification_listener f,
                                               notification_result* res,
                                               void* data) {
  f(res, data);
}

void bridge_trigger_kuzzle_response(kuzzle_response_listener f,
                                    kuzzle_response* res, void* data) {
  f(res, data);
}
