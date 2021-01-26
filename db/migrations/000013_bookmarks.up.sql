create table if not exists bookmarks(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  feedentryid uuid,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index b_feid on bookmarks (feedentryid);
create index bo_uid on bookmarks (userid);

drop trigger if exists uat_bookmarks on bookmarks;

create trigger uat_bookmarks
before update on bookmarks
for each row
  execute procedure moddatetime(uat);
