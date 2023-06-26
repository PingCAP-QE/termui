// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/gizak/termui/v3"
)

type Paragraph struct {
	Block
	Text      string
	TextStyle Style
	WrapText  bool
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		Block:     *NewBlock(),
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
	}
}

func (p *Paragraph) ScrollHalfPageUp() {}

func (p *Paragraph) ScrollHalfPageDown() {}

func (p *Paragraph) Get() string {
	return p.Text
}

func (p *Paragraph) Edit() string {
	ioutil.WriteFile("./.tmpfile", []byte(p.Text), 0644)
	cmd := exec.Command("vim", "./.tmpfile")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	tmp, err := ioutil.ReadFile("./.tmpfile")
	if err != nil {
		return ""
	}
	p.Text = string(tmp)
	return p.Text
}

func (p *Paragraph) Draw(buf *Buffer) {
	p.Block.Draw(buf)

	cells := ParseStyles(p.Text, p.TextStyle)
	if p.WrapText {
		cells = WrapCells(cells, uint(p.Inner.Dx()))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+p.Inner.Min.Y >= p.Inner.Max.Y {
			break
		}
		row = TrimCells(row, p.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(p.Inner.Min))
		}
	}
}
