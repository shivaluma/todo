-- Migration to replace username with fullname

-- Step 1: Add the new fullname column
ALTER TABLE users ADD COLUMN fullname VARCHAR(100);

-- Step 2: Populate fullname with existing username data
UPDATE users SET fullname = username;

-- Step 3: Remove any constraints or indexes related to username
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;

-- Step 4: Remove the username column
ALTER TABLE users DROP COLUMN username; 