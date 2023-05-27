package main

import (
	"fmt"

	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

type LogView struct {
	TextEdit *d.TextEdit
	logEdit  *walk.TextEdit
}

func (c *LogView) Append(value string) {
	c.logEdit.AppendText(fmt.Sprintf("\r\n%v", value))
	c.logEdit.SetTextSelection(len(c.logEdit.Text())+1, 0)
}

func NewLogView() *LogView {
	lv := &LogView{}

	lv.TextEdit = &d.TextEdit{
		AssignTo: &lv.logEdit,
		VScroll:  true,
		HScroll:  true,
	}

	return lv
}
