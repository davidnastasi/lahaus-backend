CREATE TYPE property_type as enum ('HOUSE','APARTMENT');
CREATE TYPE property_status as enum ('ACTIVE','INACTIVE', 'INVALID');

CREATE TABLE properties (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    sale_price INTEGER NOT NULL,
    administrative_fee INTEGER,
    property_type property_type NOT NULL,
    bedrooms INTEGER NOT NULL,
    bathrooms INTEGER NOT NULL,
    parking_spots INTEGER  NULL,
    area INTEGER NOT NULL,
    photos TEXT[] NULL,
    status property_status NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

-- Add procedure and triggers for field updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_update_at_timestamp
    BEFORE UPDATE ON properties
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX properties_update_at_idx ON properties (updated_at);


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email CHARACTER VARYING(320) not null,
    password CHARACTER VARYING(256) not null
);

CREATE UNIQUE INDEX users_email_idx ON users (email);

CREATE TABLE favourites (
   user_id BIGINT NOT NULL,
   property_id BIGINT NOT NULL
);

CREATE UNIQUE INDEX favourites_user_property_idx ON favourites (user_id,property_id);

-- ADD constraints
ALTER TABLE favourites
    ADD CONSTRAINT fk_favorites_users
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE favourites
    ADD CONSTRAINT fk_favorites_properties
        FOREIGN KEY (property_id)
            REFERENCES properties (id);

