
BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";	

CREATE TABLE IF NOT EXISTS "users" (
    user_id character varying(45) NOT NULL DEFAULT uuid_generate_v4(),
    user_name character varying(200) NOT NULL UNIQUE,
    email character varying(200) NOT NULL UNIQUE,
    password text NOT NULL,
    user_status character varying(45) NOT NULL DEFAULT 'email_not_verified',
    user_picture text NOT NULL DEFAULT 'http://localhost:7070/images/default_image.png',
    user_role_id int,
    user_created_at timestamp with time zone NOT NULL DEFAULT now(),
    user_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id),
    CONSTRAINT user_status_check CHECK (user_status IN ('active', 'inactive', 'email_not_verified'))
   
) 
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "roles" (
    role_id int PRIMARY KEY,
    role_name character varying(200) NOT NULL UNIQUE,
    role_status character varying(45) NOT NULL DEFAULT 'active',
    role_created_at timestamp with time zone NOT NULL DEFAULT now(),
    role_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT role_status_check CHECK (role_status IN ('active', 'inactive'))
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "posts" (
    post_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    post_title character varying(200) NOT NULL,
    post_content text NOT NULL,
    post_image text NOT NULL DEFAULT 'http://localhost:7070/images/default_image.png',
    post_author_id character varying(45) NOT NULL,
    post_status character varying(45) NOT NULL DEFAULT 'active',
    post_created_at timestamp with time zone NOT NULL DEFAULT now(),
    post_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT post_status_check CHECK (post_status IN ('active', 'inactive'))
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "tags_posts" (
    tags_posts_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    tags_posts_tag_id int NOT NULL,
    tags_posts_post_id character varying(45) NOT NULL
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "tags" (
    tag_id int PRIMARY KEY,
    tag_name character varying(200) NOT NULL UNIQUE,
    tag_status character varying(45) NOT NULL DEFAULT 'active',
    tag_created_at timestamp with time zone NOT NULL DEFAULT now(),
    tag_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT tag_status_check CHECK (tag_status IN ('active', 'inactive'))
)
WITH (
    OIDS = FALSE
);


CREATE TABLE IF NOT EXISTS "comments" (
    comment_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    comment_content text NOT NULL,
    comment_post_id character varying(45) NOT NULL,
    comment_user_id character varying(45) NOT NULL,
    comment_type character varying(45) NOT NULL,
    comment_status character varying(45) NOT NULL DEFAULT 'active',
    comment_created_at timestamp with time zone NOT NULL DEFAULT now(),
    comment_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT comment_type_check CHECK (comment_type IN ('post', 'comment')),
    CONSTRAINT comment_status_check CHECK (comment_status IN ('active', 'inactive', 'blocked'))
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "parent_child_comments" (
    parent_child_comments_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    parent_child_comments_parent_id character varying(45) NOT NULL,
    parent_child_comments_child_id character varying(45) NOT NULL,
    CONSTRAINT parent_child_comments_parent_id_unique UNIQUE (parent_child_comments_parent_id, parent_child_comments_child_id)
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "email_confirms" (
    email_confirms_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    email_confirms_email character varying(200) NOT NULL,
    email_confirms_status character varying(45) NOT NULL DEFAULT 'active',
    email_confirms_created_at timestamp with time zone NOT NULL DEFAULT now(),
    email_confirms_updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT email_confirms_status_check CHECK (email_confirms_status IN ('active', 'inactive'))
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS "contacts" (
    contacts_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    contacts_name character varying(200) NOT NULL,
    contacts_email character varying(200) NOT NULL,
    contacts_message text NOT NULL
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS "users"
    ADD FOREIGN KEY (user_role_id) 
    REFERENCES roles (role_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "posts"
    ADD FOREIGN KEY (post_author_id) 
    REFERENCES users (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "comments"
    ADD FOREIGN KEY (comment_post_id) 
    REFERENCES posts (post_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "comments"
    ADD FOREIGN KEY (comment_user_id) 
    REFERENCES users (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "parent_child_comments"
    ADD FOREIGN KEY (parent_child_comments_parent_id) 
    REFERENCES comments (comment_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "parent_child_comments"
    ADD FOREIGN KEY (parent_child_comments_child_id) 
    REFERENCES comments (comment_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "tags_posts"
    ADD FOREIGN KEY (tags_posts_tag_id) 
    REFERENCES tags (tag_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
ALTER TABLE IF EXISTS "tags_posts"
    ADD FOREIGN KEY (tags_posts_post_id) 
    REFERENCES posts (post_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;


INSERT INTO roles(role_id, role_name, role_status) VALUES(1,'admin', 'active');
INSERT INTO roles(role_id, role_name, role_status) VALUES(2,'publisher', 'active');
INSERT INTO roles(role_id, role_name, role_status) VALUES(3,'general', 'active');

INSERT INTO tags(tag_id, tag_name) VALUES(1,'javascript');
INSERT INTO tags(tag_id, tag_name) VALUES(2,'svelte');
INSERT INTO tags(tag_id, tag_name) VALUES(3,'react');
INSERT INTO tags(tag_id, tag_name) VALUES(4,'vue');
INSERT INTO tags(tag_id, tag_name) VALUES(5,'angular');
INSERT INTO tags(tag_id, tag_name) VALUES(6,'nodejs');
INSERT INTO tags(tag_id, tag_name) VALUES(7,'php');
INSERT INTO tags(tag_id, tag_name) VALUES(8,'python');
INSERT INTO tags(tag_id, tag_name) VALUES(9,'ruby');
INSERT INTO tags(tag_id, tag_name) VALUES(10,'java');
INSERT INTO tags(tag_id, tag_name) VALUES(11,'c');
INSERT INTO tags(tag_id, tag_name) VALUES(12,'c++');
INSERT INTO tags(tag_id, tag_name) VALUES(13,'c#');
INSERT INTO tags(tag_id, tag_name) VALUES(14,'swift');
INSERT INTO tags(tag_id, tag_name) VALUES(15,'kotlin');
INSERT INTO tags(tag_id, tag_name) VALUES(16,'go');
INSERT INTO tags(tag_id, tag_name) VALUES(17,'rust');
INSERT INTO tags(tag_id, tag_name) VALUES(18,'typescript');
INSERT INTO tags(tag_id, tag_name) VALUES(19,'unity');
INSERT INTO tags(tag_id, tag_name) VALUES(20,'css');
INSERT INTO tags(tag_id, tag_name) VALUES(21,'html');
INSERT INTO tags(tag_id, tag_name) VALUES(22,'sql');
INSERT INTO tags(tag_id, tag_name) VALUES(23,'mysql');
INSERT INTO tags(tag_id, tag_name) VALUES(24,'mongodb');
INSERT INTO tags(tag_id, tag_name) VALUES(25,'postgresql');
INSERT INTO tags(tag_id, tag_name) VALUES(26,'redis');
INSERT INTO tags(tag_id, tag_name) VALUES(27,'docker');
INSERT INTO tags(tag_id, tag_name) VALUES(28,'sass');

END;