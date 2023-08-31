CREATE TABLE posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  title TEXT NOT NULL CHECK(length(title) <= 100),  
  content TEXT NOT NULL CHECK(length(content) <= 10000), 
  image_url TEXT,
  privacy INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id INTEGER,
  user_id INTEGER,
  content TEXT NOT NULL CHECK(length(content) <= 1000),
  image_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Trigger to delete related post comments when a post is deleted
CREATE TRIGGER delete_related_post_comments
AFTER DELETE ON posts
FOR EACH ROW
BEGIN
  DELETE FROM comments WHERE post_id = OLD.id;
END;



CREATE TABLE post_followers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id INTEGER,
  user_id INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Trigger to delete related post followers when a post is deleted
CREATE TRIGGER delete_related_post_followers
AFTER DELETE ON posts
FOR EACH ROW
BEGIN
  DELETE FROM post_followers WHERE post_id = OLD.id;
END;