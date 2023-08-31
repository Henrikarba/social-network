CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL CHECK(length(title) <= 100),  -- max 100 characters for title
    description TEXT CHECK(length(description) <= 2000),  -- max 2,000 characters for description
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);
CREATE TABLE group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT CHECK (status IN ('invited', 'joined', 'requested', 'rejected')),
    invited_by INTEGER,  
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (invited_by) REFERENCES users(id)
);
-- Create the trigger to add the group owner as a member
CREATE TRIGGER add_group_owner_as_member
AFTER INSERT ON groups
FOR EACH ROW
BEGIN
  INSERT INTO group_members (group_id, user_id, status, invited_by, joined_at)
  VALUES (NEW.id, NEW.creator_id, "joined", NEW.creator_id, NEW.created_at);
END;

CREATE TRIGGER delete_group_followers
AFTER DELETE ON groups
FOR EACH ROW
BEGIN
  DELETE FROM group_members WHERE group_id = OLD.id;
END;

CREATE TABLE group_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,  -- the author of the post
    title TEXT NOT NULL CHECK(length(title) <= 100),  -- max 100 characters for title
    content TEXT NOT NULL CHECK(length(content) <= 10000),  -- max 10,000 characters for content
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE group_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,  
    content TEXT NOT NULL CHECK(length(content) <= 1000), -- max 1000 characters for comments
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image_url TEXT,
    FOREIGN KEY (post_id) REFERENCES group_posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);