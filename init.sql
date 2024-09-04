\c bot
Create Table products (
    ID varchar(255) PRIMARY KEY,
    Name varchar(255) NOT NULL,
    Price numeric
);

Create Table users (
    cpf varchar(255) PRIMARY KEY,
    name varchar(255),
    phone varchar(255)
);

Create Table orders(
    pedido varchar(255),
    produto varchar(255),
    FOREIGN KEY (produto) REFERENCES products(ID)
);

Create Table issues (
    issues varchar(255),
    cpf varchar(255),
    FOREIGN KEY (cpf) REFERENCES users(cpf)
);

INSERT INTO products (ID, Name, Price) VALUES ('1', 'carro', 56.00);
INSERT INTO products (ID, Name, Price) VALUES ('2', 'moto', 10.00);
INSERT INTO products (ID, Name, Price) VALUES ('3', 'barco', 70.00);
INSERT INTO products (ID, Name, Price) VALUES ('4', 'aviao', 90.00);
INSERT INTO products (ID, Name, Price) VALUES ('5', 'tartaruga', 6666.00);

INSERT INTO users (cpf, name, phone) VALUES ('123', 'Rafael', '819');
INSERT INTO users (cpf, name, phone) VALUES ('321', 'miguel', '918');
INSERT INTO users (cpf, name, phone) VALUES ('213', 'luna', '198');

INSERT INTO issues (issues, cpf) VALUES ('problema de conexao', '123');
INSERT INTO issues (issues, cpf) VALUES ('problema de rede', '321');
INSERT INTO issues (issues, cpf) VALUES ('problema de dinheiro', '213');

INSERT INTO orders (pedido,produto) values ('133','1');
INSERT INTO orders (pedido,produto) values ('133','5');
INSERT INTO orders (pedido,produto) values ('133','4');

