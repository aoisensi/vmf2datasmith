package udatasmith

type ActorMesh struct {
	Name      string `xml:",attr"`
	MeshName  string `xml:"Mesh>Name"`
	Transform Transform
}
