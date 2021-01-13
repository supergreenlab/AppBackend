create table if not exists comments(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  feedentryid uuid not null,

  replyto uuid,
  text varchar not null,
  ctype varchar not null,
  params jsonb not null default '{}'::jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index c_feid on comments (feedentryid);

drop trigger if exists uat_comments on comments;

create trigger uat_comments
before update on comments
for each row
  execute procedure moddatetime(uat);
