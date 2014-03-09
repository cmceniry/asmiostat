package asmiostat

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

var Subs map[string]string = make(map[string]string)
var Revs map[string]string = make(map[string]string)

func getRdev(dirname string) (map[uint64]string, error) {
	ret := make(map[uint64]string)
	dir, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	for _, fileinfo := range fileinfos {
		if !fileinfo.IsDir() {
			if stat, ok := fileinfo.Sys().(*syscall.Stat_t); ok {
				ret[stat.Rdev] = fileinfo.Name()
			} else {
				continue
			}
		}
	}
	return ret, nil
}

func GenerateSubs() error {
	newsubs := make(map[string]string)
	newrevs := make(map[string]string)

	asmdevs, err := getRdev("/dev/oracleasm/disks")
	if err != nil {
		return err
	}
	rawdevs, err := getRdev("/dev")
	if err != nil {
		return err
	}

	//fmt.Printf("-asmdevs: %v\n", asmdevs)
	//fmt.Printf("-rawdevs: %v\n", rawdevs)

	for asmrdev, asmname := range asmdevs {
		//fmt.Printf("-Checking %d %s\n", asmrdev, asmname)
		if rawname, found := rawdevs[asmrdev]; found {
			newsubs[rawname] = asmname
			newrevs[asmname] = rawname
		}
	}
	Subs = newsubs
	Revs = newrevs

	return nil
}

func startBackgroundSubs(e chan error) error {
	if err := GenerateSubs(); err != nil {
		return err
	}
	go (func() {
		for {
			fmt.Printf("-Subs: %s\n", Subs)
			time.Sleep(60 * time.Second)
			if err := GenerateSubs(); err != nil {
				e <- err
				break
			}
		}
	})()
	return nil
}
