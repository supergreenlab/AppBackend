create table if not exists likes(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  feedentryid uuid,
  commentid uuid,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index l_feid on likes (feedentryid);
create index l_uid on likes (userid);
create index l_cid on likes (commentid);

drop trigger if exists uat_likes on likes;

create trigger uat_likes
before update on likes
for each row
  execute procedure moddatetime(uat);
