USE school
CREATE TABLE Stu_Union (
    sno CHAR(5) NOT NULL UNIQUE,
    name VARCHAR(8),
	sex CHAR(1),
    age INT ,
    dept CHAR(20),
    CONSTRAINT PK_Stu_Union PRIMARY KEY (sno)
);

insert into Stu_Union values('10001','Alice','1',20,'CS');
insert into Stu_Union values('10002','Bob','0',21,'EE');
insert into Stu_Union values('23100','Charlie','1',22,'CS');
insert into Stu_Union values('23101','David','0',23,'EE');
insert into Stu_Union values('95001','Eve','1',24,'CS');
insert into Stu_Union values('95002','Frank','0',25,'EE');

CREATE TABLE COURSE (
    cno CHAR(4) NOT NULL UNIQUE,
    cname VARCHAR(50) NOT NULL,
    cpoint INT,
    CONSTRAINT PK_COURSE PRIMARY KEY (cno)
);

insert into COURSE values('0001','Computer Networks',3);
insert into COURSE values('0002','Database',2);

CREATE TABLE SC (
    sno CHAR(5) REFERENCES Stu_Union(sno) ON DELETE CASCADE ,
    cno CHAR(4) REFERENCES COURSE(cno) ON DELETE CASCADE ,
    grade INT,
    CONSTRAINT PK_SC PRIMARY KEY (sno, cno)
);

insert into SC values('10001','0001',1);
insert into SC values('10001','0002',1);
insert into SC values('23100','0001',2);
insert into SC values('23100','0002',2);
insert into SC values('95002','0001',3);
insert into SC values('95002','0002',3);


select * from SC;

ALTER TABLE SC DROP CONSTRAINT FK_SC_sno;

ALTER TABLE SC ADD CONSTRAINT FK_SC_sno FOREIGN KEY (sno) REFERENCES Stu_Union(sno) ON DELETE NO ACTION;

select * from SC;

delete from Stu_Union where sno='95002';
select * from SC;

ALTER TABLE SC DROP CONSTRAINT FK_SC_sno;
ALTER TABLE SC ADD CONSTRAINT FK_SC_sno FOREIGN KEY (sno) REFERENCES Stu_Union(sno) ON DELETE SET NULL;
SELECT * FROM SC;

CREATE TABLE Stu_Card(
    card_id CHAR(14),
    stu_id CHAR(10) REFERENCES students(sid) ON DELETE CASCADE,
    remaining_balance DECIMAL(10,2),
    CONSTRAINT PK_Stu_Card PRIMARY KEY (card_id)
);

insert into Stu_Card values('20240000000001','800014991 ',100.00);
insert into Stu_Card values('20240000000002','800016416 ',50.00);

CREATE TABLE ICBC_Card(
    bank_id CHAR(20),
    stu_card_id CHAR(14) REFERENCES Stu_Card(card_id) ON DELETE CASCADE,
    restored_balance DECIMAL(10,2),
    CONSTRAINT PK_ICBC_Card PRIMARY KEY (bank_id)
);

insert into ICBC_Card values('6222024000000001','20240000000001',3000.00);
insert into ICBC_Card values('6222024000000002','20240000000002',5200.00);

ALTER TABLE CHOICES DROP CONSTRAINT FK_CHOICES_STUDENTS;
ALTER TABLE CHOICES ADD CONSTRAINT FK_CHOICES_STUDENTS FOREIGN KEY (sid) REFERENCES students(sid) ON DELETE CASCADE;

BEGIN TRANSACTION T3
    ALTER TABLE ICBC_Card DROP CONSTRAINT FK__ICBC_Card__stu_c__719CDDE7;
    ALTER TABLE ICBC_Card ADD CONSTRAINT FK_Stu_Card_stu_id FOREIGN KEY (stu_card_id) REFERENCES Stu_Card(card_id) ON DELETE SET NULL;
    DELETE FROM students WHERE sid='800014991';
    SELECT * FROM ICBC_Card;
COMMIT TRANSACTION T3;