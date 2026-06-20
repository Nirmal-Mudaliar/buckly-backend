-- +goose Up
CREATE TABLE public.countries (
  id bigInt NOT NULL,
  name text NOT NULL,
  code text NOT NULL
);

ALTER TABLE public.countries ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.countries_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.countries
  ADD CONSTRAINT "PK_countries" PRIMARY KEY (id);

ALTER TABLE public.countries
  ADD CONSTRAINT "countries_code_unique" UNIQUE (code);


-- +goose Down
DROP TABLE IF EXISTS public.countries;
