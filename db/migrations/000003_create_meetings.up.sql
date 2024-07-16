CREATE TABLE IF NOT EXISTS meetings (
  id character varying NOT NULL PRIMARY KEY,

  title character varying NOT NULL,
  meeting_description character varying NOT NULL,

  start_time timestamp with time zone NOT NULL,
  end_time timestamp with time zone NOT NULL,

  created_at timestamp with time zone,
  updated_at timestamp with time zone
);
