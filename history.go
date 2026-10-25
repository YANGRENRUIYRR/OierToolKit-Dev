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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 评测历史：每次完整评测（编译并跑完全部样例）追加写入项目的
// .otk/history.json，最多保留 historyMax 条，超出时丢最旧的。
// 单点重测不进历史，避免刷屏。

const historyMax = 100

// historyEntry 是一次评测的快照。
type historyEntry struct {
	Time      string        `json:"time"`                 // RFC3339 本地时间
	Passed    int           `json:"passed"`               // 通过测试点数
	Total     int           `json:"total"`                // 测试点总数
	TimeLimit float64       `json:"time_limit"`           // 时限（秒）
	MemLimit  float64       `json:"mem_limit"`            // 内存限制（MB）
	Points    []pointResult `json:"points,omitempty"`     // 逐点明细
}

// historyPath 返回评测历史文件路径。
func historyPath(p project) string { return filepath.Join(p.metaDir(), "history.json") }

// appendHistory 把一次评测结果追加进历史。落盘失败静默忽略：
// 历史是锦上添花，不值得为它打断评测流程。
func appendHistory(p project, rep judgeReport) {
	if len(rep.Points) == 0 {
		return
	}
	entries := readHistory(p, 0)
	entry := historyEntry{
		Time:      time.Now().Format("2006-01-02 15:04:05"),
		Passed:    rep.passed(),
		Total:     len(rep.Points),
		TimeLimit: rep.TimeLimit,
		MemLimit:  rep.MemLimit,
		Points:    rep.Points,
	}
	entries = append(entries, entry)
	if len(entries) > historyMax {
		entries = entries[len(entries)-historyMax:]
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(p.metaDir(), 0755)
	_ = os.WriteFile(historyPath(p), data, 0644)
}

// readHistory 读取评测历史，limit > 0 时只返回最近 limit 条（时间正序）。
func readHistory(p project, limit int) []historyEntry {
	data, err := os.ReadFile(historyPath(p))
	if err != nil {
		return nil
	}
	var entries []historyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil
	}
	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}
	return entries
}

// cmdHistory 查看当前项目的评测历史（默认最近 10 次）。
func cmdHistory(args []string) {
	if len(args) > 1 {
		fatal("用法: otk history [次数]")
	}
	n := 10
	if len(args) == 1 {
		v, err := strconv.Atoi(args[0])
		if err != nil || v <= 0 {
			fatal("次数必须是正整数")
		}
		n = v
	}

	p := requireProject()
	entries := readHistory(p, n)
	if len(entries) == 0 {
		fmt.Println("（还没有评测记录，先 otk j 跑一次）")
		return
	}

	for _, e := range entries {
		verd := green + "AC" + reset
		if e.Passed < e.Total {
			verd = red + "未通过" + reset
			if first := firstWorstVerdict(e.Points); first != "" {
				verd = red + first + reset
			}
		}
		fmt.Printf("%s  %d/%d  %s  (%.2fs / %.0fMB)\n",
			e.Time, e.Passed, e.Total, verd, e.TimeLimit, e.MemLimit)
	}
}

// clearHistory 删除项目的全部评测历史。
func clearHistory(p project) {
	_ = os.Remove(historyPath(p))
}

// firstWorstVerdict 返回逐点明细里第一个未通过的判定（历史列表展示用）。
func firstWorstVerdict(points []pointResult) string {
	for _, pt := range points {
		if pt.Verdict != verdictAC {
			return pt.Verdict
		}
	}
	return ""
}
