CREATE TABLE IF NOT EXISTS teachers (
  id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
  user_id uuid NOT NULL UNIQUE REFERENCES public.users (id) ON DELETE CASCADE,
  description text,
  picture text
);
