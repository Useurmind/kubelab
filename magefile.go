// +build mage

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type CodeGen mg.Namespace
type Projects mg.Namespace
type UI mg.Namespace
type KubelabDev mg.Namespace

func (KubelabDev) HelmInstall() error {
	return sh.Run("helm", "upgrade", "--install", "kubelab", "--namespace", "kubelab", "--values", "./helm_chart/values-dev.yaml", "./helm_chart/")
}

func (CodeGen) Build() error {
	return sh.Run("go", "build", "-o", ".codegen/bin/codegen.exe", ".codegen/main.go")
}

func (CodeGen) Run() error {
	return codegen("projects", "ui")
}

func (UI) Build() error {
	// build webpack
	err := func() error {
		err := os.Chdir("services/ui")
		if err != nil {
			return err
		}
		defer os.Chdir("../..")

		err = sh.Run("npx", "webpack-cli", "-c", "webpack.config.js")
		if err != nil {
			return err
		}
	
		err = copyAllFiles("www", "bin/www")	
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return err
	}

	err = buildGo("ui")
	if err != nil {
		return err
	}

	return nil
}

func (UI) Run() error {
	return sh.Run("services/projects/bin/ui.exe")
}

func (UI) RunWebpack() error {
	err := os.Chdir("services/ui")
	if err != nil {
		return err
	}
	defer os.Chdir("../..")

	return sh.Run("npx", "webpack-dev-server")
}

func (UI) DockerBuildLocal() error {
	return sh.Run("docker", "build", "-f", "./services/ui/Dockerfile", "-t", "kubelab-ui:local-dev", ".")
}

func (Projects) Build() error {
	err := buildGo("projects")
	if err != nil {
		return err
	}

	err = copyAllFiles("services/projects/db", "services/projects/bin/db")	
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

func buildGo(serviceNames ...string) error {
	for _, serviceName := range serviceNames {
		err := func() error {
			err := os.Chdir(fmt.Sprintf("services/%s/api", serviceName))
			if err != nil {
				return err
			}
			defer os.Chdir("../../..")

			err = os.RemoveAll("../bin")
			if err != nil {
				return err
			}

			binName := serviceName
			if runtime.GOOS == "windows" {
				binName = binName + ".exe"
			}

			return sh.Run("go", "build", "-o", "../bin/" + binName)
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func codegen(serviceNames ...string) error {
	for _, serviceName := range serviceNames {
		err := func() error {
			err := os.Chdir("services/" + serviceName)
			if err != nil {
				return err
			}
			defer os.Chdir("../..")

			return sh.Run("../../.codegen/bin/codegen.exe")
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func copyAllFiles(srcDir, dstDir string) error {
	log.Printf("Copying all files from %s to %s", srcDir, dstDir)
	err := os.MkdirAll(dstDir, 777)
	if err != nil {
		return err
	}

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
		log.Printf(" - %s -> %s", srcPath, dstPath)
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
