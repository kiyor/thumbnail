package thumbnail

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (a *Agent) getDur(file string) (*videoDuration, error) {
	u, err := url.Parse(file)
	if err != nil {
		return nil, err
	}
	var path string
	p := strings.Split(u.Path, "/")
	if len(p) != 5 {
		path = file
	} else {
		path = u.String()
	}
	reDuration := regexp.MustCompile(`Duration: (\d{2}:\d{2}:\d{2}.\d{2}),`)
	cmd := fmt.Sprintf(`ffmpeg -i '%s' 2>&1`, path)
	log.Println(cmd)
	c := exec.Command("/bin/sh", "-c", cmd)

	var sout bytes.Buffer
	c.Stdout = &sout

	c.Run()

	p = reDuration.FindStringSubmatch(sout.String())
	var length string
	if len(p) == 2 {
		length = p[1]
	} else {
		return nil, errors.New("parse error")
	}
	dur, err := parseDuration(length)
	if err != nil {
		return nil, err
	}
	return dur, nil
}

type videoDuration struct {
	Hour        int
	Minute      int
	Second      int
	Millisecond int
	d           time.Duration
}

func (d *videoDuration) string() string {
	return fmt.Sprintf("%02d:%02d:%02d.%02d", d.Hour, d.Minute, d.Second, d.Millisecond/10)
}

func newVideoDuration(dur time.Duration) *videoDuration {
	var d videoDuration
	d.d = dur

	d.Millisecond = int(d.d.Seconds()*1000) % 1000
	d.Second = int(d.d.Seconds()) % 60
	d.Minute = int(d.d.Minutes()) % 60
	d.Hour = int(d.d.Hours()) % 60

	return &d
}

func parseDuration(in string) (*videoDuration, error) {
	p := strings.Split(in, ":")
	if len(p) != 3 {
		return nil, errors.New("error parse " + in)
	}
	var d videoDuration
	var err error
	d.Hour, err = strconv.Atoi(p[0])
	if err != nil {
		return nil, err
	}
	d.d += time.Duration(d.Hour) * time.Hour

	d.Minute, err = strconv.Atoi(p[1])
	if err != nil {
		return nil, err
	}
	d.d += time.Duration(d.Minute) * time.Minute

	last := strings.Split(p[2], ".")
	if len(last) != 2 {
		return nil, errors.New("error")
	}
	d.Second, err = strconv.Atoi(last[0])
	if err != nil {
		return nil, err
	}
	d.d += time.Duration(d.Second) * time.Second

	d.Millisecond, err = strconv.Atoi(last[1])
	if err != nil {
		return nil, err
	}
	d.Millisecond = d.Millisecond * 10
	d.d += time.Duration(d.Millisecond) * time.Millisecond

	return &d, nil
}

// Slice ...
func (d *videoDuration) slice(piece int) []*videoDuration {
	l := d.d.Nanoseconds() / int64(piece)
	var out []*videoDuration
	for i := int64(1000000000); i < d.d.Nanoseconds(); i += l {
		dur := newVideoDuration(time.Duration(i) * time.Nanosecond)
		out = append(out, dur)
	}
	return out
}
