package main


import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
)

func main() {
	/*
		db, err := sql.Open("sqlite3", executionPath+"/backend/database/db2.sqlite")
		checkErr(err)
	*/
	db, err := sql.Open("mysql", "ctf:1234@/ctfengine")

	if err != nil {
		log.Panic(err)
	}

	_, errCreateUsers := db.Query(`CREATE TABLE ctfengine.users (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
						username varchar(100) NULL,
						email varchar(100) NULL,
						password varchar(100) NULL,
						status INT DEFAULT 0 NULL,
						points INT DEFAULT 0 NOT NULL,
						flagright INT UNSIGNED DEFAULT 0 NOT NULL,
						flagfalse INT UNSIGNED DEFAULT 0 NOT NULL,
						commandid INT UNSIGNED DEFAULT 0 NOT NULL
					)
					ENGINE=InnoDB
					DEFAULT CHARSET=utf8
					COLLATE=utf8_general_ci;
				`)

	if errCreateUsers != nil {
		log.Println(errCreateUsers)
	} else {
		log.Println("Users DB created")
	}


	_, errCreatePwnedBy := db.Query(`CREATE TABLE ctfengine.pwnedby (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT primary key,
						userid INT UNSIGNED NULL,
						command_id INT UNSIGNED NULL,
						taskid INT UNSIGNED NULL,
						contestid INT UNSIGNED NULL,
						time VARCHAR(64) NOT NULL
					)
					ENGINE=InnoDB
					DEFAULT CHARSET=utf8
					COLLATE=utf8_general_ci;
				`)

	if errCreatePwnedBy != nil {
		log.Println(errCreatePwnedBy)
	} else {
		log.Println("Pwnedby DB created")
	}


	_, errCreateRating := db.Query(`CREATE TABLE ctfengine.rating (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT,
						contest_id INT UNSIGNED NULL,
						userid INT UNSIGNED NULL,
						team INT UNSIGNED NULL,
						points INT UNSIGNED NULL,
						CONSTRAINT rating_pk PRIMARY KEY (id)
					)
					ENGINE=InnoDB
					DEFAULT CHARSET=utf8
					COLLATE=utf8_general_ci;
				`)

	if errCreateRating != nil {
		log.Println(errCreateRating)
	} else {
		log.Println("Pwnedby DB created")
	}

	_, errCreateCommand := db.Query(`CREATE TABLE ctfengine.command (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT primary key,
						name varchar(100) NULL,
						captainid INT UNSIGNED NULL,
						invite varchar(1024) NOT NULL
                    )
					ENGINE=InnoDB
					DEFAULT CHARSET=utf8
					COLLATE=utf8_general_ci;
				`)

	if errCreateCommand != nil {
		log.Println(errCreateCommand)
	} else {
		log.Println("Command DB created")
	}

	_, errCreateTasks := db.Query(`CREATE TABLE ctfengine.tasks (
					  	id INT NOT NULL AUTO_INCREMENT,
					 	name varchar(100) NULL,
						contestid INT UNSIGNED NULL,
						value INT UNSIGNED NULL,
						flag varchar(100) NULL,
						description varchar(8000) NULL,
						category varchar(100) NULL,
						CONSTRAINT tasks_pk PRIMARY KEY (id)
					)
					ENGINE=InnoDB
					DEFAULT CHARSET=utf8
					COLLATE=utf8_general_ci;
		  		`)

	if errCreateTasks != nil {
		log.Println(errCreateTasks)
	} else {
		log.Println("Tasks DB created")
	}

	_, errCreateContests := db.Query(`CREATE TABLE ctfengine.contests (
							id INT UNSIGNED NOT NULL AUTO_INCREMENT,
							name varchar(100) NULL,
							type varchar(100) NULL,
							visibility BOOL DEFAULT 1 NOT NULL,
							permit BOOL DEFAULT 1 NOT NULL,
							CONSTRAINT contests_pk PRIMARY KEY (id)
						)
						ENGINE=InnoDB
						DEFAULT CHARSET=utf8
						COLLATE=utf8_general_ci;
		  			`)

	if errCreateContests != nil {
		log.Println(errCreateContests)
	} else {
		log.Println("Contests DB created")
	}

	_, errCreateAttachments := db.Query(`CREATE TABLE attachments(
    							id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
   								taskid int,
    							name varchar(100)
							);
		  				`)

	if errCreateAttachments != nil {
		log.Println(errCreateAttachments)
	} else {
		log.Println("Attachments DB created")
	}

	_, errCreateIdorTaskTable := db.Query(`CREATE TABLE articles(
    							id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
   								userid int,
    							article varchar(2000),
    							time VARCHAR(64) NOT NULL
							);
		  				`)

	if errCreateIdorTaskTable != nil {
		log.Println(errCreateIdorTaskTable)
	} else {
		log.Println("IdorTask DB created")
	}
}
