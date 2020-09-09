package docker

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	RepositoryPassphrase = "DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE"
)

type (
	// Daemon defines Docker daemon parameters.
	Daemon struct {
		Registry      string   // Docker registry
		Mirror        string   // Docker registry mirror
		Insecure      bool     // Docker daemon enable insecure registries
		StorageDriver string   // Docker daemon storage driver
		StoragePath   string   // Docker daemon storage path
		Disabled      bool     // DOcker daemon is disabled (already running)
		Debug         bool     // Docker daemon started in debug mode
		Bip           string   // Docker daemon network bridge IP address
		DNS           []string // Docker daemon dns server
		DNSSearch     []string // Docker daemon dns search domain
		MTU           string   // Docker daemon mtu setting
		IPv6          bool     // Docker daemon IPv6 networking
		Experimental  bool     // Docker daemon enable experimental mode
	}

	// Login defines Docker login parameters.
	Login struct {
		Registry string // Docker registry address
		Username string // Docker registry username
		Password string // Docker registry password
		Email    string // Docker registry email
		Config   string // Docker Auth Config
	}

	// Trust defines Docker trust parameters
	Trust struct {
		Passphrase   string // Repository Passphrase for the key
		RepoKey      string // Repo Private key
		RootKeyName  string // Repo Private key
		RootKey      string // Root Private key
		RootCert     string // Root cert
		RootCertName string // Root cert name
		Repo         string // Repo to be signed
		Tag          string // Image Tag to be signed
	}

	// Plugin defines the Docker plugin parameters.
	Plugin struct {
		Login  Login  // Docker login configuration
		Daemon Daemon // Docker daemon configuration
		Trust  Trust  // Docker trust configuration
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	// start the Docker daemon server
	if !p.Daemon.Disabled {
		p.startDaemon()
	}

	// poll the docker daemon until it is started. This ensures the daemon is
	// ready to accept connections before we proceed.
	for i := 0; i < 15; i++ {
		cmd := commandInfo()
		err := cmd.Run()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	// Create Auth Config File
	if p.Login.Config != "" {
		if err := os.MkdirAll(dockerHome, 0600); err != nil {
			log.Fatal("Error initializing docker config", err)
		}

		path := filepath.Join(dockerHome, "config.json")
		err := ioutil.WriteFile(path, []byte(p.Login.Config), 0600)
		if err != nil {
			return fmt.Errorf("error writing config.json: %s", err)
		}
	}

	// initialize local docker trust
	initDCT()

	// login to the Docker registry
	if p.Login.Password != "" {
		fmt.Println("password is --------------", p.Login.Password)
		cmd := commandLogin(p.Login)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error authenticating: %s", err)
		}
	}

	switch {
	case p.Login.Password != "":
		fmt.Println("Detected registry credentials")
	case p.Login.Config != "":
		fmt.Println("Detected registry credentials file")
	default:
		fmt.Println("Registry credentials or Docker config not provided. Guest mode enabled.")
	}

	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion()) // docker version
	cmds = append(cmds, commandInfo())    // docker info

	// Pull the repository
	dockerPull(p.Trust.Repo, p.Trust.Tag)

	if p.Trust.RootKey != "" {
		os.Setenv(RepositoryPassphrase, p.Trust.Passphrase)
		copyCerts(p)
		//os.Setenv("DOCKER_CONTENT_TRUST","1")
		cmds = append(cmds, commandTrustKeyLoad(p.Trust, p.Trust.RootCertName))
		cmds = append(cmds, commandTrustSign(p.Trust, p.Trust.Tag))
	} else {
		log.Fatal("error! please provide rootkey!")
	}

	// Cleanup images
	{
		cmds = append(cmds, commandRmi(p.Trust.Repo, p.Trust.Tag))
		cmds = append(cmds, commandPrune())
	}

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		cmd.Run()

	}
	return nil
}

func commandDaemon(daemon Daemon) *exec.Cmd {
	fmt.Println("inside docker daemon ==========")
	args := []string{
		"--data-root", daemon.StoragePath,
		"--host=unix:///var/run/docker.sock",
	}

	if daemon.StorageDriver != "" {
		args = append(args, "-s", daemon.StorageDriver)
	}
	if daemon.Insecure && daemon.Registry != "" {
		args = append(args, "--insecure-registry", daemon.Registry)
	}
	if daemon.IPv6 {
		args = append(args, "--ipv6")
	}
	if len(daemon.Mirror) != 0 {
		args = append(args, "--registry-mirror", daemon.Mirror)
	}
	if len(daemon.Bip) != 0 {
		args = append(args, "--bip", daemon.Bip)
	}
	for _, dns := range daemon.DNS {
		args = append(args, "--dns", dns)
	}
	for _, dnsSearch := range daemon.DNSSearch {
		args = append(args, "--dns-search", dnsSearch)
	}
	if len(daemon.MTU) != 0 {
		args = append(args, "--mtu", daemon.MTU)
	}
	if daemon.Experimental {
		args = append(args, "--experimental")
	}
	return exec.Command(dockerdExe, args...)
}

func commandTrustKeyLoad(trust Trust, certName string) *exec.Cmd {
	rootKeyName := dockerTrustStore + trust.RootKeyName + ".key"
	fmt.Println("rootkeyname", rootKeyName)
	return exec.Command(dockerExe, "trust", "key", "load", rootKeyName, "--name", certName)
}

// helper function to create the docker login command.
func commandLogin(login Login) *exec.Cmd {
	if login.Email != "" {
		return commandLoginEmail(login)
	}
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		login.Registry,
	)
}

func commandLoginEmail(login Login) *exec.Cmd {
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		"-e", login.Email,
		login.Registry,
	)
}

func commandVersion() *exec.Cmd {
	return exec.Command(dockerExe, "version")
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command(dockerExe, "info")
}

func copyCerts(p Plugin) {
	rootKeyPath := dockerTrustStore + p.Trust.RootKeyName + ".key"
	loadTrustKeyAsFile(rootKeyPath, p.Trust.RootKey)
}

func commandTrustSign(trust Trust, tag string) *exec.Cmd {
	repo := fmt.Sprintf("%s:%s", trust.Repo, tag)
	return exec.Command(dockerExe, "trust", "sign", repo)
}

func commandPrune() *exec.Cmd {
	return exec.Command(dockerExe, "system", "prune", "-f")
}

func commandRmi(repo, tag string) *exec.Cmd {
	repoUrl := fmt.Sprintf("%s:%s", repo, tag)
	return exec.Command(dockerExe, "rmi", repoUrl)
}

func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

func loadTrustKeyAsFile(keyPath, trustKey string) string {
	if _, err := os.Stat(keyPath); err == nil {
		os.Remove(keyPath)
	}
	tKey := fmt.Sprintf("%s", trustKey)
	key, err := base64.URLEncoding.DecodeString(tKey)
	if err != nil {
		log.Fatal("Error decoding docker trust key! ", err)
	}
	if err := ioutil.WriteFile(keyPath, key, 0600); err != nil {
		log.Fatal("Error writing docker trust key!", err)
	}
	return keyPath
}

// Initialize docker trust repository
func initDCT() {
	cmd := exec.Command("docker", "trust", "key", "generate", "jeff")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error Generating docker trust keys", err)
	}
}

func dockerPull(repo, tag string) {
	repoUrl := fmt.Sprintf("%s:%s", repo, tag)
	cmd := exec.Command("docker", "pull", repoUrl)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error pulling docker image ", err)
	}
}
