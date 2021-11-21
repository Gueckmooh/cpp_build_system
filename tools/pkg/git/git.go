package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"tools/pkg/options"
	"tools/pkg/utils"
)

type GitRepository struct {
	Url    string `xml:",chardata"`
	Target string `xml:"target,attr"`
}

func CloneRepository(url string, target string) error {
	if !utils.DirExists(path.Dir(target)) {
		if err := utils.Mkdir(path.Dir(target)); err != nil {
			return fmt.Errorf("git.CloneRepository: %s", err.Error())
		}
	}

	cmd := exec.Command("git", "clone", url, target)

	var stdBuffer bytes.Buffer
	if !options.GetOptionBool("quiet") {
		fmt.Printf("Cloning git repository \"%s\" into \"%s\"\n", url, target)

		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git.CloneRepository: %s", err.Error())
	}

	return nil
}

func (g *GitRepository) Clone() error {
	return CloneRepository(g.Url, g.Target)
}
