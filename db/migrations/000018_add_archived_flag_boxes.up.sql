alter table boxes add column archived boolean default false;
create index b_archived on boxes (archived);
