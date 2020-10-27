--
-- PostgreSQL database dump
--

-- Dumped from database version 13.0
-- Dumped by pg_dump version 13.0

-- Started on 2020-10-27 01:35:53

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;


CREATE DATABASE rss;
ALTER DATABASE rss OWNER TO postgres;

\connect rss

SET default_tablespace = '';

SET default_table_access_method = heap;


--
-- TOC entry 201 (class 1259 OID 24579)
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id integer NOT NULL,
    site_id integer,
    title text,
    link text,
    description text,
    pub_date timestamp with time zone,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- TOC entry 200 (class 1259 OID 24577)
-- Name: article_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.article_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.article_id_seq OWNER TO postgres;

--
-- TOC entry 3024 (class 0 OID 0)
-- Dependencies: 200
-- Name: article_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.article_id_seq OWNED BY public.articles.id;


--
-- TOC entry 207 (class 1259 OID 24612)
-- Name: site_info; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.site_info (
    id integer NOT NULL,
    site_id integer,
    title text,
    link text,
    description text,
    language text,
    id_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.site_info OWNER TO postgres;

--
-- TOC entry 206 (class 1259 OID 24610)
-- Name: site_info_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.site_info_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.site_info_id_seq OWNER TO postgres;

--
-- TOC entry 3025 (class 0 OID 0)
-- Dependencies: 206
-- Name: site_info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.site_info_id_seq OWNED BY public.site_info.id;


--
-- TOC entry 205 (class 1259 OID 24601)
-- Name: site_rules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.site_rules (
    id integer NOT NULL,
    site_id integer,
    article_root_name text,
    title text,
    url text,
    description text,
    pub_date text,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.site_rules OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 24599)
-- Name: site_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.site_rules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.site_rules_id_seq OWNER TO postgres;

--
-- TOC entry 3026 (class 0 OID 0)
-- Dependencies: 204
-- Name: site_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.site_rules_id_seq OWNED BY public.site_rules.id;


--
-- TOC entry 203 (class 1259 OID 24590)
-- Name: sites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sites (
    id integer NOT NULL,
    link text,
    cron text,
    task_id integer,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sites OWNER TO postgres;

--
-- TOC entry 202 (class 1259 OID 24588)
-- Name: sites_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sites_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sites_id_seq OWNER TO postgres;

--
-- TOC entry 3027 (class 0 OID 0)
-- Dependencies: 202
-- Name: sites_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sites_id_seq OWNED BY public.sites.id;


--
-- TOC entry 2872 (class 2604 OID 24582)
-- Name: articles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.article_id_seq'::regclass);


--
-- TOC entry 2878 (class 2604 OID 24615)
-- Name: site_info id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.site_info ALTER COLUMN id SET DEFAULT nextval('public.site_info_id_seq'::regclass);


--
-- TOC entry 2876 (class 2604 OID 24604)
-- Name: site_rules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.site_rules ALTER COLUMN id SET DEFAULT nextval('public.site_rules_id_seq'::regclass);


--
-- TOC entry 2874 (class 2604 OID 24593)
-- Name: sites id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sites ALTER COLUMN id SET DEFAULT nextval('public.sites_id_seq'::regclass);


--
-- TOC entry 2881 (class 2606 OID 24587)
-- Name: articles article_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT article_pkey PRIMARY KEY (id);


--
-- TOC entry 2888 (class 2606 OID 24620)
-- Name: site_info site_info_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.site_info
    ADD CONSTRAINT site_info_pkey PRIMARY KEY (id);


--
-- TOC entry 2886 (class 2606 OID 24609)
-- Name: site_rules site_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.site_rules
    ADD CONSTRAINT site_rules_pkey PRIMARY KEY (id);


--
-- TOC entry 2884 (class 2606 OID 24598)
-- Name: sites sites_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sites
    ADD CONSTRAINT sites_pkey PRIMARY KEY (id);


--
-- TOC entry 2882 (class 1259 OID 24678)
-- Name: article_title_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX article_title_idx ON public.articles USING btree (title);


-- Completed on 2020-10-27 01:35:54

--
-- PostgreSQL database dump complete
--

