CREATE FUNCTION diesel_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (
        NEW IS DISTINCT FROM OLD AND
        NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at
    ) THEN
        NEW.updated_at := current_timestamp;
    END IF;
    RETURN NEW;
END;
$$;

CREATE FUNCTION notify_data_changed() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NOTIFY data_changed, '';
    RETURN NULL;
END;
$$;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email character varying NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

CREATE TABLE objects (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL,
    image character varying NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

CREATE TABLE object_users (
    object_id integer NOT NULL,
    user_id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    PRIMARY KEY(object_id, user_id)
);


CREATE TRIGGER notify_data_changed AFTER INSERT OR DELETE OR UPDATE ON users FOR EACH STATEMENT EXECUTE PROCEDURE notify_data_changed();
CREATE TRIGGER notify_data_changed AFTER INSERT OR DELETE OR UPDATE ON objects FOR EACH STATEMENT EXECUTE PROCEDURE notify_data_changed();
CREATE TRIGGER notify_data_changed AFTER INSERT OR DELETE OR UPDATE ON object_users FOR EACH STATEMENT EXECUTE PROCEDURE notify_data_changed();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE diesel_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON objects FOR EACH ROW EXECUTE PROCEDURE diesel_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON object_users FOR EACH ROW EXECUTE PROCEDURE diesel_set_updated_at();
