-- +goose Up
-- +goose StatementBegin
ALTER TABLE meteo_station.homebridge_accessory_mapping DROP COLUMN IF EXISTS characteristic_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meteo_station.homebridge_accessory_mapping ADD COLUMN characteristic_type varchar NOT NULL DEFAULT '';
-- +goose StatementEnd