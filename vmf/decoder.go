package vmf

import (
	"bufio"
	"io"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decode() (*Node, error) {
	return d.readNode()
}

func (d *Decoder) readNode() (*Node, error) {
	node := newNode()
	for {
		if err := d.trim(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		end, err := d.maybe('}')
		if err != nil {
			return nil, err
		}
		if end {
			break
		}
		if err := d.trim(); err != nil {
			return nil, err
		}
		key, err := d.readString()
		if err != nil {
			return nil, err
		}
		if err := d.trim(); err != nil {
			return nil, err
		}
		isNode, err := d.maybe('{')
		if err != nil {
			return nil, err
		}
		if isNode {
			n, err := d.readNode()
			if err != nil {
				return nil, err
			}
			node.addNode(key, n)
		} else {
			v, err := d.readString()
			if err != nil {
				return nil, err
			}
			node.data[key] = v
		}
	}
	return node, nil
}

func (d *Decoder) readString() (string, error) {
	r, _, err := d.r.ReadRune()
	if err != nil {
		return "", err
	}
	if r == '"' {
		str, err := d.r.ReadString('"')
		if err != nil {
			return "", err
		}
		return str[:len(str)-1], nil
	} else {
		text := string([]rune{r})
		for {
			r, _, err := d.r.ReadRune()
			if err != nil {
				return "", err
			}
			if isWhite(r) {
				break
			}
			text += string([]rune{r})
		}
		return text, d.r.UnreadRune()
	}
}

func (d *Decoder) maybe(c rune) (bool, error) {
	r, _, err := d.r.ReadRune()
	if err != nil {
		return false, err
	}
	if r == c {
		return true, nil
	} else {
		return false, d.r.UnreadRune()
	}
}

func (d *Decoder) trim() error {
	for {
		r, _, err := d.r.ReadRune()
		if err != nil {
			return err
		}
		if isWhite(r) {
			continue
		}
		return d.r.UnreadRune()
	}
}
func isWhite(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}
