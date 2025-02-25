CREATE TABLE devices (
    device_id BIGSERIAL PRIMARY KEY,
    aggregate_identifier UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP(3)
);