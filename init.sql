\c bot
Create Table products (
    ID varchar(255) PRIMARY KEY,
    Name varchar(255) NOT NULL,
    Price numeric
);

INSERT INTO products (ID, Name, Price) VALUES ('1', 'carro', 56.00);
INSERT INTO products (ID, Name, Price) VALUES ('2', 'moto', 10.00);
INSERT INTO products (ID, Name, Price) VALUES ('3', 'barco', 70.00);
INSERT INTO products (ID, Name, Price) VALUES ('4', 'aviao', 90.00);
INSERT INTO products (ID, Name, Price) VALUES ('5', 'tartaruga', 6666.00);
