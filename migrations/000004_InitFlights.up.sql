CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    aircraft_id INTEGER NOT NULL,
    from_location VARCHAR(100) NOT NULL,
    to_location VARCHAR(100) NOT NULL,
    departure_at TIMESTAMP NOT NULL,
    arrive_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_flights_aircraft
        FOREIGN KEY (aircraft_id)
        REFERENCES aircrafts(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_flights_aircraft_id
    ON flights (aircraft_id);

CREATE INDEX idx_flights_departure_at
    ON flights (departure_at);

CREATE INDEX idx_flights_arrive_at
    ON flights (arrive_at);