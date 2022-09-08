CREATE DATABASE examples;
use examples;

CREATE TABLE products (
    id serial,
    name varchar(50),
);

CREATE TABLE ancestors (
    ancestor_id integer,
    product_id integer,
    depth integer
);

INSERT INTO products (id, name)
VALUES 
(1, 'Car'),
(2, 'Truck'),
(3, 'Motorcycle'),
(4, 'Bicycle'),
(5, 'Horse'),
(6, 'Boat'),
(7, 'Plane'),
(8, 'Scooter'),
(9, 'Skateboard');

INSERT INTO ancestors (ancestor_id, product_id)
VALUES 
(1, 1, 0),
(1, 2, 1),
(2, 2, 0),
(4, 4, 0),
(4, 3, 1),
(5, 5, 0),
(6, 6, 0),
(7, 7, 0),
(8, 8, 0),
(9, 9, 0);
