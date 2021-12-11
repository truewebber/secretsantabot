CREATE TABLE chats
(
    id          BIGSERIAL NOT NULL PRIMARY KEY,
    tg_chat_id  BIGINT    NOT NULL,
    tg_admin_id BIGINT    NOT NULL,
    deleted     BOOLEAN   NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS chats__tg_chat_id__uidx ON chats (tg_chat_id);
CREATE INDEX IF NOT EXISTS chats__tg_admin_id__idx ON chats (tg_admin_id);

CREATE TABLE users
(
    id         BIGSERIAL NOT NULL PRIMARY KEY,
    tg_chat_id BIGINT    NOT NULL,
    tg_user_id BIGINT    NOT NULL,
    deleted    BOOLEAN   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS participants__tg_chat_id__tg_users_id__uidx ON users (tg_chat_id, tg_user_id);

CREATE TABLE magic_chat_history
(
    id         BIGSERIAL NOT NULL PRIMARY KEY,
    chat_id    BIGINT    NOT NULL,
    version    SMALLINT  NOT NULL,
    status     SMALLINT  NOT NULL,
    deleted    BOOLEAN   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (chat_id) REFERENCES chats (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS magic_chat_history__tg_chat_id__version__uidx ON magic_chat_history (chat_id, version);

CREATE TABLE magic_participants
(
    id                    BIGSERIAL NOT NULL PRIMARY KEY,
    magic_chat_history_id BIGINT    NOT NULL,
    participant_id        BIGINT    NOT NULL,
    deleted               BOOLEAN   NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at            TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (magic_chat_history_id) REFERENCES magic_chat_history (id),
    FOREIGN KEY (participant_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS magic_participants__chat_id__participant_id__uidx
    ON magic_participants (magic_chat_history_id, participant_id);

CREATE TABLE magic_results
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
CREATE UNIQUE INDEX IF NOT EXISTS magic__history_id__giver_id__receiver_id__uidx
    ON magic_results (magic_chat_history_id, participant_giver_id, participant_receiver_id);