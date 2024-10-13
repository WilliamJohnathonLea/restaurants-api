BEGIN TRANSACTION;

CREATE TABLE orders(
    id uuid PRIMARY KEY,
    restaurant_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

CREATE TABLE line_items(
    id uuid PRIMARY KEY,
    order_id uuid NOT NULL,
    item_id uuid NOT NULL,
    name VARCHAR(255) NOT NULL,
    price NUMERIC NOT NULL,
    quantity INTEGER NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id),
    CONSTRAINT check_quantity CHECK(quantity > 0)
);

COMMIT;
