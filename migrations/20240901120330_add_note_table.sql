-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.note
(
    id serial primary key,
    rf_user_id integer references public.user(id) on delete cascade not null,
    body text,
    deleted boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.note;
-- +goose StatementEnd
