package thumbnail

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var tmpdir string

func init() {
	if d, err := os.Stat("/dev/shm"); err == nil && d.IsDir() {
		tmpdir = "/dev/shm"
	} else {
		tmpdir = os.TempDir()
	}
}

// Agent for thumbnail
type Agent struct {
	Slice int
}

// NewAgent for thumbnail
func NewAgent() *Agent {
	return &Agent{
		Slice: 60,
	}
}

// Thumbnail ...
func (a *Agent) Thumbnail(file string) error {
	tmpPath := filepath.Join(tmpdir, fmt.Sprintf("%d", time.Now().Unix()))
	err := os.Mkdir(tmpPath, 0755)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer os.RemoveAll(tmpPath)
	dur, err := a.getDur(file)
	if err != nil {
		return err
	}
	for k, v := range dur.slice(a.Slice) {
		fp := filepath.Join(tmpPath, fmt.Sprintf("img%03d.jpg", k+1))
		// 		cmd := fmt.Sprintf(`ffmpeg -ss %s -i '%s' -vframes 1 -q:v 2 -vf "scale=iw/2:ih/2" %s`, v.string(), file, fp)
		cmd := fmt.Sprintf(`ffmpeg -ss %s -i '%s' -vframes 1 -q:v 2 -vf "scale='min(480,iw)':-1" %s`, v.string(), file, fp)
		log.Println(cmd)
		c := exec.Command("/bin/sh", "-c", cmd)

		var sout, serr bytes.Buffer
		c.Stdout = &sout
		c.Stderr = &serr

		err = c.Run()
		if err != nil {
			log.Println(serr.String())
			return err
		}
	}

	var fp string
	if strings.HasPrefix(file, "http") {
		fp = filepath.Base(file)
	} else {
		fp = file
	}

	err = a.Convert(tmpPath, fp+".jpg")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
