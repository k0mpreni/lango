CREATE TABLE IF NOT EXISTS teachers (
  id uuid PRIMARY KEY,
  user_id uuid references public.users (id) UNIQUE NOT NULL,
  description text,
  picture text
);
