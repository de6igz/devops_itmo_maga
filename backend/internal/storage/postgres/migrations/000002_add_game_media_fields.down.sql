ALTER TABLE games
    DROP COLUMN IF EXISTS image_path,
    DROP COLUMN IF EXISTS description;
