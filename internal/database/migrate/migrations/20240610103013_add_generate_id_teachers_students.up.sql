CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE public.teachers
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();

ALTER TABLE public.students
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();
