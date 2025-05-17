CREATE TABLE user_movies (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    movie_id INT REFERENCES movies(id),
    rank INT CHECK (rank >= 1 AND rank <= 100),
    notes TEXT,
    UNIQUE(user_id, rank),
    UNIQUE(user_id, movie_id)
);
