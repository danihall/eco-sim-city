CREATE OR REPLACE FUNCTION mutable.func_increment_game_tick ()
    RETURNS TRIGGER AS $$
BEGIN
    SELECT COALESCE(MAX(game_tick), 0) + 1
    INTO NEW.game_tick
    FROM mutable.cities_state_history
    WHERE city_id = NEW.city_id;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS mutable.cities_state_history (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    game_tick BIGINT CHECK (game_tick >= 0),
    recorded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    electricity_consumption INT NOT NULL DEFAULT 0,
    water_consumption INT NOT NULL DEFAULT 0,
    children_education_level SMALLINT NOT NULL DEFAULT 0,
    adults_education_level SMALLINT NOT NULL DEFAULT 0,
    elderly_education_level SMALLINT NOT NULL DEFAULT 0,
    children_health_level SMALLINT NOT NULL DEFAULT 0,
    adults_health_level SMALLINT NOT NULL DEFAULT 0,
    elderly_health_level SMALLINT NOT NULL DEFAULT 0,
    children_count INT NOT NULL CHECK (children_count >= 0) DEFAULT 0,
    adults_count INT NOT NULL CHECK (adults_count >= 0) DEFAULT 0,
    elderly_count INT NOT NULL CHECK (elderly_count >= 0) DEFAULT 0,
    -- global direct pollution
    ghg_1_level INT NOT NULL,
    -- global indirect pollution from power generation
    ghg_2_level INT NOT NULL DEFAULT 0,
    -- global garbage
    ghg_3_level INT NOT NULL DEFAULT 0,
    citizens_education_level INT GENERATED ALWAYS AS (children_education_level + adults_education_level + elderly_education_level) STORED,
    citizens_health_level INT GENERATED ALWAYS AS (children_health_level + adults_health_level + elderly_health_level) STORED,
    citizens_count INT GENERATED ALWAYS AS (children_count + adults_count + elderly_count) STORED,
    climate_change_level INT GENERATED ALWAYS AS (GREATEST (0, ghg_1_level + ghg_2_level + ghg_3_level)) STORED,
    PRIMARY KEY (city_id, game_tick)
);

CREATE TRIGGER trigger_cities_state_history_game_tick
    BEFORE INSERT ON mutable.cities_state_history
    FOR EACH ROW
    EXECUTE FUNCTION mutable.func_increment_game_tick ();

---- create above / drop below ----
DROP TRIGGER IF EXISTS trigger_cities_state_history_game_tick ON mutable.cities_state_history;
DROP TABLE IF EXISTS mutable.cities_state_history;
DROP FUNCTION IF EXISTS mutable.func_increment_game_tick;
