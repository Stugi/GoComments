DROP TABLE IF EXISTS comment;
DROP SEQUENCE IF EXISTS id_seq_comment;

CREATE SEQUENCE id_seq_comment;

CREATE TABLE comment (
    id INTEGER PRIMARY KEY DEFAULT nextval('id_seq_comment'),
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    news_id INTEGER ,
    parent_id INTEGER  REFERENCES comment(id),
    status INTEGER
    )