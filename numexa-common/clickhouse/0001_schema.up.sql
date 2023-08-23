CREATE DATABASE IF NOT EXISTS numexa;

CREATE TABLE IF NOT EXISTS proxy_requests (
    request_timestamp DateTime,
    source_ip String,
    request_method String,
    request_url String,
    request_headers Nested (
        name String,
        value String
    ),
    request_body String,
    response_timestamp DateTime,
    response_status_code UInt16,
    response_headers Nested (
        name String,
        value String
    ),
    response_body String,
    provider String
) ENGINE = MergeTree()
ORDER BY request_timestamp;
