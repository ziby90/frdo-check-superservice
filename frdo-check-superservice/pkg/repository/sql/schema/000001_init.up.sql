CREATE TABLE cls_queue_status
(
    id   bigserial NOT NULL,
    name varchar(100),
    code varchar(100),
    CONSTRAINT cls_queue_status_pkey PRIMARY KEY (id)
);

-- Auto-generated SQL script #202102111138
INSERT INTO public.cls_queue_status (id, "name", code)
VALUES (1, 'документ по реквизитам найден', 'not_found');
INSERT INTO public.cls_queue_status (id, "name", code)
VALUES (2, 'документ по реквизитам не найден', 'founded');
INSERT INTO public.cls_queue_status (id, "name", code)
VALUES (3, 'необходимо дополнить запрос дополнительными данными', 'many');


CREATE TABLE queue_in
(
    id             bigserial NOT NULL,
    created_at     timestamptz NULL default now(),
    updated_at     timestamptz NULL,
    deleted_at     timestamptz NULL,
    package_id     int8,
    smev_id        int8,
    message_id     varchar(255),
    kind_action_id int4,
    type           varchar(100),
    xml            text,
    error          text,
    status         int4,
    CONSTRAINT queue_in_pkey PRIMARY KEY (id)
);

CREATE TABLE queue_out
(
    id                 bigserial NOT NULL,
    created_at         timestamptz NULL default now(),
    updated_at         timestamptz NULL,
    deleted_at         timestamptz NULL,
    package_id         int8,
    smev_id            int8,
    message_id         varchar(255),
    kind_action_id     int4,
    type               varchar(100),
    xml                text,
    error              text,
    status             int4,
    status_response_id int4,
    CONSTRAINT queue_out_pkey PRIMARY KEY (id),
    CONSTRAINT status_response_id_foreign FOREIGN KEY (status_response_id) REFERENCES cls_queue_status (id) ON DELETE RESTRICT
);
