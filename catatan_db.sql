CREATE TABLE userpass (
	id INTEGER PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	display_name TEXT NOT NULL
);