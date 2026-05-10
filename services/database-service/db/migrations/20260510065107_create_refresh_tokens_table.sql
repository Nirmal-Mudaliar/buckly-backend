-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS public.refresh_tokens;
