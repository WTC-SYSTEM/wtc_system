CREATE TABLE IF NOT EXISTS recipes
(
    id          UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    title       TEXT      NOT NULL,
    description TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMP,
    calories    INTEGER,
    takes_time  INTEGER,
    user_id     TEXT,
    hidden      BOOLEAN
);

CREATE TABLE IF NOT EXISTS steps
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id   UUID    NOT NULL,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    takes_time  INTEGER,
    required    BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT fk_recipe_id FOREIGN KEY (recipe_id)
        REFERENCES recipes (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_photos
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL,
    url       TEXT NOT NULL,
    CONSTRAINT fk_recipe_id FOREIGN KEY (recipe_id)
        REFERENCES recipes (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_tags
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID        NOT NULL,
    tag       VARCHAR(30) NOT NULL,
    CONSTRAINT fk_recipe_id FOREIGN KEY (recipe_id)
        REFERENCES recipes (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS step_photos
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    step_id UUID NOT NULL,
    url     TEXT NOT NULL,
    CONSTRAINT fk_recipe_id FOREIGN KEY (step_id)
        REFERENCES steps (id)
        ON DELETE CASCADE
);
CREATE INDEX ON steps (id, recipe_id);
CREATE INDEX ON recipes (id);
-- trigger to update updated_at on update
DROP TRIGGER IF EXISTS update_on_recipe ON recipes;
CREATE FUNCTION update_on_recipe()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_on_recipe
    BEFORE UPDATE
    ON
        recipes
    FOR EACH ROW
EXECUTE PROCEDURE update_on_recipe();