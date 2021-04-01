package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aoisensi/vmf2datasmith/datasmith/udatasmith"
	"github.com/aoisensi/vmf2datasmith/datasmith/udsmesh"
	"github.com/aoisensi/vmf2datasmith/vmf"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/spatial/r3"
)

var rootCmd = &cobra.Command{
	Use:   "vmf2datasmith [VMF INPUT] [OUTPUT]",
	Args:  cobra.RangeArgs(2, 2),
	Short: "Transform vmf to UE's datasmith",
	Run:   runRoot,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func runRoot(cmd *cobra.Command, args []string) {
	vmfFile, err := os.Open(args[0])
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	defer vmfFile.Close()
	vmf, err := vmf.NewDecoder(vmfFile).Decode()
	if err != nil {
		cmd.Println(err)
		return
	}
	/*
		_, err = os.Create(args[1]) //TODO
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	*/
	os.Chdir(filepath.Dir(args[1]))
	os.MkdirAll("Geometries", os.ModePerm)
	vmfWorld := vmf.Nodes("world")[0]
	dsScene := udatasmith.DatasmithUnrealScene{}
	dsScene.ActorMeshes = make([]udatasmith.ActorMesh, 0, vmfWorld.CountNodes("solid"))
	// Burshes
	for _, vmfSolid := range vmfWorld.Nodes("solid") {
		type Face struct {
			Plane          Plane
			Material       string
			VertexIndecies []int
		}
		faces := make([]Face, vmfSolid.CountNodes("side"))
		toolMaterial := false
		for i, vmfSide := range vmfSolid.Nodes("side") {
			plane := parse3Vec(vmfSide.String("plane"))
			faces[i].Plane = planeFromPoints(plane)
			material := strings.ToLower(vmfSide.String("material"))
			if strings.HasPrefix(material, "tools/") {
				toolMaterial = true
				break
			}
			faces[i].Material = material
			faces[i].VertexIndecies = make([]int, 0, 16)
		}
		if toolMaterial {
			continue
		}
		id := vmfSolid.ID()
		name := fmt.Sprintf("S_Solid%v", id)
		facesLen := len(faces)
		verteces := make([]r3.Vec, 0, 256)
		for i := 0; i < facesLen-2; i++ {
			faceI := faces[i]
			for j := i + 1; j < facesLen-1; j++ {
				faceJ := faces[j]
				for k := j + 1; k < facesLen; k++ {
					faceK := faces[k]
					ok := true
					v := calcIntersection(faceI.Plane, faceJ.Plane, faceK.Plane)
					if v == nil {
						continue
					}
					for _, faceL := range faces {
						if r3.Dot(faceL.Plane.V, *v)+faceL.Plane.D < -EPS {
							ok = false
							break
						}
					}
					if ok {
						verteces = append(verteces, *v)
						index := len(verteces) - 1
						faces[i].VertexIndecies = append(faces[i].VertexIndecies, index)
						faces[j].VertexIndecies = append(faces[j].VertexIndecies, index)
						faces[k].VertexIndecies = append(faces[k].VertexIndecies, index)
					}
				}
			}
		}

		mesh := &udsmesh.DSMesh{Name: name, Raw: &udsmesh.RawMesh{}}

		mesh.Raw.VertexPositions = make([][3]float32, len(verteces))
		mesh.Raw.WedgeIndices = make([]uint32, 0, 1024)
		for i, vertex := range verteces {
			mesh.Raw.VertexPositions[i] = toUEVec3(vertex)
		}

		for _, face := range faces {
			vertex := func(i int) r3.Vec {
				return verteces[face.VertexIndecies[i]]
			}
			// Calc average
			center := r3.Vec{}
			for _, i := range face.VertexIndecies {
				center = r3.Add(center, verteces[i])
			}
			center = r3.Scale(1.0/float64(len(face.VertexIndecies)), center)
			viLen := len(face.VertexIndecies)
			for n := 0; n < viLen-2; n++ {
				a := r3.Unit(r3.Sub(vertex(n), center))
				p := planeFromPoints([3]r3.Vec{vertex(n), center, r3.Add(center, face.Plane.V)})
				smallestAngle := -1.0
				smallest := -1
				for m := n + 1; m < viLen; m++ {
					side := p.Classify(vertex(m))
					if side < EPS {
						continue
					}
					b := r3.Unit(r3.Sub(vertex(m), center))
					angle := r3.Dot(a, b)
					if angle > smallestAngle {
						smallestAngle = angle
						smallest = m
					}
				}
				face.VertexIndecies[n+1], face.VertexIndecies[smallest] = face.VertexIndecies[smallest], face.VertexIndecies[n+1]
			}
			for i := 1; i < viLen-1; i++ {
				mesh.Raw.WedgeIndices = append(
					mesh.Raw.WedgeIndices,
					uint32(face.VertexIndecies[0]),
					uint32(face.VertexIndecies[i]),
					uint32(face.VertexIndecies[i+1]),
				)
			}
		}
		meshFilePath := filepath.Join("Geometries", name+".udsmesh")
		meshFile, err := os.Create(meshFilePath)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		if err := udsmesh.NewEncoder(meshFile).Encode(mesh); err != nil {
			cmd.PrintErrln(err)
			return
		}
		meshFile.Close()
		cmd.Printf("Created %v\n", meshFilePath)
	}
}
