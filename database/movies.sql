CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    tmdb_id INT NOT NULL,
    title TEXT NOT NULL,
    overview TEXT,
    poster_url TEXT,
    rank INT,
    user_id INT
);

