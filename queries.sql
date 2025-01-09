CREATE TABLE monitoring_data (
    host TEXT NOT NULL,
    type TEXT NOT NULL,
    parameter VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    insertion_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
