CREATE TABLE IF NOT EXISTS chats
(
    id            BIGINT    NOT NULL PRIMARY KEY,
    admin_user_id BIGINT    NOT NULL,
    deleted       BOOLEAN   NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at    TIMESTAMP NOT NULL DEFAULT current_timestamp
);
CREATE INDEX IF NOT EXISTS chats__admin_user_id__idx ON chats (admin_user_id);

CREATE TABLE IF NOT EXISTS magic_chat_history
(
    id         BIGSERIAL NOT NULL PRIMARY KEY,
    chat_id    BIGINT    NOT NULL,
    deleted    BOOLEAN   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (chat_id) REFERENCES chats (id)
);
CREATE INDEX IF NOT EXISTS magic_chat_history__chat_id__idx ON magic_chat_history (chat_id);

CREATE TABLE IF NOT EXISTS magic_participants
(
    id                    BIGSERIAL NOT NULL PRIMARY KEY,
    magic_chat_history_id BIGINT    NOT NULL,
    participant_user_id   BIGINT    NOT NULL,
    deleted               BOOLEAN   NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at            TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (magic_chat_history_id) REFERENCES magic_chat_history (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS magic_participants__chat_id__participant_id__uidx
    ON magic_participants (magic_chat_history_id, participant_user_id);

CREATE TABLE IF NOT EXISTS magic_results
(
    id                      BIGSERIAL NOT NULL PRIMARY KEY,
    magic_chat_history_id   BIGINT    NOT NULL,
    participant_giver_id    BIGINT    NOT NULL,
    participant_receiver_id BIGINT    NOT NULL,
    deleted                 BOOLEAN   NOT NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (magic_chat_history_id) REFERENCES magic_chat_history (id),
    FOREIGN KEY (participant_giver_id) REFERENCES magic_participants (id),
    FOREIGN KEY (participant_receiver_id) REFERENCES magic_participants (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS magic__history_id__giver_id__uidx
    ON magic_results (magic_chat_history_id, participant_giver_id);
CREATE UNIQUE INDEX IF NOT EXISTS magic__history_id__receiver_id__uidx
    ON magic_results (magic_chat_history_id, participant_receiver_id);