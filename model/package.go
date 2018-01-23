package model

//Definition defines the attributes to load the version definitions xml tag
type Definition struct {
	Action         string `xml:"action,attr"`
	Description    string `xml:"description,attr"`
	Default        string `xml:"default,attr"`
	TransformTypes string `xml:"transformTypes"`
	From           string `xml:"from"`
	To             string `xml:"to"`
	Value          string
}

//Version defines the attributes to load the package versions xml tag
type Version struct {
	Name           string       `xml:"name,attr"`
	Folder         string       `xml:"folder,attr"`
	DriverFileName string       `xml:"driver,attr"`
	Definitions    []Definition `xml:"definition"`
}

//Package defines the attributes to load the package xml file
type Package struct {
	Name           string    `xml:"name,attr"`
	Folder         string    `xml:"folder,attr"`
	DriverFileName string    `xml:"driver,attr"`
	Versions       []Version `xml:"version"`
}
