-- +goose Up
CREATE TABLE public.activity_tags (
  id bigInt NOT NULL,
  category_id bigInt NOT NULL,
  name text NOT NULL,
  description text
);

ALTER TABLE public.activity_tags ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.activity_tags_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.activity_tags
  ADD CONSTRAINT "PK_activity_tags"
  PRIMARY KEY (id);

ALTER TABLE public.activity_tags
  ADD CONSTRAINT "FK_activity_tags_activity_categories_category_id"
  FOREIGN KEY (category_id)
  REFERENCES public.activity_categories (id)
  ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS public.activity_tags;
