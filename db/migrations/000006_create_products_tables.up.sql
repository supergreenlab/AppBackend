create extension pg_trgm;

create table if not exists products(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  name varchar(256) not null,
  description varchar(4000) not null default '',

  filepath varchar not null,
  thumbnailpath varchar not null,

  categories jsonb not null default '[]'::jsonb,

  cat timestamptz default now(),
  uat timestamptz default now()
);

CREATE INDEX products_name ON products USING GIN(name gin_trgm_ops);

drop trigger if exists uat_products on products;

create trigger uat_products
before update on products
for each row
  execute procedure moddatetime(uat);

create table if not exists suppliers {
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  name varchar(256) not null,
  url varchar(256) not null,
  description varchar(4000) not null default '',
  locals jsonb not null default '[]'::jsonb,

  filepath varchar not null,
  thumbnailpath varchar not null,
  
  cat timestamptz default now(),
  uat timestamptz default now()
};

CREATE INDEX suppliers_name ON suppliers USING GIN(name gin_trgm_ops);

drop trigger if exists uat_suppliers on suppliers;

create trigger uat_suppliers
before update on suppliers
for each row
  execute procedure moddatetime(uat);

create table if not exists productsuppliers(
  id uuid primary key default uuid_generate_v4(),
  userid uuid not null,
  productid uuid not null,
  supplierid uuid not null,
  url varchar(256) not null,
  price numeric not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

drop trigger if exists uat_productsuppliers on productsuppliers;

create trigger uat_productsuppliers
before update on productsuppliers
for each row
  execute procedure moddatetime(uat);
