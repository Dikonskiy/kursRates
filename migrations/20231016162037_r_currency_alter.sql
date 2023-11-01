-- +goose Up
-- +goose StatementBegin
ALTER TABLE R_CURRENCY
ADD COLUMN U_DATE TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE R_CURRENCY
DROP COLUMN U_DATE;
-- +goose StatementEnd
