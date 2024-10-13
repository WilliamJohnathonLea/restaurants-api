BEGIN TRANSACTION;

CREATE TABLE restaurants (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

ALTER TABLE orders
ADD CONSTRAINT fk_restaurant FOREIGN KEY(restaurant_id) REFERENCES restaurants(id);

COMMIT;
