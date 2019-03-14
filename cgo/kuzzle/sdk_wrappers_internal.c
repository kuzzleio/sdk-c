#include "internal/sdk_wrappers_internal.h"
#include <errno.h>

void set_errno(int err) {
  errno = err;
}

void kuzzle_notify(kuzzle_notification_listener f, notification_result* res,
                   void* data) {
    f(res, data);
}

void kuzzle_trigger_event(int event, kuzzle_event_listener f, char* res,
                          void* data) {
    f(event, res, data);
}

void room_on_subscribe(kuzzle_subscribe_listener f, room_result* res,
                       void* data) {
    f(res, data);
}

bool kuzzle_filter_query(kuzzle_queue_filter f, const char *rq) {
  return f(rq);
}

void free_char_array(char **arr, size_t length) {
  if (arr != NULL) {
    for(int i = 0; i < length; i++) {
      free(arr[i]);
    }

    free(arr);
  }
}

void assign_geopos(double (*ptr)[2], int idx, double lon, double lat) {
  ptr[idx][0] = lon;
  ptr[idx][1] = lat;
}
