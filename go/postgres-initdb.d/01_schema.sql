--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg110+1)
-- Dumped by pg_dump version 16.9 (Debian 16.9-1.pgdg110+1)

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

--
-- Name: mahking_local; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA IF NOT EXISTS mahking_local;
SET search_path TO mahking_local;


ALTER SCHEMA mahking_local OWNER TO postgres;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: mahking_local; Owner: postgres
--

CREATE FUNCTION update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;


ALTER FUNCTION update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: atlas_schema_revisions; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE atlas_schema_revisions (
    version character varying NOT NULL,
    description character varying NOT NULL,
    type bigint DEFAULT 2 NOT NULL,
    applied bigint DEFAULT 0 NOT NULL,
    total bigint DEFAULT 0 NOT NULL,
    executed_at timestamp with time zone NOT NULL,
    execution_time bigint NOT NULL,
    error text,
    error_stmt text,
    hash character varying NOT NULL,
    partial_hashes jsonb,
    operator_version character varying NOT NULL
);


ALTER TABLE atlas_schema_revisions OWNER TO postgres;

--
-- Name: game_rules; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE game_rules (
    id bigint NOT NULL,
    game_id bigint NOT NULL,
    group_id bigint NOT NULL,
    mahjong_type integer NOT NULL,
    initial_points integer NOT NULL,
    return_points integer NOT NULL,
    ranking_points_first integer NOT NULL,
    ranking_points_second integer NOT NULL,
    ranking_points_third integer NOT NULL,
    ranking_points_fourth integer,
    fractional_calculation integer NOT NULL,
    use_bust boolean DEFAULT false NOT NULL,
    bust_point integer,
    use_chip boolean DEFAULT false NOT NULL,
    chip_point integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_game_rules_bust_point CHECK ((((use_bust = true) AND (bust_point IS NOT NULL) AND (bust_point > 0)) OR ((use_bust = false) AND (bust_point IS NULL)))),
    CONSTRAINT chk_game_rules_chip_point CHECK ((((use_chip = true) AND (chip_point IS NOT NULL) AND (chip_point > 0)) OR ((use_chip = false) AND (chip_point IS NULL)))),
    CONSTRAINT chk_game_rules_ranking_points CHECK (
CASE
    WHEN (mahjong_type = 1) THEN ((ranking_points_fourth IS NULL) AND (((ranking_points_first + ranking_points_second) + ranking_points_third) = 0))
    WHEN (mahjong_type = 2) THEN ((ranking_points_fourth IS NOT NULL) AND ((((ranking_points_first + ranking_points_second) + ranking_points_third) + ranking_points_fourth) = 0))
    ELSE NULL::boolean
END),
    CONSTRAINT game_rules_fractional_calculation_check CHECK ((fractional_calculation = ANY (ARRAY[1, 2, 3, 4, 5]))),
    CONSTRAINT game_rules_initial_points_check CHECK ((initial_points > 0)),
    CONSTRAINT game_rules_mahjong_type_check CHECK ((mahjong_type = ANY (ARRAY[1, 2]))),
    CONSTRAINT game_rules_return_points_check CHECK ((return_points > 0))
);


ALTER TABLE game_rules OWNER TO postgres;

--
-- Name: COLUMN game_rules.mahjong_type; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.mahjong_type IS '1: 三麻, 2: 四麻';


--
-- Name: COLUMN game_rules.initial_points; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.initial_points IS '持ち点 (単位: 千点)';


--
-- Name: COLUMN game_rules.return_points; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.return_points IS '返し点 (単位: 千点)';


--
-- Name: COLUMN game_rules.ranking_points_first; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.ranking_points_first IS '1位のウマ';


--
-- Name: COLUMN game_rules.ranking_points_second; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.ranking_points_second IS '2位のウマ';


--
-- Name: COLUMN game_rules.ranking_points_third; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.ranking_points_third IS '3位のウマ';


--
-- Name: COLUMN game_rules.ranking_points_fourth; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.ranking_points_fourth IS '4位のウマ';


--
-- Name: COLUMN game_rules.fractional_calculation; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.fractional_calculation IS '1: 切り上げ, 2: 切り捨て, 3: 四捨五入, 4: 10点未満切り上げ, 5: 10点未満切り捨て';


--
-- Name: COLUMN game_rules.use_bust; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.use_bust IS '飛び設定';


--
-- Name: COLUMN game_rules.bust_point; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.bust_point IS '飛び賞のポイント';


--
-- Name: COLUMN game_rules.use_chip; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.use_chip IS 'チップ設定';


--
-- Name: COLUMN game_rules.chip_point; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_rules.chip_point IS 'チップのポイント';


--
-- Name: game_rules_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE game_rules ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME game_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: game_scores; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE game_scores (
    id bigint NOT NULL,
    game_id bigint NOT NULL,
    group_id bigint NOT NULL,
    member_id bigint NOT NULL,
    seat integer NOT NULL,
    ranking integer NOT NULL,
    raw_score integer NOT NULL,
    point numeric(10,1) NOT NULL,
    chip_count integer,
    is_busted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT game_scores_ranking_check CHECK (((ranking >= 1) AND (ranking <= 4))),
    CONSTRAINT game_scores_seat_check CHECK (((seat >= 1) AND (seat <= 4)))
);


ALTER TABLE game_scores OWNER TO postgres;

--
-- Name: COLUMN game_scores.seat; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.seat IS '1: 東, 2: 南, 3: 西, 4: 北';


--
-- Name: COLUMN game_scores.ranking; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.ranking IS '順位';


--
-- Name: COLUMN game_scores.raw_score; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.raw_score IS '素点 (100点単位, 例: 32400)';


--
-- Name: COLUMN game_scores.point; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.point IS '計算後のポイント';


--
-- Name: COLUMN game_scores.chip_count; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.chip_count IS 'チップ枚数';


--
-- Name: COLUMN game_scores.is_busted; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN game_scores.is_busted IS '飛びフラグ';


--
-- Name: game_scores_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE game_scores ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME game_scores_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: games; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE games (
    id bigint NOT NULL,
    group_id bigint NOT NULL,
    note text,
    played_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE games OWNER TO postgres;

--
-- Name: COLUMN games.note; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN games.note IS '対局メモ';


--
-- Name: COLUMN games.played_at; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN games.played_at IS '対局日時';


--
-- Name: games_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE games ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME games_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: groups; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE groups (
    id bigint NOT NULL,
    uid uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE groups OWNER TO postgres;

--
-- Name: groups_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE groups ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: members; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE members (
    id bigint NOT NULL,
    group_id bigint NOT NULL,
    name character varying(10) DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE members OWNER TO postgres;

--
-- Name: members_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE members ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: rules; Type: TABLE; Schema: mahking_local; Owner: postgres
--

CREATE TABLE rules (
    id bigint NOT NULL,
    group_id bigint NOT NULL,
    mahjong_type integer DEFAULT 2 NOT NULL,
    initial_points integer DEFAULT 25 NOT NULL,
    return_points integer DEFAULT 30 NOT NULL,
    ranking_points_first integer DEFAULT 20 NOT NULL,
    ranking_points_second integer DEFAULT 10 NOT NULL,
    ranking_points_third integer DEFAULT '-10'::integer NOT NULL,
    ranking_points_fourth integer,
    fractional_calculation integer DEFAULT 0 NOT NULL,
    use_bust boolean DEFAULT false NOT NULL,
    bust_point integer,
    use_chip boolean DEFAULT false NOT NULL,
    chip_point integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_bust_point_if_use_bust CHECK ((((use_bust = true) AND (bust_point IS NOT NULL) AND (bust_point > 0)) OR ((use_bust = false) AND (bust_point IS NULL)))),
    CONSTRAINT chk_chip_point_if_use_chip CHECK ((((use_chip = true) AND (chip_point IS NOT NULL) AND (chip_point > 0)) OR ((use_chip = false) AND (chip_point IS NULL)))),
    CONSTRAINT chk_ranking_points_by_mahjong_type CHECK (
CASE
    WHEN (mahjong_type = 1) THEN ((ranking_points_fourth IS NULL) AND (((ranking_points_first + ranking_points_second) + ranking_points_third) = 0))
    WHEN (mahjong_type = 2) THEN ((ranking_points_fourth IS NOT NULL) AND ((((ranking_points_first + ranking_points_second) + ranking_points_third) + ranking_points_fourth) = 0))
    ELSE NULL::boolean
END),
    CONSTRAINT rules_fractional_calculation_check CHECK ((fractional_calculation = ANY (ARRAY[1, 2, 3, 4, 5]))),
    CONSTRAINT rules_initial_points_check CHECK ((initial_points > 0)),
    CONSTRAINT rules_mahjong_type_check CHECK ((mahjong_type = ANY (ARRAY[1, 2]))),
    CONSTRAINT rules_return_points_check CHECK ((return_points > 0))
);


ALTER TABLE rules OWNER TO postgres;

--
-- Name: COLUMN rules.mahjong_type; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.mahjong_type IS '1: 三麻, 2: 四麻';


--
-- Name: COLUMN rules.initial_points; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.initial_points IS '持ち点 (単位: 千点)';


--
-- Name: COLUMN rules.return_points; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.return_points IS '返し点 (単位: 千点)';


--
-- Name: COLUMN rules.ranking_points_first; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.ranking_points_first IS '1位のウマ';


--
-- Name: COLUMN rules.ranking_points_second; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.ranking_points_second IS '2位のウマ';


--
-- Name: COLUMN rules.ranking_points_third; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.ranking_points_third IS '3位のウマ';


--
-- Name: COLUMN rules.ranking_points_fourth; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.ranking_points_fourth IS '4位のウマ';


--
-- Name: COLUMN rules.fractional_calculation; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.fractional_calculation IS '1: 切り上げ, 2: 切り捨て, 3: 四捨五入, 4: 10点未満切り上げ, 5: 10点未満切り捨て';


--
-- Name: COLUMN rules.use_bust; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.use_bust IS '飛び設定';


--
-- Name: COLUMN rules.bust_point; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.bust_point IS '飛び賞のポイント';


--
-- Name: COLUMN rules.use_chip; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.use_chip IS 'チップ設定';


--
-- Name: COLUMN rules.chip_point; Type: COMMENT; Schema: mahking_local; Owner: postgres
--

COMMENT ON COLUMN rules.chip_point IS 'チップのポイント';


--
-- Name: rules_id_seq; Type: SEQUENCE; Schema: mahking_local; Owner: postgres
--

ALTER TABLE rules ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: atlas_schema_revisions atlas_schema_revisions_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY atlas_schema_revisions
    ADD CONSTRAINT atlas_schema_revisions_pkey PRIMARY KEY (version);


--
-- Name: game_rules game_rules_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_rules
    ADD CONSTRAINT game_rules_pkey PRIMARY KEY (id);


--
-- Name: game_scores game_scores_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT game_scores_pkey PRIMARY KEY (id);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY games
    ADD CONSTRAINT games_pkey PRIMARY KEY (id);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);


--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: rules rules_pkey; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY rules
    ADD CONSTRAINT rules_pkey PRIMARY KEY (id);


--
-- Name: game_rules uk_game_rules_game_id; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_rules
    ADD CONSTRAINT uk_game_rules_game_id UNIQUE (game_id);


--
-- Name: game_scores uk_game_scores_game_member; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT uk_game_scores_game_member UNIQUE (game_id, member_id);


--
-- Name: game_scores uk_game_scores_game_ranking; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT uk_game_scores_game_ranking UNIQUE (game_id, ranking);


--
-- Name: game_scores uk_game_scores_game_seat; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT uk_game_scores_game_seat UNIQUE (game_id, seat);


--
-- Name: groups uk_groups_uid; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY groups
    ADD CONSTRAINT uk_groups_uid UNIQUE (uid);


--
-- Name: members uk_members_group_id_name; Type: CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT uk_members_group_id_name UNIQUE (group_id, name);


--
-- Name: idx_game_rules_group_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_game_rules_group_id ON game_rules USING btree (group_id);


--
-- Name: idx_game_scores_group_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_game_scores_group_id ON game_scores USING btree (group_id);


--
-- Name: idx_game_scores_member_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_game_scores_member_id ON game_scores USING btree (member_id);


--
-- Name: idx_games_group_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_games_group_id ON games USING btree (group_id);


--
-- Name: idx_members_group_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_members_group_id ON members USING btree (group_id);


--
-- Name: idx_rules_group_id; Type: INDEX; Schema: mahking_local; Owner: postgres
--

CREATE INDEX idx_rules_group_id ON rules USING btree (group_id);


--
-- Name: game_rules on_update_time_game_rules; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_game_rules BEFORE UPDATE ON game_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: game_scores on_update_time_game_scores; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_game_scores BEFORE UPDATE ON game_scores FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: games on_update_time_games; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_games BEFORE UPDATE ON games FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: groups on_update_time_groups; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_groups BEFORE UPDATE ON groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: members on_update_time_members; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_members BEFORE UPDATE ON members FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: rules on_update_time_rules; Type: TRIGGER; Schema: mahking_local; Owner: postgres
--

CREATE TRIGGER on_update_time_rules BEFORE UPDATE ON rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: game_rules fk_game_rules_game_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_rules
    ADD CONSTRAINT fk_game_rules_game_id FOREIGN KEY (game_id) REFERENCES games(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: game_rules fk_game_rules_group_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_rules
    ADD CONSTRAINT fk_game_rules_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: game_scores fk_game_scores_game_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT fk_game_scores_game_id FOREIGN KEY (game_id) REFERENCES games(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: game_scores fk_game_scores_group_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT fk_game_scores_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: game_scores fk_game_scores_member_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY game_scores
    ADD CONSTRAINT fk_game_scores_member_id FOREIGN KEY (member_id) REFERENCES members(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: games fk_games_group_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY games
    ADD CONSTRAINT fk_games_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: members fk_members_group_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT fk_members_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rules fk_members_group_id; Type: FK CONSTRAINT; Schema: mahking_local; Owner: postgres
--

ALTER TABLE ONLY rules
    ADD CONSTRAINT fk_members_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: game_rules; Type: ROW SECURITY; Schema: mahking_local; Owner: postgres
--

ALTER TABLE game_rules ENABLE ROW LEVEL SECURITY;

--
-- Name: game_scores; Type: ROW SECURITY; Schema: mahking_local; Owner: postgres
--

ALTER TABLE game_scores ENABLE ROW LEVEL SECURITY;

--
-- Name: games; Type: ROW SECURITY; Schema: mahking_local; Owner: postgres
--

ALTER TABLE games ENABLE ROW LEVEL SECURITY;

--
-- Name: game_rules group_policy; Type: POLICY; Schema: mahking_local; Owner: postgres
--

CREATE POLICY group_policy ON game_rules USING ((group_id = (current_setting('app.group_id'::text))::integer));


--
-- Name: game_scores group_policy; Type: POLICY; Schema: mahking_local; Owner: postgres
--

CREATE POLICY group_policy ON game_scores USING ((group_id = (current_setting('app.group_id'::text))::integer));


--
-- Name: games group_policy; Type: POLICY; Schema: mahking_local; Owner: postgres
--

CREATE POLICY group_policy ON games USING ((group_id = (current_setting('app.group_id'::text))::integer));


--
-- Name: members group_policy; Type: POLICY; Schema: mahking_local; Owner: postgres
--

CREATE POLICY group_policy ON members USING ((group_id = (current_setting('app.group_id'::text))::integer));


--
-- Name: rules group_policy; Type: POLICY; Schema: mahking_local; Owner: postgres
--

CREATE POLICY group_policy ON rules USING ((group_id = (current_setting('app.group_id'::text))::integer));


--
-- Name: members; Type: ROW SECURITY; Schema: mahking_local; Owner: postgres
--

ALTER TABLE members ENABLE ROW LEVEL SECURITY;

--
-- Name: rules; Type: ROW SECURITY; Schema: mahking_local; Owner: postgres
--

ALTER TABLE rules ENABLE ROW LEVEL SECURITY;

--
-- PostgreSQL database dump complete
--

