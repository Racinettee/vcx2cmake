package main

import (
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"vcx2cmake/pkg/cmake"
	"vcx2cmake/pkg/vcxproj"
)

type OutputConfig struct {
	OutputDir  string `toml:"OutputDirectory"`
	CopySource bool   `toml:"CopySource"`
}

func main() {
	// The user can specify the vcx project they are interested in converting
	var args struct {
		InFile string `long:"input" description:"the file name of the vcxproj to convert" required:"true"`
	}
	var outConfig OutputConfig
	_, err := flags.Parse(&args)

	if err != nil {
		panic(err)
	}
	// A configuration file specifies the output parameters
	_, err = toml.DecodeFile("vcx2cm_config.toml", &outConfig)

	if err != nil {
		panic(err)
	}

	file, err := os.Open(args.InFile)

	if err != nil {
		panic(err)
	}
	xmlDecoder := xml.NewDecoder(file)

	var project vcxproj.Project

	err = xmlDecoder.Decode(&project)

	if err != nil {
		log.Panic("could not parse the project file")
	}

	fmt.Printf("%+v", project)

	cmakeFile := cmake.CMakeListsFile{
		MinRequiredVersion: 3.0,
		ProjectName:        strings.TrimSuffix(path.Base(args.InFile), filepath.Ext(args.InFile)),
	}

	// for starters, we'll just read the first property group that contains configuration type and generate an exe/lib for that
	exe := func() *cmake.Executable {
		for _, pg := range project.PropertyGroups {
			if len(pg.ConfigurationType) != 0 {
				switch pg.ConfigurationType {
				case "Application":
					return cmakeFile.AddExecutable(strings.TrimSuffix(path.Base(args.InFile), filepath.Ext(args.InFile)))
				}
			}
		}
		return nil
	}()
	exe.Files = func() []string {
		result := make([]string, 0)
		for _, ig := range project.ItemGroups {
			if len(ig.CompileTargets) != 0 {
				for _, ct := range ig.CompileTargets {
					result = append(result, ct.Include)
				}
			}
		}
		return result
	}()
	outDir := "./"
	if outConfig.OutputDir != "" {
		outDir = outConfig.OutputDir
	}
	err = os.MkdirAll(outConfig.OutputDir, 0o744)
	of, err := os.Create(outDir + "/CMakeLists.txt")
	if err != nil {
		panic(err)
	}
	cmakeFile.Encode(of)
	if outConfig.CopySource {
		for _, file := range exe.Files {
			copyFile(file, fmt.Sprintf("%s/%s", outDir, file))
		}
	}
}

func copyFile(src, dest string) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Panicf("%v\n", err)
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Panicf("%v is not a regular file\n", src)
		return
	}

	source, err := os.Open(src)
	if err != nil {
		log.Panicf("%v\n", err)
		return
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		log.Panicf("%v\n", err)
		return
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Panicf("%v\n", err)
	}
}
