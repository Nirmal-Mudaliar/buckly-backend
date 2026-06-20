-- +goose Up
CREATE TABLE public.states (
  id bigInt NOT NULL,
  country_id bigInt NOT NULL,
  name text NOT NULL
);

ALTER TABLE public.states ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.states_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.states
  ADD CONSTRAINT "PK_states" PRIMARY KEY (id);

ALTER TABLE public.states
  ADD CONSTRAINT "FK_states_countries_country_id" 
  FOREIGN KEY (country_id) 
  REFERENCES public.countries (id) 
  ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS public.states;
