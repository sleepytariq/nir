package player

import (
	"fmt"
	"os/exec"
)

type Player struct {
	Name string
	Args []string
}

var MPV = Player{
	Name: "mpv",
	Args: []string{"--keep-open=no"},
}

var VLC = Player{
	Name: "vlc",
	Args: []string{"--play-and-exit"},
}

func (p *Player) Play(URL string, filename string) error {
	if p.Name == "mpv" {
		p.Args = append(p.Args, fmt.Sprintf("--force-media-title=%s", filename))
	}

	if p.Name == "vlc" {
		p.Args = append(p.Args, fmt.Sprintf("--meta-title=%s", filename))
	}

	cmd := exec.Command(p.Name, append(p.Args, URL)...)

	fmt.Printf("Playing %s with %s...\n", filename, p.Name)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to play stream %s <%s>", filename, err)
	}

	return nil
}
