create table users
(
  id            serial       not null unique,
  name          varchar(255) not null,
  username      varchar(255) not null unique,
  password_hash varchar(255) not null
);

create table timeslots_lists
(
  id    serial       not null unique,
  title varchar(255) not null
);

create table users_lists
(
  id      serial not null unique,
  user_id int    references users (id) on delete cascade not null,
  list_id int    references timeslots_lists (id) on delete cascade not null
);

create table timeslots_items
(
  id          serial       not null unique,
  title       varchar(255) not null,
  description varchar(255),
  begining    timestamp    not null,
  finish      timestamp    not null
);

create table lists_items
(
  id       serial not null unique,
  item_id  int    references timeslots_items (id) on delete cascade not null,
  lists_id int    references timeslots_lists (id) on delete cascade not null
)