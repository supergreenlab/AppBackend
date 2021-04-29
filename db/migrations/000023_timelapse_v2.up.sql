alter table timelapses add column ttype varchar not null default 'dropbox';
alter table timelapses add column settings jsonb not null default '{}'::jsonb;

update timelapses set settings=jsonb_build_object('controllerid', t.controllerid, 'rotate', t.rotate, 'name', t.name, 'strain', t.strain, 'dropboxtoken', t.dropboxtoken, 'uploadname', t.uploadname) from (select * from timelapses) as t;

alter table timelapses drop column controllerid;
alter table timelapses drop column rotate;
alter table timelapses drop column name;
alter table timelapses drop column strain;
alter table timelapses drop column dropboxtoken;
alter table timelapses drop column uploadname;

create table if not exists timelapseframes(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  timelapseid uuid not null,

  filepath varchar not null,

  meta jsonb not null default '{}'::jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index tf_uid on timelapseframes (userid);
create index tf_tid on timelapseframes (timelapseid);

drop trigger if exists uat_timelapseframes on timelapseframes;

create trigger uat_timelapseframes
before update on timelapseframes
for each row
  execute procedure moddatetime(uat);
