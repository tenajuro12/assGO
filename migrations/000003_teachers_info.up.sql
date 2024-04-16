CREATE TABLE IF NOT EXISTS teachers_info (
                                             id SERIAL PRIMARY KEY,
                                             name VARCHAR(100),
                                             surname VARCHAR(100),
                                             email VARCHAR(100),
                                             modules INT REFERENCES module_info(id)
);