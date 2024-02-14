CREATE TABLE public.posts (
    id integer NOT NULL,
    title character varying(255),
    content character varying(500),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);
--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.posts
ALTER COLUMN id
ADD GENERATED ALWAYS AS IDENTITY (
        SEQUENCE NAME public.posts_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1
    );
insert into public.posts (
        title,
        content,
        created_at,
        updated_at
    )
values(
        'First post',
        'It is first mock post',
        '2024-01-20 00:00:00',
        '2024-01-20 00:00:00'
    );
ALTER TABLE ONLY public.posts
ADD CONSTRAINT posts_pkey PRIMARY KEY (id);