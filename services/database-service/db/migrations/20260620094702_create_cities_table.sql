-- +goose Up
CREATE TABLE public.cities (
  id bigInt NOT NULL,
  state_id bigInt NOT NULL,
  name text NOT NULL
);

ALTER TABLE public.cities ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.cities_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.cities
  ADD CONSTRAINT "PK_cities" PRIMARY KEY (id);

ALTER TABLE public.cities
  ADD CONSTRAINT "FK_cities_states_state_id"
  FOREIGN KEY (state_id)
  REFERENCES public.states (id)
  ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS public.cities;
