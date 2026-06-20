-- +goose Up
CREATE TABLE public.activity_categories (
  id bigInt NOT NULL,
  name text NOT NULL
);

ALTER TABLE public.activity_categories ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.activity_categories_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.activity_categories
  ADD CONSTRAINT "PK_activity_categories"
  PRIMARY KEY (id);

-- +goose Down
DROP TABLE IF EXISTS public.activity_categories;
