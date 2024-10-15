BEGIN TRANSACTION;

CREATE TABLE menus(
    id uuid PRIMARY KEY,
    restaurant_id uuid NOT NULL,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT fk_restaurant FOREIGN KEY(restaurant_id) REFERENCES restaurants(id)
);

CREATE TABLE menu_items(
    id uuid PRIMARY KEY,
    menu_id uuid NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    CONSTRAINT fk_menu FOREIGN KEY(menu_id) REFERENCES menus(id)
);

COMMIT;
