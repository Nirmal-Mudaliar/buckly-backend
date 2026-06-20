-- +goose Up
CREATE TABLE public.bucket_list_items (
  id bigInt NOT NULL,
  user_id bigInt NOT NULL,
  activity_tag_id bigInt NOT NULL,
  country_id bigInt NOT NULL,
  state_id bigInt NOT NULL,
  city_id bigInt NOT NULL,
  timeframe_start_date date NOT NULL,
  timeframe_end_date date NOT NULL,
  note text,
  is_public boolean NOT NULL DEFAULT true,
  status text NOT NULL DEFAULT 'planned',
  insert_ts timestamp with time zone NOT NULL,
  modified_ts timestamp with time zone NOT NULL 
);

ALTER TABLE public.bucket_list_items ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.bucket_list_items_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "PK_bucket_list_items" PRIMARY KEY (id);

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "FK_bucket_list_items_users_user_id"
  FOREIGN KEY (user_id)
  REFERENCES public.users (id)
  ON DELETE CASCADE;

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "FK_bucket_list_items_activity_tags_activity_tag_id"
  FOREIGN KEY (activity_tag_id)
  REFERENCES public.activity_tags (id)
  ON DELETE CASCADE;

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "FK_bucket_list_items_countries_country_id"
  FOREIGN KEY (country_id)
  REFERENCES public.countries (id)
  ON DELETE CASCADE;

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "FK_bucket_list_items_states_state_id"
  FOREIGN KEY (state_id)
  REFERENCES public.states (id)
  ON DELETE CASCADE;

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT "FK_bucket_list_items_cities_city_id"
  FOREIGN KEY (city_id)
  REFERENCES public.cities (id)
  ON DELETE CASCADE;

CREATE INDEX idx_bucket_list_items_user_id ON public.bucket_list_items (user_id);
CREATE INDEX idx_bucket_list_items_activity_tag_id ON public.bucket_list_items (activity_tag_id);

-- +goose Down
DROP TABLE IF EXISTS public.bucket_list_items;
