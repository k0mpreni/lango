CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE public.courses
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();
