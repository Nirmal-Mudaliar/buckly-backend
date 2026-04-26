-- +goose Up
ALTER TABLE public.users
ADD COLUMN password_hash TEXT NOT NULL;

-- +goose Down
ALTER TABLE public.users
DROP COLUMN password_hash;
