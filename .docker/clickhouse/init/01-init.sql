CREATE DATABASE IF NOT EXISTS freterapido;

USE freterapido;

CREATE TABLE quotes (
    id UUID DEFAULT generateUUIDv7(),
    name String,
    service String,
    deadline UInt8,
    price Decimal(10, 2),
    timestamp DateTime DEFAULT NOW()
) ENGINE = MergeTree()
ORDER BY
    (service, timestamp);
