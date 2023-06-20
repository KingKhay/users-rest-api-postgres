CREATE TABLE users (
                       id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       name VARCHAR(255),
                       email VARCHAR(255)
);