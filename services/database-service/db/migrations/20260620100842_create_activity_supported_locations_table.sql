-- +goose Up
CREATE TABLE public.activity_supported_locations (
  id bigInt NOT NULL,
  activity_tag_id bigInt NOT NULL,
  country_id bigInt NOT NULL,
  state_id bigInt,
  city_id bigInt,
  is_active boolean NOT NULL DEFAULT true
);

ALTER TABLE public.activity_supported_locations ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.activity_supported_locations_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.activity_supported_locations
  ADD CONSTRAINT "PK_activity_supported_locations" PRIMARY KEY (id);

ALTER TABLE public.activity_supported_locations
  ADD CONSTRAINT "FK_activity_supported_locations_activity_tags_activity_tag_id"
  FOREIGN KEY (activity_tag_id)
  REFERENCES public.activity_tags (id)
  ON DELETE CASCADE;

ALTER TABLE public.activity_supported_locations
  ADD CONSTRAINT "FK_activity_supported_locations_countries_country_id"
  FOREIGN KEY (country_id)
  REFERENCES public.countries (id)
  ON DELETE CASCADE;

ALTER TABLE public.activity_supported_locations
  ADD CONSTRAINT "FK_activity_supported_locations_states_state_id"
  FOREIGN KEY (state_id)
  REFERENCES public.states (id)
  ON DELETE CASCADE;

ALTER TABLE public.activity_supported_locations
  ADD CONSTRAINT "FK_activity_supported_locations_cities_city_id"
  FOREIGN KEY (city_id)
  REFERENCES public.cities (id)
  ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS public.activity_supported_locations;