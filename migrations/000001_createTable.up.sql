CREATE TABLE IF NOT EXISTS Users  (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    referrer_id INTEGER NOT NULL DEFAULT 0,
    password_hash VARCHAR(64) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Tasks  (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    given_points INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Completed_tasks  (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    task_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (task_id) REFERENCES Tasks(id)
);

INSERT INTO Tasks (description, given_points) VALUES
('The first task', 10),
('The second task', 20),
('The third task', 30);