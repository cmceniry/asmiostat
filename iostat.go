package asmiostat

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func substitute(line string) string {
	for orig, replace := range Subs {
		if re, err := regexp.Compile("^" + orig + strings.Repeat(" ", 16-len(orig))); err == nil {
			line = re.ReplaceAllString(line, replace+strings.Repeat(" ", 16-len(replace)))
		}
	}
	return line
}

func startIostat(options optionset, output chan string, errors chan error) error {
	args := make([]string, len(options.IostatArgs))
	for idx, arg := range options.IostatArgs {
		if sub, found := Revs[arg]; found {
			args[idx] = sub
		} else {
			args[idx] = arg
		}
	}

	cmd := exec.Command("iostat", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	out := bufio.NewReader(stdout)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			if buf, err := out.ReadString('\n'); err != nil {
				errors <- err
			} else {
				output <- substitute(buf)
			}
		}
	}()
	return nil
}
