create extension if not exists pgcrypto;
create extension if not exists moddatetime;
create extension if not exists "uuid-ossp";

create table if not exists users(
  id uuid primary key default uuid_generate_v4(),
  nickname varchar(64) not null,
  password varchar(60) not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_users on users;

create trigger uat_users
before update on users
for each row
  execute procedure moddatetime(uat);

create table if not exists boxes(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  deviceid uuid,
  devicebox int,
  name varchar(64) not null,

  settings jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index b_uid on boxes (userid);

drop trigger if exists uat_boxes on boxes;

create trigger uat_boxes
before update on boxes
for each row
  execute procedure moddatetime(uat);

create table if not exists plants(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  boxid uuid not null,
  feedid uuid not null,
  name varchar(64) not null,
  single boolean not null default true,

  settings jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index p_uid on plants (userid);

drop trigger if exists uat_plants on plants;

create trigger uat_plants
before update on plants
for each row
  execute procedure moddatetime(uat);

create table if not exists timelapses(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  plantid uuid not null,
  controllerid varchar(64) not null,
  rotate varchar(5) not null default 'false',
  name varchar(64) not null,
  strain varchar(64) not null,
  dropboxtoken varchar(64) not null,
  uploadname varchar(64) not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index t_uid on timelapses (userid);
create index t_pid on timelapses (plantid);

drop trigger if exists uat_timelapses on timelapses;

create trigger uat_timelapses
before update on timelapses
for each row
  execute procedure moddatetime(uat);

create table if not exists devices(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  identifier varchar(16) not null,
  name varchar(24) not null,
  ip varchar(15) not null,
  mdns varchar(64) not null,
  config varchar not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index d_uid on devices (userid);

drop trigger if exists uat_devices on devices;

create trigger uat_devices
before update on devices
for each row
  execute procedure moddatetime(uat);

create table if not exists feeds(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  name varchar(24) not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index f_uid on feeds (userid);

drop trigger if exists uat_feeds on feeds;

create trigger uat_feeds
before update on feeds
for each row
  execute procedure moddatetime(uat);

create table if not exists feedentries(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  feedid uuid not null,
  etype varchar(24) not null,
  createdat timestamptz not null,

  params jsonb not null default '{}'::jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index fe_fid on feedentries (feedid);

drop trigger if exists uat_feedentries on feedentries;

create trigger uat_feedentries
before update on feedentries
for each row
  execute procedure moddatetime(uat);

create table if not exists feedmedias(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  feedentryid uuid not null,
  filepath varchar not null,
  thumbnailpath varchar not null,

  params jsonb not null default '{}'::jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index fm_feid on feedmedias (feedentryid);

drop trigger if exists uat_feedmedias on feedmedias;

create trigger uat_feedmedias
before update on feedmedias
for each row
  execute procedure moddatetime(uat);

create table if not exists plantsharings(
  userid uuid not null,
  plantid uuid not null,
  touserid uuid not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index ps_uuids on plantsharings (userid, plantid);

drop trigger if exists uat_plantsharings on plantsharings;

create trigger uat_plantsharings
before update on plantsharings
for each row
  execute procedure moddatetime(uat);

--
-- userend tables
--

create table if not exists userends(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index ue_uuids on userends (userid);

drop trigger if exists uat_userends on userends;

create trigger uat_userends
before update on userends
for each row
  execute procedure moddatetime(uat);

create table if not exists userend_boxes(
  userendid uuid not null,
  boxid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index ueb_uuids on userend_boxes (userendid, boxid);
create index ueb_sent on userend_boxes (sent);
create index ueb_dirty on userend_boxes (dirty);

drop trigger if exists uat_userend_boxes on userend_boxes;

create table if not exists userend_plants(
  userendid uuid not null,
  plantid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index uep_uuids on userend_plants (userendid, plantid);
create index uep_sent on userend_plants (sent);
create index uep_dirty on userend_plants (dirty);

drop trigger if exists uat_userend_plants on userend_plants;

create trigger uat_userend_plants
before update on userend_plants
for each row
  execute procedure moddatetime(uat);

create table if not exists userend_timelapses(
  userendid uuid not null,
  timelapseid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index uet_uuids on userend_timelapses (userendid, timelapseid);
create index uet_sent on userend_timelapses (sent);
create index uet_dirty on userend_timelapses (dirty);

drop trigger if exists uat_userend_plants on userend_plants;

create trigger uat_userend_plants
before update on userend_plants
for each row
  execute procedure moddatetime(uat);


create table if not exists userend_devices(
  userendid uuid not null,
  deviceid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index ued_uuids on userend_devices (userendid, deviceid);
create index ued_sent on userend_devices (sent);
create index ued_dirty on userend_devices (dirty);

drop trigger if exists uat_userend_devices on userend_devices;

create trigger uat_userend_devices
before update on userend_devices
for each row
  execute procedure moddatetime(uat);

create table if not exists userend_feeds(
  userendid uuid not null,
  feedid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index uef_uuids on userend_feeds (userendid, feedid);
create index uef_sent on userend_feeds (sent);
create index uef_dirty on userend_feeds (dirty);

drop trigger if exists uat_userend_feeds on userend_feeds;

create trigger uat_userend_feeds
before update on userend_feeds
for each row
  execute procedure moddatetime(uat);

create table if not exists userend_feedentries(
  userendid uuid not null,
  feedentryid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index uefe_uuids on userend_feedentries (userendid, feedentryid);
create index uefe_sent on userend_feedentries (sent);
create index uefe_dirty on userend_feedentries (dirty);

drop trigger if exists uat_userend_feedentries on userend_feedentries;

create trigger uat_userend_feedentries
before update on userend_feedentries
for each row
  execute procedure moddatetime(uat);

create table if not exists userend_feedmedias(
  userendid uuid not null,
  feedmediaid uuid not null,

  sent boolean not null default false,
  dirty boolean not null default false,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index uefm_uuids on userend_feedmedias (userendid, feedmediaid);
create index uefm_sent on userend_feedmedias (sent);
create index uefm_dirty on userend_feedmedias (dirty);

drop trigger if exists uat_userend_feedmedias on userend_feedmedias;

create trigger uat_userend_feedmedias
before update on userend_feedmedias
for each row
  execute procedure moddatetime(uat);
