--
-- PostgreSQL database dump
--

-- \restrict REwvFs23gpYgUCvfzUnEGLbT4XImpb9RZPSTxCAALPhfIQBgfcBkleFzMixRaiF


-- Dumped from database version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    restaurant_id bigint,
    name character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.categories OWNER TO ilx;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: ilx
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO ilx;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ilx
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: menu; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.menu (
    id bigint NOT NULL,
    category_id bigint NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    price_cent integer NOT NULL,
    is_available boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT check_isavailable_constraint CHECK (((is_available = true) OR (is_available = false)))
);


ALTER TABLE public.menu OWNER TO ilx;

--
-- Name: menu_id_seq; Type: SEQUENCE; Schema: public; Owner: ilx
--

CREATE SEQUENCE public.menu_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.menu_id_seq OWNER TO ilx;

--
-- Name: menu_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ilx
--

ALTER SEQUENCE public.menu_id_seq OWNED BY public.menu.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.permissions (
    id bigint NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.permissions OWNER TO ilx;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: ilx
--

CREATE SEQUENCE public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO ilx;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ilx
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: restaurant; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.restaurant (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    country character varying(50) NOT NULL,
    full_address text NOT NULL,
    cuisine character varying(50) NOT NULL,
    status character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT check_status_constarint CHECK ((((status)::text = 'open'::text) OR ((status)::text = 'closed'::text)))
);


ALTER TABLE public.restaurant OWNER TO ilx;

--
-- Name: restaurant_id_seq; Type: SEQUENCE; Schema: public; Owner: ilx
--

CREATE SEQUENCE public.restaurant_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurant_id_seq OWNER TO ilx;

--
-- Name: restaurant_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ilx
--

ALTER SEQUENCE public.restaurant_id_seq OWNED BY public.restaurant.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.tokens (
    hash bytea NOT NULL,
    user_id bigint,
    expiry timestamp with time zone NOT NULL,
    scope text NOT NULL
);


ALTER TABLE public.tokens OWNER TO ilx;

--
-- Name: users; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    role character varying(20) DEFAULT 'customer'::character varying,
    restaurant_id bigint,
    is_active boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    password_hash text NOT NULL,
    last_updated timestamp with time zone DEFAULT now(),
    CONSTRAINT check_role_constarint CHECK ((((role)::text = 'customer'::text) OR ((role)::text = 'seller'::text) OR ((role)::text = 'admin'::text)))
);


ALTER TABLE public.users OWNER TO ilx;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: ilx
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO ilx;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ilx
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users_permissions; Type: TABLE; Schema: public; Owner: ilx
--

CREATE TABLE public.users_permissions (
    user_id bigint NOT NULL,
    permission_id bigint NOT NULL
);


ALTER TABLE public.users_permissions OWNER TO ilx;

--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: menu id; Type: DEFAULT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.menu ALTER COLUMN id SET DEFAULT nextval('public.menu_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: restaurant id; Type: DEFAULT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.restaurant ALTER COLUMN id SET DEFAULT nextval('public.restaurant_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: menu menu_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT menu_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: restaurant restaurant_name_key; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.restaurant
    ADD CONSTRAINT restaurant_name_key UNIQUE (name);


--
-- Name: restaurant restaurant_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.restaurant
    ADD CONSTRAINT restaurant_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (hash);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: categories categories_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurant(id) ON DELETE CASCADE;


--
-- Name: menu menu_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT menu_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: tokens tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_permissions users_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: users_permissions users_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users users_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ilx
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurant(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

-- \unrestrict REwvFs23gpYgUCvfzUnEGLbT4XImpb9RZPSTxCAALPhfIQBgfcBkleFzMixRaiF

