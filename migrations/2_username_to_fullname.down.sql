-- Migration to revert fullname back to username

-- Step 1: Add the username column
ALTER TABLE users ADD COLUMN username VARCHAR(50);

-- Step 2: Populate username with existing fullname data
UPDATE users SET username = fullname;

-- Step 3: Add the unique constraint back to username
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);

-- Step 4: Make username NOT NULL
ALTER TABLE users ALTER COLUMN username SET NOT NULL;

-- Step 5: Remove the fullname column
ALTER TABLE users DROP COLUMN fullname; 