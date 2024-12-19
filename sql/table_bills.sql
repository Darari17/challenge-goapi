CREATE TABLE bills (
    id SERIAL PRIMARY KEY,
    bill_date DATE NOT NULL,
    entry_date DATE NOT NULL,
    finish_date DATE,
    employee_id INT NOT NULL,
    customer_id INT NOT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);