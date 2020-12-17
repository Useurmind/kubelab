// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type CodeGen mg.Namespace

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