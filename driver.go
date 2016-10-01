package db

import (
	"fmt"
)

type driver interface {
	portString() string
	portNumber() int64
	imageName() string
	envString() []string
	dumpCommands() []string
	restoreCommands() []string
	healthcheckCommand() string
}

type mysqlDriver struct {
	image    string
	port     int64
	protocol string
	name     string
	user     string
	pass     string
}

func newMariaDBDriver(image, tag, name, user, pass string) driver {
	return &mysqlDriver{
		image:    fmt.Sprintf("%s:%s", image, tag),
		port:     3306,
		protocol: "tcp",
		name:     name,
		user:     user,
		pass:     pass,
	}
}

func newMySQLDriver(image, tag, name, user, pass string) driver {
	return &mysqlDriver{
		image:    fmt.Sprintf("%s:%s", image, tag),
		port:     3306,
		protocol: "tcp",
		name:     name,
		user:     user,
		pass:     pass,
	}
}

func (driver *mysqlDriver) portString() string {
	return fmt.Sprintf("%d/%s", driver.port, driver.protocol)
}

func (driver *mysqlDriver) portNumber() int64 {
	return driver.port
}

func (driver *mysqlDriver) imageName() string {
	return driver.image
}

func (driver *mysqlDriver) envString() []string {
	return []string{
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", driver.pass),
		fmt.Sprintf("MYSQL_DATABASE=%s", driver.name),
		fmt.Sprintf("MYSQL_USER=%s", driver.user),
		fmt.Sprintf("MYSQL_PASSWORD=%s", driver.pass),
	}
}

func (driver *mysqlDriver) dumpCommands() []string {
	return []string{
		"mysqldump",
		fmt.Sprintf("-u%s", driver.user),
		fmt.Sprintf("-p%s", driver.pass),
		fmt.Sprintf("%s", driver.name),
	}
}

func (driver *mysqlDriver) restoreCommands() []string {
	return []string{
		"mysql",
		fmt.Sprintf("-u%s", driver.user),
		fmt.Sprintf("-p%s", driver.pass),
		fmt.Sprintf("%s", driver.name),
	}
}

func (driver *mysqlDriver) healthcheckCommand() string {
	//return "while ! mysqladmin ping --silent -udb -pdb; do sleep 1; done"
	return "while ! [ -e /var/run/mysqld/mysqld.sock ]; do sleep 1; done"
}
