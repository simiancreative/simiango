CREATE DATABASE examples;
use examples;

CREATE TABLE products (
    id serial,
    name varchar(50)
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

