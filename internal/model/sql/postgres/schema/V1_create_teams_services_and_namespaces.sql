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