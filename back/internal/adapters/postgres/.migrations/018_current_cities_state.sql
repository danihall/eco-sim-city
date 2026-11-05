CREATE OR REPLACE FUNCTION computed.func_update_cities_state_latest()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO computed.current_cities_state (
        city_id,
        game_tick,
        electricity_consumption,
        water_consumption,
        children_education_level,
        adults_education_level,
        elderly_education_level,
        children_health_level,
        adults_health_level,
        elderly_health_level,
        children_count,
        adults_count,
        elderly_count,
        ghg_1_level ,
        ghg_2_level ,
        ghg_3_level
    )
    VALUES (
        NEW.city_id,
        NEW.game_tick,
        NEW.electricity_consumption,
        NEW.water_consumption,
        NEW.children_education_level,
        NEW.adults_education_level,
        NEW.elderly_education_level,
        NEW.children_health_level,
        NEW.adults_health_level,
        NEW.elderly_health_level,
        NEW.children_count,
        NEW.adults_count,
        NEW.elderly_count,
        NEW.ghg_1_level ,
        NEW.ghg_2_level ,
        NEW.ghg_3_level
    )
    ON CONFLICT (city_id) DO UPDATE SET
        game_tick = EXCLUDED.game_tick,
        electricity_consumption = EXCLUDED.electricity_consumption,
        water_consumption = EXCLUDED.water_consumption,
        children_education_level = EXCLUDED.children_education_level,
        adults_education_level = EXCLUDED.adults_education_level,
        elderly_education_level = EXCLUDED.elderly_education_level,
        children_health_level = EXCLUDED.children_health_level,
        adults_health_level = EXCLUDED.adults_health_level,
        elderly_health_level = EXCLUDED.elderly_health_level,
        children_count = EXCLUDED.children_count,
        adults_count = EXCLUDED.adults_count,
        elderly_count = EXCLUDED.elderly_count,
        ghg_1_level = EXCLUDED.ghg_1_level ,
        ghg_2_level = EXCLUDED.ghg_2_level ,
        ghg_3_level = EXCLUDED.ghg_3_level;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE UNLOGGED TABLE IF NOT EXISTS computed.current_cities_state (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    game_tick BIGINT CHECK (game_tick >= 0),
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
    ghg_2_level INT NOT NULL,
    -- global garbage
    ghg_3_level INT NOT NULL,
    citizens_education_level INT GENERATED ALWAYS AS (children_education_level + adults_education_level + elderly_education_level) STORED,
    citizens_health_level INT GENERATED ALWAYS AS (children_health_level + adults_health_level + elderly_health_level) STORED,
    citizens_count INT GENERATED ALWAYS AS (children_count + adults_count + elderly_count) STORED,
    climate_change_level INT GENERATED ALWAYS AS (GREATEST (0, ghg_1_level + ghg_2_level + ghg_3_level)) STORED,
    PRIMARY KEY (city_id)
);

CREATE TRIGGER trigger_update_cities_state_latest
    AFTER INSERT ON mutable.cities_state_history
    FOR EACH ROW
    EXECUTE FUNCTION computed.func_update_cities_state_latest();

---- create above / drop below ----
DROP TRIGGER IF EXISTS trigger_update_cities_state_latest ON mutable.cities_state_history;
DROP TABLE IF EXISTS computed.current_cities_state;
DROP FUNCTION IF EXISTS computed.func_update_cities_state_latest;
