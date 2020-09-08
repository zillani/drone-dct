package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	docker "github.com/zillani/drone-dct"
	"os"
)

var (
	version = "unknown"
)

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}
	app := cli.NewApp()
	app.Name = "dct plugin"
	app.Usage = "dct plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "dry run disables docker push",
			EnvVar: "PLUGIN_DRY_RUN",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "daemon.mirror",
			Usage:  "docker daemon registry mirror",
			EnvVar: "PLUGIN_MIRROR",
		},
		cli.StringFlag{
			Name:   "daemon.storage-driver",
			Usage:  "docker daemon storage driver",
			EnvVar: "PLUGIN_STORAGE_DRIVER",
		},
		cli.StringFlag{
			Name:   "daemon.storage-path",
			Usage:  "docker daemon storage path",
			Value:  "/var/lib/docker",
			EnvVar: "PLUGIN_STORAGE_PATH",
		},
		cli.StringFlag{
			Name:   "daemon.bip",
			Usage:  "docker daemon bride ip address",
			EnvVar: "PLUGIN_BIP",
		},
		cli.StringFlag{
			Name:   "daemon.mtu",
			Usage:  "docker daemon custom mtu setting",
			EnvVar: "PLUGIN_MTU",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns",
			Usage:  "docker daemon dns server",
			EnvVar: "PLUGIN_CUSTOM_DNS",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns-search",
			Usage:  "docker daemon dns search domains",
			EnvVar: "PLUGIN_CUSTOM_DNS_SEARCH",
		},
		cli.BoolFlag{
			Name:   "daemon.insecure",
			Usage:  "docker daemon allows insecure registries",
			EnvVar: "PLUGIN_INSECURE",
		},
		cli.BoolFlag{
			Name:   "daemon.ipv6",
			Usage:  "docker daemon IPv6 networking",
			EnvVar: "PLUGIN_IPV6",
		},
		cli.BoolFlag{
			Name:   "daemon.experimental",
			Usage:  "docker daemon Experimental mode",
			EnvVar: "PLUGIN_EXPERIMENTAL",
		},
		cli.BoolFlag{
			Name:   "daemon.debug",
			Usage:  "docker daemon executes in debug mode",
			EnvVar: "PLUGIN_DEBUG,DOCKER_LAUNCH_DEBUG",
		},
		cli.BoolFlag{
			Name:   "daemon.off",
			Usage:  "don't start the docker daemon",
			EnvVar: "PLUGIN_DAEMON_OFF",
		},
		cli.StringFlag{
			Name:   "docker.username",
			Usage:  "docker username",
			EnvVar: "DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.password",
			Usage:  "docker password",
			EnvVar: "DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.email",
			Usage:  "docker email",
			EnvVar: "DOCKER_EMAIL",
		},
		cli.StringFlag{
			Name:   "docker.config",
			Usage:  "docker json dockerconfig content",
			EnvVar: "PLUGIN_CONFIG",
		},
		cli.StringFlag{
			Name:   "passphrase",
			Usage:  "DCT Repo passphrase",
			EnvVar: "PLUGIN_PASSPHRASE",
		},
		cli.StringFlag{
			Name:   "rootkey",
			Usage:  "Docker Container Trust Root Private key",
			EnvVar: "PLUGIN_ROOTKEY",
		},
		cli.StringFlag{
			Name:   "repokey",
			Usage:  "Docker Container Trust Repo Private Key",
			EnvVar: "PLUGIN_REPOKEY",
		},
		cli.StringFlag{
			Name:   "rootkeyname",
			Usage:  "Docker Container Trust Root Private key name",
			EnvVar: "PLUGIN_ROOTKEYNAME",
		},
		cli.StringFlag{
			Name:   "rootcert",
			Usage:  "Docker Container Trust Root cert",
			EnvVar: "PLUGIN_ROOTCERT",
		},
		cli.StringFlag{
			Name:   "rootcertname",
			Usage:  "Docker Container Trust Repo cert name",
			EnvVar: "PLUGIN_ROOTCERTNAME",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "Repo to be signed",
			EnvVar: "PLUGIN_REPO",
		},
		cli.StringFlag{
			Name:   "tag",
			Usage:  "Tag to be signed",
			EnvVar: "PLUGIN_TAG",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := docker.Plugin{
		Login: docker.Login{
			Registry: c.String("docker.registry"),
			Username: c.String("docker.username"),
			Password: c.String("docker.password"),
			Email:    c.String("docker.email"),
			Config:   c.String("docker.config"),
		},
		Daemon: docker.Daemon{
			Registry:      c.String("docker.registry"),
			Mirror:        c.String("daemon.mirror"),
			StorageDriver: c.String("daemon.storage-driver"),
			StoragePath:   c.String("daemon.storage-path"),
			Insecure:      c.Bool("daemon.insecure"),
			Disabled:      c.Bool("daemon.off"),
			IPv6:          c.Bool("daemon.ipv6"),
			Debug:         c.Bool("daemon.debug"),
			Bip:           c.String("daemon.bip"),
			DNS:           c.StringSlice("daemon.dns"),
			DNSSearch:     c.StringSlice("daemon.dns-search"),
			MTU:           c.String("daemon.mtu"),
			Experimental:  c.Bool("daemon.experimental"),
		},
		Trust: docker.Trust{
			Passphrase:   c.String("passphrase"),
			RepoKey:      c.String("repokey"),
			RootKey:      c.String("rootkey"),
			RootKeyName:  c.String("rootkeyname"),
			RootCert:     c.String("rootcert"),
			RootCertName: c.String("rootcertname"),
			Repo:         c.String("repo"),
			Tag:          c.String("tag"),
		},
	}
	return plugin.Exec()
}
