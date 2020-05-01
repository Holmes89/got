package core

import (
	"errors"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrRepositoryDoesNotExist = errors.New("repository does not exist")
	ErrConfigurationDoesNotExist = errors.New("configuration does not exist")
)

type Repository struct {
	Worktree string
	GotDir string
	Conf Config
}

func NewRepository(dir string) (r *Repository, err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, errors.New("invalid path")
	}

	r = &Repository{
		Worktree: dir,
		GotDir: path.Join(dir, ".got"),
	}
	_, err = os.Stat(r.GotDir)
	if os.IsNotExist(err) {
		return nil, ErrRepositoryDoesNotExist
	}

	// TODO cleanup errors
	f, err := os.Open(path.Join(r.GotDir, "config"))
	if err != nil {
		return nil, errors.New("unable to open config")
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := toml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}

	//TODO check version

	return r, nil
}


//NewFile create new file within the got config repository
func (r *Repository) NewFile(name string) (*os.File, error) {
	name = path.Join(r.GotDir, name)
	if err := r.NewDirectory(filepath.Dir(name)); err != nil {
		return nil, err
	}
	return os.Create(name)
}

//NewDirectory creates new directory within the got repository
func (r *Repository) NewDirectory(name string) error {
	name = path.Join(r.GotDir, name)
	return os.MkdirAll(name, 0777)
}

type Config struct {
	Core ConfigCore `toml:"core"`
}

type ConfigCore struct {
	Version int `toml:"repositoryformatversion"`
	FileMode bool `toml:"filemode"`
	Bare bool `toml:"bare"`
}

func CreateRepository(dir string) (r *Repository, err error) {
	// Make sure folder exists and is a folder
	fi, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		if !fi.IsDir() {
			return r, errors.New("not a directory")
		}
	} else {
		os.MkdirAll(dir, 0777)
	}

	// TODO err checks
	gotDir := filepath.Join(dir, ".got")
	// TODO check if already a repo.
	os.MkdirAll(gotDir, 0777)
	f, err :=os.Create(path.Join(gotDir, "config"))
	config := Config{
		Core: ConfigCore{
			Version:  0,
			FileMode: false,
			Bare:     false,
		},
	}
	b, _ := toml.Marshal(config)
	f.Write(b)
	f.Close()

	r, err = NewRepository(dir)
	if err != nil && err != ErrConfigurationDoesNotExist{
		return nil, err
	}

	if err := r.NewDirectory("branches"); err != nil {
		return nil, errors.New("unable to create branches directory")
	}

	if err := r.NewDirectory("objects"); err != nil {
		return nil, errors.New("unable to create objects directory")
	}

	if err := r.NewDirectory("refs/tags"); err != nil {
		return nil, errors.New("unable to create /refs/tags directory")
	}

	if err := r.NewDirectory("refs/heads"); err != nil {
		return nil, errors.New("unable to create /refs/heads directory")
	}

	f, err = r.NewFile("description")
	if err != nil {
		return r, errors.New("unable to create description file")
	}
	f.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))
	f.Close()

	f, err = r.NewFile("HEAD")
	if err != nil {
		return r, errors.New("unable to create HEAD file")
	}
	f.Write([]byte("ref: refs/heads/master\n"))
	f.Close()
	return r, nil
}
