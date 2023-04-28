database-schema
---------------


-- illa_supervisor-database-design


-- init
create database illa_supervisor;
\c illa_supervisor;
create user illa_supervisor with encrypted password 'illa2022';
grant all privileges on database illa_supervisor to illa_supervisor;
CREATE EXTENSION pg_trgm;
CREATE EXTENSION btree_gin;


/**
 * TEAM Management
 *
 *
 */

-- teams
create table if not exists teams (
    id                       bigserial                               not null primary key,
    uid                      uuid         default gen_random_uuid()  not null,
    name                     varchar(255)                            not null,
    identifier               varchar(255) unique                     not null,
    icon                     varchar(255)                            not null,
    permission               jsonb                                   not null,
    created_at               timestamp                               not null,
    updated_at               timestamp                               not null,
    constraint               teams_ukey unique (id, uid)
);

CREATE INDEX teams_uid ON teams (uid);

alter table
    teams owner to illa_supervisor;

-- users
create table if not exists users (
    id                       bigserial                         not null primary key,
    uid                      uuid    default gen_random_uuid() not null,
    nickname                 varchar(15)                       not null, 
    password_digest          varchar(60)                       not null,
    email                    varchar(255)                      not null,
    avatar                   varchar(255)                      not null,
    sso_config               jsonb                             not null, 
    customization            jsonb                             not null, 
    created_at               timestamp                         not null,
    updated_at               timestamp                         not null,
    constraint               users_ukey unique (id, uid)
);

CREATE INDEX users_uid ON users (uid);
CREATE INDEX users_nickname_fulltext ON users USING gin (to_tsvector('english', nickname));
CREATE INDEX users_email_fulltext ON users USING gin (to_tsvector('english', email));

alter table
    users owner to illa_supervisor;

-- team_members
create table if not exists team_members (
    id                       bigserial                         not null primary key,
    team_id                  bigserial                         not null,
    user_id                  bigserial                         not null,  
    user_role                smallint                          not null, 
    permission               jsonb                            ,         
    status                   smallint                          not null, 
    created_at               timestamp                         not null,
    updated_at               timestamp                         not null
);

CREATE INDEX team_members_team_and_user_id ON team_members (team_id, user_id);

alter table
    team_members owner to illa_supervisor;

-- invites
create table if not exists invites (
    id                       bigserial                            not null primary key,
    uid                      uuid       default gen_random_uuid() not null, 
    category                 smallint                             not null,  
    team_id                  bigserial                            not null,
    team_member_id           bigserial                            not null,  
    email                    varchar(255)                        ,          
    email_status             boolean default false                not null,  
    user_role                smallint                             not null,  
    status                   smallint                             not null,  
    created_at               timestamp                            not null,
    updated_at               timestamp                            not null,
    constraint               invite_ukey unique (id, uid)
);

CREATE INDEX invites_uid ON invites (uid);
CREATE INDEX invites_email ON invites (email);
CREATE INDEX invites_user_role ON invites (user_role);

alter table
    invites owner to illa_supervisor;


/**
 * Role Management
 *
 *
 */

-- roles
create table if not exists roles (
    id                       bigserial                            not null primary key,
    uid                      uuid       default gen_random_uuid() not null,
    name                     varchar(255)                         not null,         
    team_id                  bigserial                            not null, 
    permissions              jsonb                                not null,
    created_at               timestamp                            not null,
    updated_at               timestamp                            not null
);
CREATE INDEX roles_id_team_id ON roles(id, team_id);
CREATE INDEX roles_name_fulltext ON roles USING gin (to_tsvector('english', name));
alter table roles owner to illa_supervisor;

-- user_role_relations
create table if not exists user_role_relations (
    id                       bigserial                            not null primary key,
    uid                      uuid       default gen_random_uuid() not null,
    team_id                  bigserial                            not null, 
    role_id                  bigserial                            not null, 
    user_id                  bigserial                            not null,
    created_at               timestamp                            not null,
    updated_at               timestamp                            not null
);
CREATE INDEX user_role_relations_team_role_user_id ON user_role_relations(team_id, role_id, user_id);
alter table user_role_relations owner to illa_supervisor;

-- unit_role_relations
create table if not exists unit_role_relations (
    id                       bigserial                            not null primary key,
    uid                      uuid       default gen_random_uuid() not null,
    team_id                  bigserial                            not null, 
    role_id                  bigserial                            not null, 
    unit_id                  bigserial                            not null,
    unit_type                smallint                             not null,
    created_at               timestamp                            not null,
    updated_at               timestamp                            not null
);
CREATE INDEX unit_role_relations_team_role_unit_id_and_unit_type ON unit_role_relations(team_id, role_id, unit_id, unit_type);
alter table unit_role_relations owner to illa_supervisor;


/**
 * DDL
 *
 *
 */

INSERT INTO teams ( 
    id, uid, name, identifier, icon, permission, created_at, updated_at
) VALUES (
    0, '00000000-0000-0000-0000-000000000000', 'my-team'    , '0'  , 'https://cdn.illacloud.com/email-template/people.png', to_jsonb('{"allowEditorInvite": true, "allowViewerInvite": true, "inviteLinkEnabled": true, "allowEditorManageTeamMember": true, "allowViewerManageTeamMember": true}'::text), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

INSERT INTO users (
    id, uid, nickname, password_digest, email, avatar, sso_config, customization, created_at, updated_at
) VALUES (
    1, '00000000-0000-0000-0000-000000000000', 'root', '$2a$10$iVIxJRgy1K6RIV389AYg3OiMIbuDyuCIja1xrHGkCljdg/6gdmWXa'::text, 'root', '', to_jsonb('{"default": ""}'::text) , to_jsonb('{"Language": "en-US", "IsSubscribed": false}'::text), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

INSERT INTO team_members (
    id, team_id, user_id, user_role, permission, status, created_at, updated_at   
) VALUES (      
    1, 0, 1, 1, to_jsonb('{"Config": 0}'::text), 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);
