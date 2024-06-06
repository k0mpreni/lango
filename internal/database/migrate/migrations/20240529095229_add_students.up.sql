CREATE TABLE IF NOT EXISTS students (
  id uuid PRIMARY KEY,
  user_id uuid references public.users (id) UNIQUE NOT NULL
);
