package main

import (
	"net/http"
	"log"
	"strings"
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
	"flag"
)

func main() {

	smbShare := "//127.0.0.1/share"
	targetDir := filepath.Join(os.Getenv("SNAP_COMMON"), "shares")
	listenAddr := ":9090"
	flag.StringVar(&smbShare, "smbShare", smbShare, "SMB target share")
	flag.StringVar(&targetDir, "targetDir", targetDir, "Local mount target")
	flag.StringVar(&listenAddr, "listenAddr", listenAddr, "Listen address for the minimal rest api")
	flag.Parse()

	mountHandler := NewSmbMountHandler(smbShare, targetDir)

	http.HandleFunc("/", mountHandler.handle)   // set router
	err := http.ListenAndServe(listenAddr, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type SmbMountHandler struct {
	smbShare  string
	targetDir string
}

func NewSmbMountHandler(smbShare string, targetDir string) *SmbMountHandler {
	handler := &SmbMountHandler{}
	handler.smbShare = smbShare
	handler.targetDir = targetDir
	return handler
}

func (this *SmbMountHandler) mount() (string, error) {
	err := ensureDir(this.targetDir)

	if err != nil {
		log.Printf("ERROR during mount ensuring target directory exists: %s", err.Error())
		return "", err
	}

	command := "mount.cifs"
	args := []string{"-o", "user=anonymous,pass=whatever,dom=WORKGROUP", this.smbShare, this.targetDir}

	return execute(command, args...)
}

func (this *SmbMountHandler) show() (string, error) {

	command := "mount"
	args := []string{"-t","cifs"}

	return execute(command, args...)
}

func (this *SmbMountHandler) umount() (string, error) {

	command := "umount"
	args := []string{"-l", "-t", "cifs", this.targetDir}
	return execute(command, args...)

}

func (this *SmbMountHandler) handle(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if (strings.Contains(path, "umount")) {
		output, err := this.umount()
		var response string

		if err != nil {
			response = fmt.Sprintf("UMOUNT\nERROR: %s\nOUTPUT: %s\n", err.Error(), output)
		} else {
			response = fmt.Sprintf("UMOUNT\nOUTPUT: %s\n", output)
		}

		log.Print(response)
		fmt.Fprint(w, response)
	} else if (strings.Contains(path, "mount")) {
		output, err := this.mount()
		var response string

		if err != nil {
			response = fmt.Sprintf("MOUNT\nERROR: %s\nOUTPUT: %s\n", err.Error(), output)
		} else {
			response = fmt.Sprintf("MOUNT\nOUTPUT: %s\n", output)
		}

		log.Print(response)
		fmt.Fprint(w, response)
	} else {
		output, err := this.show()
		var response string

		if err != nil {
			response = fmt.Sprintf("SHOW CIFS MOUNTS\nERROR: %s\nOUTPUT: %s\n", err.Error(), output)
		} else {
			response = fmt.Sprintf("SHOW CIFS MOUNTS\nOUTPUT: %s\n", output)
		}

		log.Print(response)
		fmt.Fprint(w, response)
	}
}

func ensureDir(dir string) error {
	err := os.Mkdir(dir, 0777)

	if err != nil && ! os.IsExist(err) {
		return err
	} else {
		return nil
	}
}

func execute(command string, args ...string) (string, error) {
	var cmdOut []byte
	var err error
	cmdPrint := fmt.Sprintf("Executing command: %s %s\n\n", command, strings.Join(args, " "))
	log.Print(cmdPrint)

	if cmdOut, err = exec.Command(command, args...).CombinedOutput(); err != nil {
		errorString := fmt.Sprintf("Error during command execution %s %s", err.Error(), string(cmdOut))
		log.Printf("ERROR during command execution: %s", errorString)
		return cmdPrint + string(cmdOut), err
	} else {
		return cmdPrint + string(cmdOut), nil
	}
}
