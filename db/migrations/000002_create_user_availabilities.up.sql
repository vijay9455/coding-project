CREATE TABLE IF NOT EXISTS user_availabilities (
  id character varying NOT NULL PRIMARY KEY,
  user_id character varying NOT NULL,
  day_of_week integer NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone
);

CREATE INDEX idx_user_id_on_user_availabilities on user_availabilities USING btree (user_id);
