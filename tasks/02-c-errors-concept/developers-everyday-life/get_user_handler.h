#ifndef GET_USER_HANDLER_H
#define GET_USER_HANDLER_H

#include <stdlib.h>
#include "db.h"
#include "marshalers.h"

typedef enum {
    HTTP_ERR_OK = 0,                   // 200
    HTTP_ERR_BAD_REQUEST = 1,          // 400
    HTTP_ERR_UNPROCESSABLE_ENTITY = 2, // 422
    HTTP_ERR_NOT_FOUND = 3,            // 404
    HTTP_ERR_INTERNAL_SERVER = 4,      // 500

} http_error_t;

const char* const HTTP_ERR_STRS[] = {
    "200 OK",
    "400 Bad Request",
    "422 Unprocessable Entity",
    "404 Not Found",
    "500 Internal Server Error",
};

const char *http_error_str(http_error_t err)
{
    return HTTP_ERR_STRS[err];
}

http_error_t get_user_handler(char *request_data, char **response_data)
{   
    request_t *req = NULL;
    int unmarshal_err = unmarshal_request(request_data, &req);
    if (unmarshal_err == -1) {
        return HTTP_ERR_BAD_REQUEST;
    }

    if (req->user_id <= 0) {
        free(req);

        return HTTP_ERR_UNPROCESSABLE_ENTITY;
    }

    user_t *user = NULL;
    db_error_t db_err = db_get_user_by_id(req->user_id, &user);
    if (db_err == DB_ERR_NOT_FOUND) {
        free(req);
        
        return HTTP_ERR_NOT_FOUND;

    } else if (db_err == DB_ERR_INTERNAL) {
        free(req);

        return HTTP_ERR_INTERNAL_SERVER;
    }

    int marshal_err = marshal_response((response_t){user}, response_data);
    if (marshal_err == -1) {
        free(req);
        free(user->email);
        free(user);

        return HTTP_ERR_INTERNAL_SERVER;
    }

    free(req);
    free(user->email);
    free(user);

    return HTTP_ERR_OK;
}

#endif
