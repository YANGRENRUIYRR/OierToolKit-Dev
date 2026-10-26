// Copyright (C) 2026 魇珩Gr3yPh4ntom
//
// This file is part of OierToolKit.
//
// OierToolKit is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// OierToolKit is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with OierToolKit.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 题面展示页：直接用 Fyne 内置的 Markdown 渲染，支持标题、粗体、斜体、
// 行内代码、代码块、列表、链接与表格（渲染范围以 Fyne 实现为准）。

// statementFileNames 是题面文件的候选名，按优先级排列。
var statementFileNames = []string{"problem.md", "statement.md", "题面.md", "README.md"}

// findStatement 在项目目录里找题面文件，返回路径或空串。
func findStatement(p project) string {
	for _, name := range statementFileNames {
		path := filepath.Join(p.Dir, name)
		if fileExists(path) {
			return path
		}
	}
	return ""
}

// statementPage 是一个只读的题面展示页，直接嵌进项目标签页。
type statementPage struct {
	pg      *projectPage
	content fyne.CanvasObject
	rich    *widget.RichText
	title   *widget.Label
}

func newStatementPage(pg *projectPage) *statementPage {
	s := &statementPage{pg: pg}
	s.title = hint("")
	s.rich = widget.NewRichText()
	s.rich.Wrapping = fyne.TextWrapWord
	reloadBtn := widget.NewButtonWithIcon(tr("重新加载"), theme.ViewRefreshIcon(), s.reload)
	head := container.NewBorder(nil, nil, nil, reloadBtn, container.NewPadded(s.title))
	s.content = container.NewBorder(
		container.NewVBox(head, widget.NewSeparator()),
		nil, nil, nil,
		container.NewVScroll(container.NewPadded(s.rich)),
	)
	s.reload()
	return s
}

// reload 读取题面文件并渲染；没有题面时给放置提示。
func (s *statementPage) reload() {
	path := findStatement(s.pg.p)
	if path == "" {
		s.title.SetText("没有找到题面文件")
		s.rich.ParseMarkdown("把题面保存为 `problem.md` / `statement.md` / `题面.md` 放在项目目录里，即可在这里渲染。")
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		s.title.SetText("读取失败：" + err.Error())
		s.rich.Segments = nil
		s.rich.Refresh()
		return
	}
	s.title.SetText(path)
	s.rich.ParseMarkdown(string(data))
}
