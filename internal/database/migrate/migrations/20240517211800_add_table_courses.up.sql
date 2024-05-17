CREATE TABLE IF NOT EXISTS courses (
  id uuid PRIMARY KEY,
  created_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    teacher_id uuid references public.users (id) NOT NULL,
    student_id uuid references public.users (id) NOT NULL,
    title varchar(200),
    description text,
    date timestamp(0)
  with
    time zone NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    canceled BOOLEAN NOT NULL DEFAULT FALSE
);
