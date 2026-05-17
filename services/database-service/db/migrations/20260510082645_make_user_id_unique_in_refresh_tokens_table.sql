-- +goose Up
ALTER TABLE public.refresh_tokens
  ADD CONSTRAINT "UQ_refresh_tokens_user_id" UNIQUE (user_id);

-- +goose Down
ALTER TABLE public.refresh_tokens
  DROP CONSTRAINT "UQ_refresh_tokens_user_id";
