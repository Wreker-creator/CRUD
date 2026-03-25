-- Now why did we create this folder and file?

-- Well, in a real project your database schema evolves over time.
-- You add a column, rename a table, add indexes.
-- Keeping changes in numbered SQL files (migrations) 
-- allows you to track and apply changes in a controlled manner.:
-- 1. Anyone can recreate the exact schema from scratch
-- 2. We have a visible history
-- 3. Toold like golang-migrate can apply these automatically.

CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY, -- our autoincrement integer, unique, no null
    title VARCHAR(255) NOT NULL, -- a string with max length 255, cannot be null
    description TEXT, -- a string with no length limit, can be null
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), --useful for debugging
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() -- our last modifcation
);