package udatasmith

type DatasmithUnrealScene struct {
	ActorMeshes  []ActorMesh  `xml:"ActorMesh"`
	StaticMeshes []StaticMesh `xml:"StaticMesh"`
}
