package core

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

type Repository struct {
	Worktree string
	GotDir string
	Conf *viper.Viper
}

func NewRepository(p string, force bool) (r Repository, err error) {
	r = Repository{
		Worktree: p,
		GotDir: path.Join(p, ".git"),
		Conf: viper.New(),
	}

	_, err = os.Stat(r.GotDir)
	if !(force || os.IsNotExist(err)) {
		return r, errors.New("not a repository")
	}
	confFile := RepoFile(r, "config")
	r.Conf.AddConfigPath(confFile)

	if err := r.Conf.ReadInConfig(); err != nil {
		return r, errors.New("unable to retrieve configuration file")
	}

	if !force {
		cmap := r.Conf.Get("core")
		if cmap != nil {
			return r, errors.New("unsupported version format")
		}
		c, ok := cmap.(map[string]interface{})
		if !ok {
			return r, errors.New("unsupported version format")
		}
		v, ok := c["version"]
		if !ok || v != 0 {
			return r, errors.New("unsupported version format")
		}
	}

	return r, nil
}

func CreateRepository(p string) (r Repository, err error) {

	// Make sure folder exists and is a folder
	// The example checks to make sure it's empty.. I'm not ogoing to do that
	fi, err := os.Stat(p)
	if !os.IsNotExist(err) {
		if !fi.IsDir() {
			return r, errors.New("not a directory")
		}
	} else {
		os.MkdirAll(p, 0777)
	}

	r, err = NewRepository(p, true)

	if d := RepoDir(r, true, "/branches"); d == "" {
		return r, errors.New("unable to create branches directory")
	}

	if d := RepoDir(r, true, "/objects"); d == "" {
		return r, errors.New("unable to create objects directory")
	}

	if d := RepoDir(r, true, "/refs/tags"); d == "" {
		return r, errors.New("unable to create /refs/tags directory")
	}

	if d := RepoDir(r, true, "/refs/heads"); d == "" {
		return r, errors.New("unable to create /refs/heads directory")
	}

	f, err := os.Create(RepoFile(r, "description"))
	if err != nil {
		return r, errors.New("unable to create description file")
	}
	defer f.Close()
	f.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))

	f, err = os.Create(RepoFile(r, "HEAD"))
	if err != nil {
		return r, errors.New("unable to create HEAD file")
	}
	defer f.Close()
	f.Write([]byte("ref: refs/heads/master\n"))

	config := map[string]string{
		"repositoryformatversion": "0",
		"filemode": "false",
		"bare": "false",
	}
	r.Conf.Set("core", config)
	if err := r.Conf.WriteConfig(); err != nil {
		return r, errors.New("unable to write config")
	}
	return r, nil
}

// Rethink these functions...
func RepoPath(repo Repository, p string) string {
	return path.Join(repo.GotDir, p)
}

func RepoFile(repo Repository, fp string) string {
	d := filepath.Dir(fp)
	if RepoDir(repo, true, d) != "" {
		return RepoPath(repo, fp)
	}
	return ""
}

func RepoDir(repo Repository, createDir bool, dirPath string) string {
	dirPath = path.Join(repo.GotDir, dirPath)
	fi, err := os.Stat(dirPath)
	if !os.IsNotExist(err) {
		if fi.IsDir() {
			return dirPath
		}
	}

	if createDir {
		os.MkdirAll(dirPath, 0777)
		return dirPath
	}

	return ""
}