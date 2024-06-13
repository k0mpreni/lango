CREATE TABLE IF NOT EXISTS students (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id uuid references public.users (id) UNIQUE NOT NULL
);
