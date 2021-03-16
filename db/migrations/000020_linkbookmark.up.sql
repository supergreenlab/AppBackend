create table if not exists linkbookmarks(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,

  url varchar(4096) not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index lbo_uid on linkbookmarks (userid);

drop trigger if exists uat_linkbookmarks on linkbookmarks;

create trigger uat_linkbookmarks
before update on linkbookmarks
for each row
  execute procedure moddatetime(uat);
