CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES product(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);

ALTER TABLE transaction_details ENABLE ROW LEVEL SECURITY;