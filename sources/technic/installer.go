package technic

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"log"
	"os"
	"os/exec"
)

func Install(pack *TechnicPack) error {
	instructions := config.Global.Packs[pack.Id]
	if len(instructions) == 0 {
		return fmt.Errorf("instructions for pack are too short/not available")
	}

	err := os.Mkdir(pack.Name, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir(pack.Name)
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = util.Download(pack.ServerPackURL, "serverpack.zip")
	if err != nil {
		return fmt.Errorf("failed to download: %s", err)
	}

	output, err := exec.Command("unzip", "-q", "serverpack.zip").Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to unzip: %s", err)
	}

	err = os.Remove("serverpack.zip")
	if err != nil {
		return fmt.Errorf("failed to delete server pack: %s", err)
	}

	err = util.Instructions(instructions)
	if err != nil {
		return fmt.Errorf("failed to run instructions: %s", err)
	}

	zipname := fmt.Sprintf("%s.zip", pack.Version)
	err = exec.Command("zip", "-qr", zipname, "./").Run()
	if err != nil {
		return fmt.Errorf("failed to zip pack: %s", err)
	}

	folder := fmt.Sprintf("%stechnic/%s/", config.Global.Target, pack.Name)
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create target folder: %s", err)
		}
	}

	err = os.Rename(zipname, folder+zipname)
	if err != nil {
		return fmt.Errorf("failed to move zip to target folder: %s", err)
	}

	log.Printf("Finished installing %s version %s", pack.DisplayName, pack.Version)

	return nil
}