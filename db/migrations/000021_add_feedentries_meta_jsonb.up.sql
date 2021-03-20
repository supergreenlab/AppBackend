alter table feedentries add column meta jsonb not null default '{}'::jsonb;
