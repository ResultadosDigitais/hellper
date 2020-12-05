CREATE TABLE IF NOT EXISTS public.team (
    id serial NOT null primary key,
    name varchar(256) NOT NULL,
    created_at timestamp not null DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.person (
    email text NOT NULL primary key,
    slack_member_id text NOT NULL,
    created_at timestamp not null DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.team_member (
    person_email text NOT NULL REFERENCES public.person (email),
    team_id integer NOT NULL REFERENCES public.team (id),
    created_at timestamp not null DEFAULT now(),

    primary key(person_email, team_id)
);

CREATE TABLE IF NOT EXISTS public.service (
    id serial NOT null primary key,
    name varchar(256) NOT NULL,
    created_at timestamp not null DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.service_instance (
    id serial NOT null primary key,
    name varchar(512) NOT NULL,
    service_id integer NOT NULL REFERENCES public.service (id),
    owner_team_id integer NOT NULL REFERENCES public.team (id),

    created_at timestamp not null DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.service_instance_stakeholder (
    service_instance_id integer NOT NULL REFERENCES public.service_instance (id),
    team_id integer NOT NULL REFERENCES public.team (id),
    created_at timestamp not null DEFAULT now(),

    primary key(service_instance_id, team_id)
);

ALTER TABLE public.team OWNER TO hellper;
ALTER TABLE public.person OWNER TO hellper;
ALTER TABLE public.team_member OWNER TO hellper;
ALTER TABLE public.service OWNER TO hellper;
ALTER TABLE public.service_instance OWNER TO hellper;
ALTER TABLE public.service_instance_stakeholder OWNER TO hellper;

-- TODO: Remove this from here or replace with a more complete bootstrap
insert into service (name, created_at) values ('Maestro', now());
insert into service (name, created_at) values ('Matchmaker', now());
insert into service (name, created_at) values ('Remote Config', now());
insert into team (name, created_at) values ('Game Services & Analytics - MIRC Squad', now());
insert into team (name, created_at) values ('Game Services & Analytics - PAC Squad', now());
insert into team (name, created_at) values ('Tennis Clash', now());
insert into team (name, created_at) values ('Zooba', now());
insert into service_instance (name, service_id, owner_team_id, created_at) values ('Zooba', 1, 1, now());
insert into service_instance (name, service_id, owner_team_id, created_at) values ('Tennis', 1, 1, now());
insert into service_instance (name, service_id, owner_team_id, created_at) values ('Zooba', 2, 1, now());
insert into service_instance (name, service_id, owner_team_id, created_at) values ('Tennis', 2, 1, now());
