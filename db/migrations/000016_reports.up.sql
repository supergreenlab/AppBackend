create table if not exists reports(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,

  feedentryid uuid,
  commentid uuid,
  plantid uuid,

  rtype varchar(50) not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index r_feid on reports (feedentryid);
create index r_cid on reports (commentid);
create index r_pid on reports (plantid);

drop trigger if exists uat_reports on reports;
create trigger uat_reports
before update on reports
for each row
  execute procedure moddatetime(uat);
