-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS meteo_station.homebridge_accessory_mapping
(
    id                  UUID DEFAULT gen_random_uuid() NOT NULL,
    device_id           varchar                        NOT NULL,
    metric_type         varchar                        NOT NULL,
    accessory_unique_id varchar                        NOT NULL,
    characteristic_type varchar                        NOT NULL,
    CONSTRAINT homebridge_accessory_mapping_pk PRIMARY KEY (id),
    CONSTRAINT homebridge_accessory_mapping_unique UNIQUE (device_id, metric_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meteo_station.homebridge_accessory_mapping;
-- +goose StatementEnd