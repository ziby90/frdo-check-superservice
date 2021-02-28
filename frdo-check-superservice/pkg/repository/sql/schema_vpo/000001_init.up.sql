-- recipient with one document
INSERT INTO persons.recipients (id,"name",created,author_id,surname,patronymic,birthday,sex,nationaly_id,snils,birthplace)
VALUES (999999999,'TESTING_FRDO_CHECK_NAME','2021-02-11 13:18:45.128',1317,'testing_frdo_check_surname','testing_frdo_check_patronymic','1980-02-01',1,151,'111-222-333 44 ','Moscow');
INSERT INTO persons.documents ("name",document_type_id,status_id,education_level_id,recipient_id,series,"number",register_number,organization_id,document_package_id,specialty,training_period,admission_year,graduation_year,confirmation_loss,confirmation_exchange,created,author_id,issue_date,specialty_id,qualification_id,code,forms_education_id,first_high_education,source_financing_id,v_doc)
VALUES ('Диплом',17,1,4,999999999,'11111111111','11111111111','11111111111',4757,1136,'Юриспруденция','4',2015,2019,false,false,'2020-10-28 11:34:28.194',1,'2019-06-24 00:00:00.000',44301,34,'40.03.01',1,1,3,1);

-- recipient with many documents
INSERT INTO persons.recipients (id,"name",created,author_id,surname,patronymic,birthday,sex,nationaly_id,snils,birthplace)
VALUES (999999998,'TESTING_FRDO_CHECK_NAME2','2021-02-11 13:18:45.128',1317,'testing_frdo_check_surname2','testing_frdo_check_patronymic2','1981-02-01',1,151,'112-222-333 44 ','Moscow');
INSERT INTO persons.documents ("name",document_type_id,status_id,education_level_id,recipient_id,series,"number",register_number,organization_id,document_package_id,specialty,training_period,admission_year,graduation_year,confirmation_loss,confirmation_exchange,created,author_id,issue_date,specialty_id,qualification_id,code,forms_education_id,first_high_education,source_financing_id,v_doc)
VALUES ('Диплом',17,1,4,999999998,'22222222222','22222222222','22222222222',4757,1136,'Юриспруденция','4',2015,2019,false,false,'2020-10-28 11:34:28.194',1,'2019-06-24 00:00:00.000',44301,34,'40.03.01',1,1,3,1);
INSERT INTO persons.documents ("name",document_type_id,status_id,education_level_id,recipient_id,series,"number",register_number,organization_id,document_package_id,specialty,training_period,admission_year,graduation_year,confirmation_loss,confirmation_exchange,created,author_id,issue_date,specialty_id,qualification_id,code,forms_education_id,first_high_education,source_financing_id,v_doc)
VALUES ('Диплом',17,1,4,999999998,'11111111111','11111111111','11111111111',4757,1136,'Юриспруденция','4',2015,2019,false,false,'2020-10-28 11:34:28.194',1,'2019-06-24 00:00:00.000',44301,34,'40.03.01',1,1,3,1);
INSERT INTO persons.documents ("name",document_type_id,status_id,education_level_id,recipient_id,series,"number",register_number,organization_id,document_package_id,specialty,training_period,admission_year,graduation_year,confirmation_loss,confirmation_exchange,created,author_id,issue_date,specialty_id,qualification_id,code,forms_education_id,first_high_education,source_financing_id,v_doc)
VALUES ('Диплом',17,1,4,999999998,'33333333333','33333333333','33333333333',4757,1136,'Юриспруденция','4',2015,2019,false,false,'2020-10-28 11:34:28.194',1,'2019-06-24 00:00:00.000',44301,34,'40.03.01',1,1,3,1);

-- recipient without documents
INSERT INTO persons.recipients (id,"name",created,author_id,surname,patronymic,birthday,sex,nationaly_id,snils,birthplace)
VALUES (999999990,'TESTING_FRDO_CHECK_NAME3','2021-02-11 13:18:45.128',1317,'testing_frdo_check_surname3','testing_frdo_check_patronymic3','1982-02-01',1,151,'113-222-333 44 ','Moscow');