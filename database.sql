CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    saldo NUMERIC(10, 2) NOT NULL CHECK (saldo >= 0)
);

CREATE TABLE transferencias (
    id SERIAL PRIMARY KEY,
    id_origen INT NOT NULL REFERENCES usuarios(id),
    id_destino INT NOT NULL REFERENCES usuarios(id),
    monto NUMERIC(10, 2) NOT NULL CHECK (monto > 0),
);

INSERT INTO usuarios (nombre, saldo) 
VALUES 
    ('Jose Manuel', 3000.00),
    ('Iram Maximiliano', 8000.00),
    ('Diego Benito', 1000.00),
    ('Jesus Alexis', 5000.00),
    ('Aaron', 7000.00);