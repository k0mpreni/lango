CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
  created_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    password_hash bytea,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1,
    email text UNIQUE NOT NULL,
    provider text UNIQUE NOT NULL,
    provider_id text UNIQUE,
    role smallint DEFAULT 1
);
