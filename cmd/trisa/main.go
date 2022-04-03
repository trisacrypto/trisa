/*
 The TRISA CLI client allows you to create and execute TRISA requests from the command
 line for development or testing purposes. For more information on how to use the CLI,
 run `trisa --help` or see the documenation at https://trisa.dev.
*/
package main

import (
	"fmt"
	"os"

	"github.com/trisacrypto/trisa/pkg"
	client "github.com/trisacrypto/trisa/pkg/trisa/cli"
	cli "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := cli.NewApp()
	app.Name = "trisa"
	app.Usage = "create and execute TRISA requests"
	app.Version = pkg.Version()
	app.Before = loadProfiles
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Usage:   "specify the path to a configuration file to use",
			EnvVars: []string{"TRISA_CONF"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "install",
			Usage:  "create a profile configuration for TRISA operations",
			Action: install,
			Flags:  []cli.Flag{},
		},
		{
			Name:    "profiles",
			Aliases: []string{"profile"},
			Usage:   "manage profiles",
			Action:  manageProfiles,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list available profiles and exit",
				},
				&cli.BoolFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "print the path to the configuration and exit",
				},
				&cli.StringFlag{
					Name:    "activate",
					Aliases: []string{"a"},
					Usage:   "activate the specified profile",
				},
				&cli.StringFlag{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "create a profile with the specified name",
				},
				&cli.StringFlag{
					Name:    "update",
					Aliases: []string{"u"},
					Usage:   "update profile with specified name",
				},
			},
		},
	}
	app.Run(os.Args)
}

//====================================================================================
// RPC Commands
//====================================================================================

var profiles *client.Profiles

func install(c *cli.Context) (err error) {
	// Create new profiles with the default editor, setting the path if specified.
	profiles = client.New()
	if conf := c.String("conf"); conf != "" {
		profiles.SetPath(conf)
	}

	editor := profiles.NewEditor(client.ProfileDefault, client.EditInstall)
	if err = editor.Edit(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func manageProfiles(c *cli.Context) (err error) {
	if c.Bool("list") {
		fmt.Println("available profiles:")
		for name := range profiles.Profiles {
			if name == profiles.Active {
				fmt.Printf("  * %s\n", name)
			} else {
				fmt.Printf("  - %s\n", name)
			}
		}
		return nil
	}

	if c.Bool("path") {
		path, _ := profiles.Path()
		fmt.Println(path)
		return nil
	}

	if name := c.String("activate"); name != "" {
		if err = profiles.SetActive(name); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	if name := c.String("create"); name != "" {
		editor := profiles.NewEditor(name, client.EditCreate)
		if err = editor.Edit(); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	if name := c.String("update"); name != "" {
		editor := profiles.NewEditor(name, client.EditUpdate)
		if err = editor.Edit(); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	var data []byte
	if data, err = yaml.Marshal(profiles.GetActive()); err != nil {
		return cli.Exit("could not print active profile", 1)
	}
	fmt.Print(string(data))
	return nil
}

//====================================================================================
// Helper Commands
//====================================================================================

func loadProfiles(c *cli.Context) (err error) {
	// do not load profiles in some cases
	args := c.Args()
	if c.NArg() == 0 || args.First() == "install" || args.First() == "help" {
		return nil
	}

	if path := c.String("conf"); path != "" {
		if profiles, err = client.LoadPath(path); err != nil {
			if err == client.ErrProfileNotFound {
				return cli.Exit("no profiles configured run trisa install first", 1)
			}
			return cli.Exit(err, 1)
		}
		return nil
	}

	if profiles, err = client.Load(); err != nil {
		if err == client.ErrProfileNotFound {
			return cli.Exit("no profiles configured run trisa install first", 1)
		}
		return cli.Exit(err, 1)
	}
	return nil
}
