CREATE TABLE public.users (
  id bigInt NOT NULL,
  first_name character varying(100) NOT NULL,
  last_name character varying(100) NOT NULL,
  email character varying(254) NOT NULL,
  phone_no character varying(15),
  date_of_birth date NOT NULL,
  gender text NOT NULL,
  bio character varying(500),
  profile_photo_url text,
  home_country_id bigInt,
  home_state_id bigInt,
  home_city_id bigInt,
  password_hash text NOT NULL,
  is_phone_verified boolean DEFAULT false,
  trust_score integer DEFAULT 0,
  status text DEFAULT 'active',
  insert_ts timestamp with time zone NOT NULL,
  modified_ts timestamp with time zone NOT NULL
);

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.users_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.users
  ADD CONSTRAINT "PK_users" PRIMARY KEY (id);

ALTER TABLE public.users
  ADD CONSTRAINT "users_email_unique" UNIQUE (email);

ALTER TABLE public.users
  ADD CONSTRAINT "users_phone_no_unique" UNIQUE (phone_no);


CREATE TABLE public.refresh_tokens (
  id bigInt NOT NULL,
  user_id bigInt NOT NULL,
  token text NOT NULL,
  expires_at timestamp with time zone NOT NULL,
  insert_ts timestamp with time zone NOT NULL
);

ALTER TABLE public.refresh_tokens ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.refresh_tokens_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
);

ALTER TABLE public.refresh_tokens
  ADD CONSTRAINT "PK_refresh_tokens" PRIMARY KEY (id);

ALTER TABLE public.refresh_tokens 
  ADD CONSTRAINT "FK_refresh_token_users_id" FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens(expires_at);


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

ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT status_check_constraint
  CHECK (status IN ('planned', 'in_progress', 'completed', 'abandoned'));

CREATE INDEX idx_bucket_list_items_user_id ON public.bucket_list_items (user_id);
CREATE INDEX idx_bucket_list_items_activity_tag_id ON public.bucket_list_items (activity_tag_id);