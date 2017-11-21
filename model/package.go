package model

type Definition struct {
	Action         string `xml:"action,attr"`
	Description    string `xml:"description,attr"`
	Default        string `xml:"default,attr"`
	TransformTypes string `xml:"transformTypes"`
	From           string `xml:"from"`
	To             string `xml:"to"`
	Value          string
}

type Version struct {
	Name           string       `xml:"name,attr"`
	Folder         string       `xml:"folder,attr"`
	DriverFileName string       `xml:"driver,attr"`
	Definitions    []Definition `xml:"definition"`
}

type Package struct {
	Name           string    `xml:"name,attr"`
	Folder         string    `xml:"folder,attr"`
	DriverFileName string    `xml:"driver,attr"`
	Versions       []Version `xml:"version"`
}
