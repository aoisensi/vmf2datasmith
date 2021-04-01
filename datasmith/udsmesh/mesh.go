package udsmesh

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type DSMesh struct {
	Name string
	Raw  *RawMesh
}

type Encoder struct {
	w io.WriteSeeker
}

func NewEncoder(w io.WriteSeeker) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(mesh *DSMesh) error {
	var err error
	if err := e.write(uint32(1)); err != nil {
		return err
	}
	if err := e.write(uint32(0)); err != nil { //file size
		return err
	}
	fileStart, err := e.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if err := e.writeStr(mesh.Name); err != nil {
		return err
	}
	if err := e.write(byte(0)); err != nil {
		return err
	}
	if err := e.write(uint32(1)); err != nil {
		return err
	}
	if err := e.writeStr("SouceModels"); err != nil {
		return err
	}
	if err := e.writeStr("StructProperty"); err != nil {
		return err
	}
	if err := e.writeNull(8); err != nil {
		return err
	}
	if err := e.writeStr("DatasmithMeshSourceModel"); err != nil {
		return err
	}
	if err := e.writeNull(25); err != nil {
		return err
	}
	sizeLoc, err := e.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if err := e.write(uint64(0)); err != nil { // size
		return err
	}
	if err := e.write(uint32(125)); err != nil {
		return err
	}
	if err := e.write(uint32(0)); err != nil {
		return err
	}
	meshStart, err := e.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if err := e.write(uint32(1)); err != nil { //raw mesh version
		return err
	}
	if err := e.write(uint32(0)); err != nil { // licensee version
		return err
	}
	if err := e.writeSlice(mesh.Raw.FaceMaterialIndices); err != nil {
		return err
	}
	if err := e.writeSlice(mesh.Raw.FaceSmoothingMasks); err != nil {
		return err
	}
	if err := e.writeVectors(mesh.Raw.VertexPositions); err != nil {
		return err
	}
	if err := e.writeSlice(mesh.Raw.WedgeIndices); err != nil {
		return err
	}
	if err := e.writeVectors(mesh.Raw.WedgeTangentX); err != nil {
		return err
	}
	if err := e.writeVectors(mesh.Raw.WedgeTangentY); err != nil {
		return err
	}
	if err := e.writeVectors(mesh.Raw.WedgeTangentZ); err != nil {
		return err
	}
	if err := e.writeVector2s(mesh.Raw.WedgeTexCoords); err != nil {
		return err
	}
	if err := e.write(uint32(0)); err != nil { // skip colors
		return err
	}
	if err := e.write(uint32(0)); err != nil { // skip material index
		return err
	}
	meshEnd, err := e.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if err := e.writeNull(20); err != nil {
		return err
	}
	fileEnd, err := e.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if _, err := e.w.Seek(sizeLoc, io.SeekStart); err != nil {
		return err
	}
	if err := e.writeTwice(uint32(meshEnd - meshStart)); err != nil {
		return err
	}
	if _, err := e.w.Seek(4, io.SeekStart); err != nil {
		return err
	}
	return e.write(uint32(fileEnd - fileStart))
}

func (e *Encoder) writeVectors(data [][3]float32) error {
	if data == nil {
		return e.write(uint32(0))
	}
	for _, vv := range data {
		for _, v := range vv {
			if err := e.write(float32(v)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Encoder) writeVector2s(data [][2]float32) error {
	if data == nil {
		return e.write(uint32(0))
	}
	for _, vv := range data {
		for _, v := range vv {
			if err := e.write(float32(v)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Encoder) writeSlice(data interface{}) error {
	if data == nil {
		return e.write(uint32(0))
	}
	switch datas := data.(type) {
	case []uint32:
		if err := e.write(uint32(len(datas))); err != nil {
			return err
		}
		return e.write(datas)
	case []int32:
		if err := e.write(uint32(len(datas))); err != nil {
			return err
		}
		return e.write(datas)
	case [][3]float32:
		if err := e.write(uint32(len(datas))); err != nil {
			return err
		}
		return e.write(datas)
	case [][2]float32:
		if err := e.write(uint32(len(datas))); err != nil {
			return err
		}
		return e.write(datas)
	default:
		panic(fmt.Sprintf("%v type is not supported in writeSlice.", reflect.TypeOf(data)))
	}
}

func (e *Encoder) writeNull(size int) error {
	return e.write(make([]byte, size))
}

func (e *Encoder) writeStr(data string) error {
	bdata := []byte(data)
	if err := e.write(uint32(len(bdata) + 1)); err != nil {
		return err
	}
	if err := e.write(bdata); err != nil {
		return err
	}
	return e.write(byte(0))
}

func (e *Encoder) writeTwice(data interface{}) error {
	if err := e.write(data); err != nil {
		return err
	}
	return e.write(data)
}

func (e *Encoder) write(data interface{}) error {
	return binary.Write(e.w, binary.LittleEndian, data)
}
