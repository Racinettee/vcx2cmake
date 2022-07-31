package vcxproj

type ClCompile struct {
	// Include is the source file that is to be included by the compiler
	Include string `xml:"Include,attr"`
}

type ItemGroup struct {
	CompileTargets []ClCompile `xml:"ClCompile"`
}

type PropertyGroup struct {
	Label             string `xml:"Label,attr"`
	ConfigurationType string `xml:"ConfigurationType"`
	CharacterSet      string `xml:"CharacterSet"`
	Condition         string `xml:"Condition,attr"`
}

type Project struct {
	ItemGroups     []ItemGroup     `xml:"ItemGroup"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
}
