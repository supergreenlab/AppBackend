alter table comments add admin_checked boolean not null default false;

create index c_ac on comments (admin_checked);
