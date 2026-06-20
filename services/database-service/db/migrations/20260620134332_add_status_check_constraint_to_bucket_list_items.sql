-- +goose Up
ALTER TABLE public.bucket_list_items
  ADD CONSTRAINT status_check_constraint
  CHECK (status IN ('planned', 'in_progress', 'completed', 'abandoned'));

-- +goose Down
ALTER TABLE public.bucket_list_items
  DROP CONSTRAINT status_check_constraint;