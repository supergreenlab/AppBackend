create table if not exists follows(
  id uuid primary key default uuid_generate_v4(),

  userid uuid not null,
  plantid uuid not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index fo_uid on follows (userid);
create index fo_pid on follows (plantid);

drop trigger if exists uat_follows on follows;

create trigger uat_follows
before update on follows
for each row
  execute procedure moddatetime(uat);

