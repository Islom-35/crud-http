create table if not exists movies (
    id bigserial primary key,
    name varchar(30) not null
);

CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    title varchar(30),
    author varchar(30),
    publish_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rating int
);
