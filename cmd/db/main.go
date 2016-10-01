package main

import (
	"os"
	"time"

	"github.com/girigiribauer/db-cli"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "db"
	app.Version = "0.0.24"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "girigiribauer",
			Email: "girigiribauer@gmail.com",
		},
	}
	app.Copyright = "(c) 2016 girigiribauer"
	app.Usage = "the command line tool with docker (required Docker)"
	app.UsageText = "db [options]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "override `CONTAINER_NAME`, auto increment with prefix (default: db0, db1 ...)",
		},
		cli.StringFlag{
			Name:   "dbname, b",
			Usage:  "override `DB_NAME`",
			EnvVar: "DBCLI_DB_NAME",
			Value:  "db",
		},
		cli.StringFlag{
			Name:   "dbuser, u",
			Usage:  "override `DB_USER`",
			EnvVar: "DBCLI_DB_USER",
			Value:  "db",
		},
		cli.StringFlag{
			Name:   "dbpass, p",
			Usage:  "override `DB_PASS`",
			EnvVar: "DBCLI_DB_PASS",
			Value:  "db",
		},
		cli.StringFlag{
			Name:   "image, i",
			Usage:  "override `DOCKER_IMAGE`",
			EnvVar: "DBCLI_DOCKER_IMAGE",
			Value:  "mariadb",
		},
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "override docker image `DOCKER_IMAGE_TAG`",
			Value: "latest",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "delete one container db0, db1 ... (auto incrementation)",
		},
		cli.StringFlag{
			Name:  "delete",
			Usage: "delete container `CONTAINER_NAME`",
		},
		cli.BoolFlag{
			Name:  "delete-all",
			Usage: "delete all db containers (without use docker command directly)",
		},
		cli.BoolFlag{
			Name:  "o",
			Usage: "output dump file in default directory (default: \"~/[CONTAINER_NAME].sql\")",
		},
		cli.StringFlag{
			Name:  "dump",
			Usage: "output dump file `FILE_PATH`",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "restore with file `FILE_PATH`",
		},
		cli.BoolFlag{
			Name:  "list",
			Usage: "list all db containers (without use docker command directly)",
		},
		cli.StringFlag{
			Name:   "container-prefix",
			EnvVar: "DBCLI_CONTAINER_PREFIX",
			Hidden: true,
			Value:  "db",
		},
		cli.IntFlag{
			Name:   "container-max",
			EnvVar: "DBCLI_CONTAINER_MAX",
			Hidden: true,
			Value:  100,
		},
		cli.StringFlag{
			Name:   "directory",
			EnvVar: "DBCLI_DIRECTORY",
			Hidden: true,
			Value:  "~/",
		},
	}

	app.Action = func(c *cli.Context) error {
		name := c.String("name")
		action := "create"
		filepath := ""

		if c.Bool("list") {
			action = "list"
		} else if c.Bool("delete-all") {
			action = "delete"
			name = "*"
		} else if c.Bool("d") || c.String("delete") != "" {
			action = "delete"
			if c.String("name") != "" && c.String("delete") != "" {
				cli.NewExitError("both name and delete options", 1)
			}
			name = c.String("delete")
		} else if c.Bool("o") || c.String("dump") != "" {
			action = "dump"
			filepath = c.String("dump")
		} else if c.String("file") != "" {
			filepath = c.String("file")
		}

		incrementation := true
		if name != "" {
			incrementation = false
		}

		db.DB(action, &db.ContainerOptions{
			Incrementation:  incrementation,
			Name:            name,
			DBName:          c.String("dbname"),
			DBUser:          c.String("dbuser"),
			DBPass:          c.String("dbpass"),
			Filepath:        filepath,
			Image:           c.String("image"),
			Tag:             c.String("tag"),
			ContainerPrefix: c.String("container-prefix"),
			ContainerMax:    c.Int("container-max"),
			Directory:       c.String("directory"),
		})
		return nil
	}

	app.Run(os.Args)
}
