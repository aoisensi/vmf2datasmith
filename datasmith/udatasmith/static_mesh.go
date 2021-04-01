package udatasmith

type StaticMesh struct {
	FilePath string `xml:"File>Path"`
	Label    string `xml:",attr"`
	Name     string `xml:",attr"`
}
