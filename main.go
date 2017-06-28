package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/mcbernie/t3-scheduler/schedule"
)

// Simple Scheduler for Typo3 Database
// Scans pages table for in scheduler.txt specified page-ids and looks in
// tt_content for updates / changes
// if some changes are found, send a mail to specified usernames
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "run with: t3_schedule TYPO3PATH")
		return
	}

	t3Path := os.Args[1]

	if stat, err := os.Stat(t3Path); err != nil || !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Typo3 Installation in %s not found!\n", t3Path)
		return
	}

	fileadminPath := filepath.Join(t3Path, "fileadmin")
	if stat, err := os.Stat(fileadminPath); err != nil || !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "fileadmin Path in Typo3 Installation %s not found!\n", fileadminPath)
		return
	}

	configurationFile := filepath.Join(fileadminPath, "scheduler.txt")
	if stat, err := os.Stat(configurationFile); err != nil || stat.IsDir() {
		fmt.Fprintf(os.Stderr, "configuration file scheduler.ini in %s not found!\n", fileadminPath)
		return
	}

	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, configurationFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on loading configuration from %s!\n(%s)\n", configurationFile, err.Error())
		return
	}

	s := schedule.Create(cfg, t3Path, fileadminPath)
	_, myerr := s.Run()
	if myerr != nil {
		fmt.Println(myerr.Error())
	}

	//fmt.Printf("nummer: %d", number)
}
