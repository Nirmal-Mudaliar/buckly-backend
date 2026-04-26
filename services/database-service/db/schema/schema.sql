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