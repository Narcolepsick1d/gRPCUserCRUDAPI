    create table users(
        id int generated by default as identity ,
        name varchar,
        age int,
        phone varchar
    );
    insert into users (name, age, phone) values ('garik',22,'3030');