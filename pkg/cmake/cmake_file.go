package cmake

import (
	"fmt"
	"io"
)

type Executable struct {
	// add_executable(<name> ...)
	Name string
	// add_executable(... <files>)
	Files []string
	//  target_link_libraries
	Libraries []string
}

func (exe *Executable) Encode(writer io.Writer) {
	fmt.Fprintf(writer, "add_executable(%v\n", exe.Name)
	for _, file := range exe.Files {
		fmt.Fprintf(writer, "%v\n", file)
	}
	fmt.Fprintln(writer, ")")
}

type CMakeListsFile struct {
	MinRequiredVersion float32
	ProjectName        string
	Executables        []*Executable
}

func (cmakeFile *CMakeListsFile) AddExecutable(name string) *Executable {
	result := &Executable{
		Name: name,
	}
	cmakeFile.Executables = append(cmakeFile.Executables, result)
	return result
}

func (cmakeFile *CMakeListsFile) Encode(writer io.Writer) {
	fmt.Fprintf(writer, "cmake_minimum_required(VERSION %.1f)\n", cmakeFile.MinRequiredVersion)
	fmt.Fprintf(writer, "project(%v)\n", cmakeFile.ProjectName)

	for _, exe := range cmakeFile.Executables {
		exe.Encode(writer)
	}
}
