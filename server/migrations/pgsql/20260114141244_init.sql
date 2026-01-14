-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS meteo_station;
CREATE TABLE IF NOT EXISTS meteo_station.meteo_data
(
    id          UUID DEFAULT gen_random_uuid() not null,
    metric_type varchar                        not null,
    datetime    timestamptz                    not null,
    device_id   varchar                        not null,
    value       decimal                        not null,
    CONSTRAINT meteo_data_pk PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meteo_station.meteo_data;
DROP SCHEMA IF EXISTS meteo_station;
-- +goose StatementEnd
