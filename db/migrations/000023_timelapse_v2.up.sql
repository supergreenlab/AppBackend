alter table timelapses add column ttype varchar not null default 'dropbox';
alter table timelapses add column settings jsonb not null default '{}'::jsonb;

update timelapses set settings=jsonb_build_object('controllerid', t.controllerid, 'rotate', t.rotate, 'name', t.name, 'strain', t.strain, 'dropboxtoken', t.dropboxtoken, 'uploadname', t.uploadname) from (select * from timelapses) as t;

alter table timelapses drop column controllerid;
alter table timelapses drop column rotate;
alter table timelapses drop column name;
alter table timelapses drop column strain;
alter table timelapses drop column dropboxtoken;
alter table timelapses drop column uploadname;
