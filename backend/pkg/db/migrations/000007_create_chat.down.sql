
DROP TABLE chatrooms;

DROP TABLE chatroom_participants;

DROP TABLE messages;


DROP INDEX IF EXISTS idx_messages_chatroom ON messages(chatroom_id);
DROP INDEX IF EXISTS idx_messages_sender ON messages(sender_id);
DROP INDEX IF EXISTS idx_messages_recipient ON messages(recipient_id);