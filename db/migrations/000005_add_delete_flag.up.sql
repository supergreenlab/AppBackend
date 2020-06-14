alter table boxes add column deleted boolean default false;
alter table plants add column deleted boolean default false;
alter table timelapses add column deleted boolean default false;
alter table devices add column deleted boolean default false;
alter table feeds add column deleted boolean default false;
alter table feedentries add column deleted boolean default false;
alter table feedmedias add column deleted boolean default false;
