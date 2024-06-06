-- 1. Drop the existing foreign key constraint
ALTER TABLE public.courses
DROP CONSTRAINT courses_teacher_id_fkey;

-- 2. Add the new foreign key constraint
ALTER TABLE public.courses ADD CONSTRAINT courses_teacher_id_fkey FOREIGN KEY (teacher_id) REFERENCES public.users (id);
