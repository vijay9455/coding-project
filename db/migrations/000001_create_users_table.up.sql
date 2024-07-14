CREATE TABLE IF NOT EXISTS users (
  id character varying NOT NULL,
  email character varying NOT NULL,
  first_name character varying NOT NULL,
  last_name character varying NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone
);

CREATE UNIQUE INDEX unique_idx_users_email on users USING btree (email);
