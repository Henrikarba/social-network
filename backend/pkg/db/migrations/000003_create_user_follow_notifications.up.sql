CREATE TABLE user_notifications (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL, 
    group_id INTEGER DEFAULT 0,
    type VARCHAR(20),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (sender_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    UNIQUE (user_id, sender_id, type)
);


