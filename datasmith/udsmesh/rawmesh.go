package udsmesh

type RawMesh struct {
	FaceMaterialIndices []int32
	FaceSmoothingMasks  []uint32
	VertexPositions     [][3]float32
	WedgeIndices        []uint32
	WedgeTangentX       [][3]float32
	WedgeTangentY       [][3]float32
	WedgeTangentZ       [][3]float32
	WedgeTexCoords      [][2]float32
}
