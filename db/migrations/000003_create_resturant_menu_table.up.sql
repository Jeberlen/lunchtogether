CREATE TABLE IF NOT EXISTS restaurant_menu (
    restaurant_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    FOREIGN KEY (restaurant_id) REFERENCES restaurant(id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_item(id),
    PRIMARY KEY (restaurant_id, menu_item_id)
);