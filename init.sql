CREATE ROLE root WITH LOGIN PASSWORD 'postgres';

CREATE TABLE IF NOT EXISTS restaurant (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS menu_item (
    id SERIAL PRIMARY KEY,
    dayOfWeek TEXT NOT NULL,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS restaurant_menu (
    restaurant_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    FOREIGN KEY (restaurant_id) REFERENCES restaurant(id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_item(id),
    PRIMARY KEY (restaurant_id, menu_item_id)
);