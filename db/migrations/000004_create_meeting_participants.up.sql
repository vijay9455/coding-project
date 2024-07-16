CREATE TABLE IF NOT EXISTS meeting_participants (
  id character varying NOT NULL PRIMARY KEY,
  meeting_id character varying NOT NULL,
  user_id character varying NOT NULL,
  accept_status character varying NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone
);

CREATE UNIQUE INDEX unique_idx_user_meeting_participants on meeting_participants USING btree (user_id, meeting_id);
CREATE INDEX idx_meeting_id_meeting_participants on meeting_participants USING btree (meeting_id);
