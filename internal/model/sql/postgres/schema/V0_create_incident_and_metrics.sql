-- public.incident definition
-- Drop table
-- DROP TABLE public.incident;
CREATE TABLE IF NOT EXISTS public.incident (
	id serial NOT NULL,
	title text NULL,
	start_ts timestamptz NULL,
	end_ts timestamptz NULL,
	identification_ts timestamptz NULL,
	responsibility varchar(50) NULL,
	team text NULL,
	functionality text NULL,
	root_cause text NULL,
	customer_impact int4 NULL,
	meeting_url text NULL,
	status_page_url text NULL,
	post_mortem_url text NULL,
	status varchar(50) NULL,
	product varchar(50) NULL,
	severity_level int4 NULL,
	channel_name text NULL,
	updated_at timestamp NOT NULL DEFAULT now(),
  snoozed_until timestamptz NULL,
	description_started text NULL,
	description_cancelled text NULL,
	description_resolved text NULL,
	channel_id varchar(50) NULL,
  commander_id text NULL,
  commander_email text NULL,
	CONSTRAINT firstkey PRIMARY KEY (id)
);

-- View table
-- DROP VIEW public.metrics;
CREATE OR REPLACE VIEW public.metrics
  AS SELECT
    incident.id,
    incident.title,
    incident.channel_id,
    incident.product,
    incident.team,
    incident.responsibility,
    incident.functionality,
    incident.root_cause,
    incident.severity_level,
    incident.status_page_url,
    incident.post_mortem_url,
    incident.channel_name,
    incident.status,
    incident.description_started,
    incident.description_cancelled,
    incident.description_resolved,
    incident.end_ts::date AS date,
    incident.commander_id,
    incident.commander_email,
    to_char(incident.start_ts, 'YYYY-MM-DD HH24:MI:SS'::text) AS start_ts,
    to_char(incident.end_ts, 'YYYY-MM-DD HH24:MI:SS'::text) AS end_ts,
    to_char(incident.identification_ts, 'YYYY-MM-DD HH24:MI:SS'::text) AS identification_ts,
    to_char(incident.updated_at, 'YYYY-MM-DD HH24:MI:SS'::text) AS updated_at,
    to_char(incident.snoozed_until,'YYYY-MM-DD HH24:MM:SS'::text) AS snoozed_until,
    COALESCE(incident.customer_impact, 0) AS customer_impact,
    COALESCE(date_part('epoch'::text, incident.identification_ts - incident.start_ts), 0::double precision) AS acknowledgetime,
    COALESCE(date_part('epoch'::text, incident.end_ts - incident.identification_ts), 0::double precision) AS solutiontime,
    COALESCE(date_part('epoch'::text, incident.end_ts - incident.start_ts), 0::double precision) AS downtime
   FROM incident
  WHERE incident.start_ts IS NOT NULL AND incident.end_ts IS NOT NULL AND incident.identification_ts IS NOT NULL AND incident.end_ts::date >= '2020-01-01'::date
  ORDER BY (incident.end_ts::date);

-- Permissions
ALTER TABLE public.metrics OWNER TO hellper;
GRANT ALL ON TABLE public.metrics TO hellper;
