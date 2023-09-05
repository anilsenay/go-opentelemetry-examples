CREATE TABLE todos (
  id BIGSERIAL CONSTRAINT todo_pk PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  completed BOOLEAN DEFAULT FALSE
);

-- Test Data
INSERT INTO todos (title, completed) VALUES ('Todo 1', true), ('Todo 2', false), ('Todo 3', false);