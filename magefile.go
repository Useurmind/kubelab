// +build mage

package main

import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type CodeGen mg.Namespace
type Projects mg.Namespace
type KubelabDev mg.Namespace

func (KubelabDev) HelmInstall() error {
	return sh.Run("helm", "upgrade", "--install", "kubelab", "--namespace", "kubelab", "--values", "./helm_chart/values-dev.yaml", "./helm_chart/")
}

func (CodeGen) Build() error {
	return sh.Run("go", "build", "-o", ".codegen/bin/codegen.exe", ".codegen/main.go")
}

func (CodeGen) Projects() error {
	err := os.Chdir("services/projects")
	if err != nil {
		return err
	}
	defer os.Chdir("../..")

	return sh.Run("../../.codegen/bin/codegen.exe")
}

func (Projects) Build() error {
	err := os.Chdir("services/projects/api")
	if err != nil {
		return err
	}
	defer os.Chdir("../../..")

	err = os.RemoveAll("../bin")
	if err != nil {
		return err
	}

	binName := "projects"
	if runtime.GOOS == "windows" {
		binName = binName + ".exe"
	}

	err = sh.Run("go", "build", "-o", "../bin/" + binName)
	if err != nil {
		return err
	}

	err = os.MkdirAll("../bin/db", 777)
	if err != nil {
		return err
	}

	err = copyAllFiles("../db", "../bin/db")	
	if err != nil {
		return err
	}

	return nil
}

func (Projects) Run() error {
	return sh.Run("services/projects/bin/projects.exe")
}

func (Projects) DockerBuildLocal() error {
	return sh.Run("docker", "build", "-f", "./services/projects/Dockerfile", "-t", "kubelab-projects:local-dev", ".")
}

func copyAllFiles(srcDir, dstDir string) error {
	fileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range fileInfos {
		if file.IsDir() {
			continue
		}

		srcPath := fmt.Sprintf("%s/%s", srcDir, file.Name())
		dstPath := fmt.Sprintf("%s/%s", dstDir, file.Name())
		_, err = copyFile(srcPath, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
			return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
			return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
			return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}