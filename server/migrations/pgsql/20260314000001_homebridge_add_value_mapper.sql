-- +goose Up
-- +goose StatementBegin
ALTER TABLE meteo_station.homebridge_accessory_mapping ADD COLUMN IF NOT EXISTS value_mapper varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meteo_station.homebridge_accessory_mapping DROP COLUMN IF EXISTS value_mapper;
-- +goose StatementEnd