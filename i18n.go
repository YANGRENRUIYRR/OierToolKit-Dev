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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

// 界面语言偏好：auto（跟随默认中文）/ zh / en。
// 词条以中文原文为键，en 字典给出英文；tr() 在英文模式下查表，
// 没有对应词条时回退中文原文。语言改动需要重启生效。

const prefLang = "gui_lang"

// languageChoices 是设置页里的可选语言。
var languageChoices = []string{"跟随系统", "中文", "English"}

// enUI 是中 → 英词条表，覆盖界面骨架与主要操作按钮。
var enUI = map[string]string{
	// 项目功能页
	"概览": "Overview", "源码": "Source", "题面": "Statement",
	"样例": "Samples", "评测": "Judge", "对拍": "Stress", "提交": "Submit",
	// 模拟赛页
	"榜单": "Board", "题目配置": "Problems", "选手与提交": "Roster", "评测日志": "Log",
	"开始评测": "Judge All", "停止": "Stop", "导出 HTML 榜单": "Export HTML",
	"导出 CSV": "Export CSV", "返回比赛列表": "Back",
	// 对拍页
	"开始对拍": "Start", "存为样例": "Save as Sample", "管理反例": "Manage",
	"打开 .stress 目录": "Open .stress", "发现反例后继续收集": "Collect more after WA",
	// 源码页
	"保存": "Save", "重新加载": "Reload", "外部编辑器": "External Editor",
	"gdb 调试": "gdb Debug", "统计": "Stats",
	// 样例页
	"保存样例": "Save Sample", "同步": "Sync", "从剪贴板导入": "Import Clipboard",
	// 概览页
	"同步样例": "Sync Samples", "打开目录": "Open Folder",
	"在终端打开": "Open Terminal", "移除项目…": "Remove Project…",
	"配置远端": "Remotes…",
	// 设置
	"跟随系统": "Auto", "窗口模式": "Window Mode",
	"全屏启动": "Start fullscreen", "窗口化启动（恢复上次窗口大小与位置）": "Windowed (restore last size & position)",
}

// uiLangIsEnglish 判断界面是否使用英文。
func uiLangIsEnglish() bool {
	return guiPref(prefLang, "auto") == "en"
}

// tr 翻译界面文案：英文模式查 enUI 字典，否则原样返回中文。
func tr(s string) string {
	if uiLangIsEnglish() {
		if en, ok := enUI[s]; ok {
			return en
		}
	}
	return s
}

// trSelect 翻译语言选项显示值（设置页里把"跟随系统"在英文下也译出）。
func trSelect(s string) string {
	if uiLangIsEnglish() {
		switch s {
		case "跟随系统":
			return "Auto"
		case "中文":
			return "Chinese"
		}
	}
	return s
}

// installTranslations 把内置中文词条注册到 Fyne 语言包。加载失败不影响
// 使用（只会退回英文原文），因此错误被忽略。界面语言选英文时不注册，
// 让 Fyne 内部控件（文件选择器、确认框等）保持英文。
func installTranslations() {
	if uiLangIsEnglish() {
		return
	}
	_ = lang.AddTranslationsForLocale(guiZhLangData, fyne.Locale("zh-Hans"))
}

// languageValueFromChoice 把下拉框显示值转成偏好存储值。
func languageValueFromChoice(display string) string {
	switch display {
	case "中文":
		return "zh"
	case "English":
		return "en"
	}
	return "auto"
}

// languageChoiceFromValue 把偏好存储值转回下拉框显示值。
func languageChoiceFromValue(v string) string {
	switch v {
	case "zh":
		return "中文"
	case "en":
		return "English"
	}
	return "跟随系统"
}
