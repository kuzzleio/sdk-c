#ifndef _H_KUZZLE_CGO_BRIDGES
#define _H_KUZZLE_CGO_BRIDGES

void call_bridge(int event, char* res, void* data);
void call_bridge_once(int event, char* res, void* data);
void call_notification_bridge(notification_result* result, void* data);
void bridge_protocol_add_listener(void (*f)(int, kuzzle_event_listener, void*),
                                  int event, kuzzle_event_listener listener,
                                  void* data);
void bridge_protocol_once(void (*f)(int, kuzzle_event_listener, void*),
                          int event, kuzzle_event_listener listener,
                          void* data);
void bridge_remove_listener(void (*f)(int, kuzzle_event_listener, void*),
                            int event, kuzzle_event_listener listener,
                            void* data);
void bridge_remove_all_listeners(void (*f)(int, void*), int event, void* data);
void bridge_once(void (*f)(int, kuzzle_event_listener, void*), int event,
                 kuzzle_event_listener listener, void* data);
int bridge_listener_count(int (*f)(int, void*), int event, void* data);
char* bridge_connect(char* (*f)(void*), void* data);
kuzzle_response* bridge_send(
    kuzzle_response* (*f)(const char*, query_options*, char*, void*),
    const char* query,
    query_options* options,
    char* request_id,
    void* data);
char* bridge_close(char* (*f)(void*), void* data);
int bridge_get_state(int (*f)(void*), void* data);
void bridge_emit_event(void (*f)(int, void*, void*), int event,
                       void* res, void* data);
void bridge_register_sub(
    void (*f)(const char*, const char*, const char*, bool,
              kuzzle_notification_listener, void*),
    const char* channel,
    const char* room_id,
    const char* filters,
    bool subscribe_to_self,
    void* data);
void bridge_unregister_sub(void (*f)(char*, void*), char* id, void* data);
void bridge_cancel_subs(void (*f)(void*), void* data);
void bridge_start_queuing(void (*f)(void*), void* data);
void bridge_stop_queuing(void (*f)(void*), void* data);
void bridge_play_queue(void (*f)(void*), void* data);
void bridge_clear_queue(void (*f)(void*), void* data);
bool bridge_queue_filter(kuzzle_queue_filter f, const char* data);
bool bridge_is_auto_queue(bool (*f)(void*), void* data);
bool bridge_is_auto_reconnect(bool (*f)(void*), void* data);
bool bridge_is_auto_resubscribe(bool (*f)(void*), void* data);
const char* bridge_get_host(const char* (*f)(void*), void* data);
unsigned int bridge_get_port(unsigned int (*f)(void*), void* data);
unsigned long long bridge_get_reconnection_delay(unsigned long long (*f)(void*),
                                                 void* data);
bool bridge_is_ssl_connection(bool (*f)(void*), void* data);
kuzzle_event_listener get_bridge_fptr();
kuzzle_event_listener get_bridge_once_fptr();
kuzzle_notification_listener get_bridge_notification_listener_fptr();

#endif
