CREATE DATABASE RINHA_BACKEND;

CREATE TABLE IF NOT EXISTS person (
    user_id VARCHAR(64) PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    user_nick VARCHAR(32) NOT NULL,
    user_birth DATE NOT NULL,
    user_stack VARCHAR(1024)
)
