# db-cli

`db` command line tools (Docker based)

Type **only 3 strokes**, you can create Database.



## Installation

### step1. Install Docker

<http://www.docker.com/>

### step2. Install db-cli

If you are developer (case: Linux, MacOS)

	$ go get github.com/girigiribauer/db-cli

	$ cd $GOPATH/src/github.com/girigiribauer/db-cli

	$ go install

If you are end-user, you can download release binaries.

(Not yet...)

### step3. Check

	$ db -h



## Basic usage

### Create

	$ db

You can create MariaDB Database (Docker Container).

Maybe container name is "db0", port number is 3306.

These are **automatic assignment.**

### Delete

	$ db -d

You can delete this database.

### dump

	$ db -o

You can get dump file on "~/db0.sql"

### restore

	$ db --file=~/db0.sql

You can restore database with dump file.

### help

	$ db -h

or

	$ db --help



## Options

see `db -h`

	--name CONTAINER_NAME, -n CONTAINER_NAME

override CONTAINER_NAME, auto increment with prefix (default: db0, db1 ...)

	--dbname DB_NAME, -b DB_NAME

override DB_NAME (default: "db") [$DBCLI_DB_NAME]

	--dbuser DB_USER, -u DB_USER

override DB_USER (default: "db") [$DBCLI_DB_USER]

	--dbpass DB_PASS, -p DB_PASS

override DB_PASS (default: "db") [$DBCLI_DB_PASS]

	--image DOCKER_IMAGE, -i DOCKER_IMAGE

override DOCKER_IMAGE (default: "mariadb") [$DBCLI_DOCKER_IMAGE]

	--tag DOCKER_IMAGE_TAG, -t DOCKER_IMAGE_TAG

override docker image DOCKER_IMAGE_TAG (default: "latest")

	-d

delete one container db0, db1 ... (auto incrementation)

	--delete CONTAINER_NAME

delete container CONTAINER_NAME

	--delete-all

delete all db containers (without use docker command directly)

	-o

output dump file in default directory (default: "~/[CONTAINER_NAME].sql")

	--dump FILE_PATH

output dump file FILE_PATH

	--file FILE_PATH, -f FILE_PATH

restore with file FILE_PATH

	--list

list all db containers (without use docker command directly)



## Environment

You can set .bashrc, .zshrc or the other.

	DBCLI_CONTAINER_MAX

default: 100 (db0, db1 ... db99)

	DBCLI_CONTAINER_PREFIX

default: "db"

	DBCLI_DB_NAME

default: "db"

	DBCLI_DB_USER

default: "db"

	DBCLI_DB_PASS

default: "db"

	DBCLI_DOCKER_IMAGE

default: "mariadb"

	DBCLI_DIRECTORY

default: "~/"



## soon...

* distribute release binaries
* add PostgreSQL driver
