CREATE TABLE public.comments (
    comment_id uuid NOT NULL,
    post_id uuid NOT NULL,
    comment_on_id uuid,
    subscriber_id uuid NOT NULL,
    content character varying(255) NOT NULL,
    created_at timestamptz DEFAULT(now()) NOT NULL,
    created_by character varying(255) NOT NULL,
    updated_at timestamptz DEFAULT(now()) NOT NULL,
    updated_by character varying(255) 
);

ALTER TABLE public.comments OWNER TO postgres;

--
-- TOC entry 201 (class 1259 OID 32573)
-- Name: post; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    post_id uuid NOT NULL,
    post_category_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    introduction text NULL,
    content text NULL,
    main_image_alt character varying(255)  NULL,
    main_image_path character varying(255)  NULL,
    thumbnail_image_path character varying(255) NULL,
    thumbnail_image_alt character varying(255) NULL,
    author character varying(255) NOT NULL,
    author_image_path character varying(255) NOT NULL,
    author_image_alt character varying(255) NOT NULL,
    published boolean NOT NULL default false,
    published_at timestamptz DEFAULT(now()) NULL,
    published_by character varying(255) NULL,
    created_at timestamptz DEFAULT(now()) NOT NULL,
    created_by character varying(255) NOT NULL,
    updated_at timestamptz DEFAULT(now()) NOT NULL,
    updated_by character varying(255)
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- TOC entry 202 (class 1259 OID 32589)
-- Name: post_detail; Type: TABLE; Schema: public; Owner: postgres
--

/* CREATE TABLE public.posts_detail ( */
/*      post_detail_id uuid NOT NULL, */
/*      post_id uuid NOT NULL, */
/*      line_order integer DEFAULT 0 NOT NULL, */
/*      detail_type character varying(255) NOT NULL, */
/*      param1 character varying(255) NULL, */
/*      param2 character varying(255) NULL, */
/*      param3 character varying(255) NULL, */
/*      param4 character varying(255) NULL, */
/*      content text NULL,  */
/*      created_at timestamptz DEFAULT(now()) NOT NULL, */
/*      created_by character varying(255) NOT NULL, */
/*      updated_at timestamptz DEFAULT(now()) NOT NULL, */
/*      updated_by character varying(255) */
/* ); */


-- ALTER TABLE public.posts_detail OWNER TO postgres; 

--
-- TOC entry 200 (class 1259 OID 32564)
-- Name: post_category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.post_categories (
    post_category_id uuid NOT NULL,
    definition character varying(255) NOT NULL,
    created_at timestamptz DEFAULT(now()) NOT NULL,
    created_by character varying(255) NOT NULL,
    updated_at timestamptz DEFAULT(now()) NOT NULL,
    updated_by character varying(255)
);

ALTER TABLE public.post_categories OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 32604)
-- Name: subscriber; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subscribers (
    subscriber_id uuid NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamptz DEFAULT(now()) NOT NULL,
    created_by character varying(255) NOT NULL,
    updated_at timestamptz DEFAULT(now()) NOT NULL,
    updated_by character varying(255)
);


ALTER TABLE public.subscribers OWNER TO postgres;

--
-- TOC entry 205 (class 1259 OID 32637)
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    email character varying(255) UNIQUE NOT NULL,
    hashed_password varchar NOT NULL,
    password_changed_at timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    profile_image_path character varying(255),
    salt integer DEFAULT 0 NOT NULL,
    last_login timestamptz DEFAULT(now()) NOT NULL,
    created_at timestamptz DEFAULT(now()) NOT NULL,
    created_by character varying(255) NOT NULL,
    updated_at timestamptz DEFAULT(now()) NOT NULL,
    updated_by character varying(255)
);

ALTER TABLE public.users OWNER TO postgres;

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");

--
-- TOC entry 2894 (class 2606 OID 32621)
-- Name: comment comment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comment_pkey PRIMARY KEY (comment_id);

--
-- TOC entry 2890 (class 2606 OID 32598)
-- Name: post_detail post_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

/* ALTER TABLE ONLY public.posts_detail */
/*      ADD CONSTRAINT post_detail_pkey PRIMARY KEY (post_detail_id); */

--
-- TOC entry 2886 (class 2606 OID 32581)
-- Name: post post_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT post_pkey PRIMARY KEY (post_id);


--
-- TOC entry 2888 (class 2606 OID 32583)
-- Name: post post_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT post_slug_key UNIQUE (slug);


--
-- TOC entry 2884 (class 2606 OID 32572)
-- Name: post_category post_category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post_categories
    ADD CONSTRAINT post_category_pkey PRIMARY KEY (post_category_id);


--
-- TOC entry 2892 (class 2606 OID 32612)
-- Name: subscriber subscriber_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscribers
    ADD CONSTRAINT subscriber_pkey PRIMARY KEY (subscriber_id);


--
-- TOC entry 2896 (class 2606 OID 32647)
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT user_pkey PRIMARY KEY (email);


--
-- TOC entry 2900 (class 2606 OID 32627)
-- Name: comment comment_on_comment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comment_on_comment_id_fkey FOREIGN KEY (comment_on_id) REFERENCES public.comments(comment_id);


--
-- TOC entry 2899 (class 2606 OID 32622)
-- Name: comment comment_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comment_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(post_id);


--
-- TOC entry 2901 (class 2606 OID 32632)
-- Name: comment comment_subscriber_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comment_subscriber_id_fkey FOREIGN KEY (subscriber_id) REFERENCES public.subscribers(subscriber_id);


--
-- TOC entry 2898 (class 2606 OID 32599)
-- Name: post_detail post_detail_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

/* ALTER TABLE ONLY public.posts_detail */
/*      ADD CONSTRAINT post_detail_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(post_id);  */


--
-- TOC entry 2897 (class 2606 OID 32584)
-- Name: post post_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT post_category_id_fkey FOREIGN KEY (post_category_id) REFERENCES public.post_categories(post_category_id);


-- Completed on 2022-10-26 17:27:35

--
-- PostgreSQL database dump complete
--

insert into post_categories(post_category_id, definition,created_by) values ('4e86fd8f-c75b-498d-982d-48f7033a3a47', 'C#', 'Francesc PUjol');
insert into post_categories(post_category_id, definition,created_by) values ('5323dd9f-854e-4977-8264-8b901a70cb88', 'Golang', 'Francesc PUjol');