create extension if not exists pgcrypto;
create extension if not exists moddatetime;
create extension if not exists "uuid-ossp";

create table if not exists users(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  nickname varchar(64),
  password varchar(60),

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_users on users;

create trigger uat_users
before update on users
for each row
  execute procedure moddatetime(uat);

create table if not exists plants(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  feedID uuid not null,
  deviceID uuid,
  deviceBox int,
  name varchar(64),

  settings jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_plants on plants;

create trigger uat_plants
before update on plants
for each row
  execute procedure moddatetime(uat);

create table if not exists timelapses(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  plantID uuid not null,
  controllerID varchar(64),
  rotate varchar(5),
  name varchar(64),
  strain varchar(64),
  uploadName varchar(64),

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_timelapses on timelapses;

create trigger uat_timelapses
before update on timelapses
for each row
  execute procedure moddatetime(uat);

create table if not exists devices(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  identifier varchar(16),
  name varchar(24),
  ip varchar(15),
  mdns varchar(64),

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_devices on devices;

create trigger uat_devices
before update on devices
for each row
  execute procedure moddatetime(uat);

create table if not exists feeds(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name varchar(24),

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_feeds on feeds;

create trigger uat_feeds
before update on feeds
for each row
  execute procedure moddatetime(uat);

create table if not exists feedEntries(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  feedID uuid,
  createdAt timestamptz,

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_feedEntries on feedEntries;

create trigger uat_feedEntries
before update on feedEntries
for each row
  execute procedure moddatetime(uat);
