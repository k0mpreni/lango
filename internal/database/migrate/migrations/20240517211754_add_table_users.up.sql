CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY,
  created_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    email text UNIQUE NOT NULL,
    user_id uuid,
    role smallint DEFAULT 1
);
