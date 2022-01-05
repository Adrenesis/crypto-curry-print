package model

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func LaunchPyReader(inputFile string, outputFile string, n int) {
	path, errwd := os.Getwd()
	if errwd != nil {
		log.Println(errwd)
	}
	fmt.Println(path)
	var stdout, stderr bytes.Buffer
	writer := io.MultiWriter(os.Stdout, os.Stderr)
	cmd := exec.Command("python", "./pyReader/main.py", inputFile, outputFile, fmt.Sprintf("%d", n))
	cmd.Stdout = writer
	cmd.Stderr = writer
	//if errout != nil {
	//	fmt.Println(errout)
	//}
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}
	//var p []byte
	//a, errread := stdout.Read(p)
	//print(a)
	//print(errread)
	//fmt.Println(p)
	fmt.Println("wtf")
	//time.Sleep(20* time.Second)
	//stdout.Close()
	if err := cmd.Wait(); err != nil {
		print("aaaa")
		if exiterr, ok := err.(*exec.ExitError); ok {

			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			status, ok := exiterr.Sys().(syscall.WaitStatus)
			if ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	print(string(stdout.Bytes()))
	print(string(stderr.Bytes()))
}
