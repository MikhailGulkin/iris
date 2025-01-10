create keyspace if not exists chat with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS chat.chats
(
    id                UUID PRIMARY KEY,
    pinned_message_id UUID,
    conversation_id   UUID,
    type              TEXT,
    name              TEXT,
    deleted_at        TIMESTAMP,
    created_at        TIMESTAMP,
) WITH CDC = {'enabled': 'true'};
CREATE TABLE IF NOT EXISTS chat.messages
(
    id              UUID,
    chat_id         UUID,
    user_id         BIGINT,
    conversation_id UUID,
    text            TEXT,
    content_type    TEXT,
    send_at         TIMESTAMP,
    created_at      TIMESTAMP,
    deleted_at      TIMESTAMP,
    edit_at         TIMESTAMP,
    PRIMARY KEY ( id, chat_id, send_at)
) WITH CLUSTERING ORDER BY (chat_id ASC, send_at DESC) AND CDC = {'enabled': 'true'};
