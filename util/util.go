package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func Download(url string, filename string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get jar: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create jar: %s", err)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write jar: %s", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close body: %s", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close jar: %s", err)
	}

	return nil
}

func Instructions(instructions []string) error {
	for _, instruction := range instructions {
		parts := strings.Split(instruction, ":")
		switch parts[0] {
		case "delete":
			target := parts[1]
			if strings.HasSuffix(target, "/") {
				err := os.RemoveAll(target)
				if err != nil {
					return fmt.Errorf("failed to execute remove folder instruction %s: %s", instruction, err)
				}
			} else {
				err := os.Remove(target)
				if err != nil {
					return fmt.Errorf("failed to execute remove instruction %s: %s", instruction, err)
				}
			}
		case "deletereg":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					err = os.Remove(file.Name())
					if err != nil {
						return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
					}
				}
			}
		case "rename":
			err := os.Rename(parts[1], parts[2])
			if err != nil {
				return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
			}
		case "renamereg":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					err = os.Rename(file.Name(), parts[2])
					if err != nil {
						return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
					}
				}
			}
		case "cdreg":
			folders, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, folder := range folders {
				if folder.IsDir() && r.MatchString(folder.Name()) {
					err = os.Chdir(folder.Name())
					if err != nil {
						return fmt.Errorf("failed to execute cd instruction %s: %s", instruction, err)
					}
				}
			}
		case "javarun":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					output, err := exec.Command("java", "-jar", file.Name(), parts[2]).Output()
					if err != nil {
						log.Println(string(output))
						return fmt.Errorf("failed to javarun for %s: %s", instruction, err)
					}
				}
			}
		case "forgegrep":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			forgeversion := ""
			forgeurl := ""

			for _, line := range strings.Split(string(file), "\n") {
				if strings.HasPrefix(line, "FORGE_VERSION") {
					forgeversion = strings.Split(line, "=")[1]
				}
				if strings.HasPrefix(line, "FORGE_URL") {
					forgeurl = strings.Split(line, "\"")[1]
				}
			}

			url := strings.ReplaceAll(forgeurl, "$FORGE_VERSION", forgeversion)
			err = Download(url, "forge-installer.jar")
			if err != nil {
				return fmt.Errorf("failed to download forge installer: %s", err)
			}

			output, err := exec.Command("java", "-jar", "forge-installer.jar", "--installServer").Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to run installer: %s", err)
			}

			err = os.Remove("forge-installer.jar")
			if err != nil {
				return fmt.Errorf("failed to remove installer: %s", err)
			}

			err = os.Remove("installer.log")
			if err != nil {
				return fmt.Errorf("failed to remove installer log: %s", err)
			}
		case "bashrun":
			output, err := exec.Command("bash", parts[1]).Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to bashrun for %s: %s", instruction, err)
			}
		}
	}
	return nil
}
