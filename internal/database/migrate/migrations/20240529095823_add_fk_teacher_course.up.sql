ALTER TABLE public.courses
DROP CONSTRAINT courses_teacher_id_fkey;

ALTER TABLE public.courses ADD CONSTRAINT courses_teacher_id_fkey FOREIGN KEY (teacher_id) REFERENCES public.teachers (id);
