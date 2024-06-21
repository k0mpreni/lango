CREATE TABLE IF NOT EXISTS students (
  id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
  user_id uuid NOT NULL UNIQUE REFERENCES public.users (id) ON DELETE CASCADE
);
