package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const title = ` __  __     ______     __   __     ______     ______     __   __        ______     ______     ______     ______     _____    
/\ \/ /    /\  __ \   /\ "-.\ \   /\  == \   /\  __ \   /\ "-.\ \      /\  == \   /\  __ \   /\  __ \   /\  == \   /\  __-.  
\ \  _"-.  \ \  __ \  \ \ \-.  \  \ \  __<   \ \  __ \  \ \ \-.  \     \ \  __<   \ \ \/\ \  \ \  __ \  \ \  __<   \ \ \/\ \ 
 \ \_\ \_\  \ \_\ \_\  \ \_\\"\_\  \ \_____\  \ \_\ \_\  \ \_\\"\_\     \ \_____\  \ \_____\  \ \_\ \_\  \ \_\ \_\  \ \____- 
  \/_/\/_/   \/_/\/_/   \/_/ \/_/   \/_____/   \/_/\/_/   \/_/ \/_/      \/_____/   \/_____/   \/_/\/_/   \/_/ /_/   \/____/ 
                                                                                                                             `

func reactivateList(lists []*widgets.List, active int) {
	for i,list := range lists {
		if i == active {
			list.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
		} else {
			list.SelectedRowStyle = ui.NewStyle(ui.ColorWhite)
		}
		ui.Render(list)
	}
}

func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}

func main() {
    if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	header := widgets.NewParagraph()
	header.Text = title

	todo := widgets.NewList()
	todo.Title = "To Do"
	todo.TitleStyle.Fg = ui.ColorMagenta
	todo.Rows = []string{
		"[1] Automatic watering system",
		"[2] Automatic temp system",
		"[6] VIM cheatsheet",
		"[8] Prepare for Lightning Talk",
	}
	todo.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	todo.WrapText = false
	todo.SetRect(0, 0, 25, 8)

	inprog := widgets.NewList()
	inprog.Title = "In Progress"
	inprog.TitleStyle.Fg = ui.ColorBlue
 	inprog.Rows = []string{
		"[7] CLI Kanband board",
	}
	inprog.WrapText = false
	inprog.SetRect(0, 0, 25, 8)

	done := widgets.NewList()
	done.Title = "Done"
	done.TitleStyle.Fg = ui.ColorGreen
	done.Rows = []string{
		"[3] Time reporting tool",
		"[4] Dockerize things",
		"[5] Email Anna",
	}
	done.WrapText = false
	done.SetRect(0, 0, 25, 8)

	info := widgets.NewParagraph()
	info.Text = "Press 'n' to progress to the next column, 'b' to move to previous column.\nMove around with 'h', 'j', 'k' and 'l'"
	info.Border = false

	grid.Set(
		ui.NewRow(
			1.0/4,
			ui.NewCol(1.0, header),
		),
		ui.NewRow(
			1.0/8,
			ui.NewCol(1.0, info),
		),
		ui.NewRow(
			5.0/8,
			ui.NewCol(1.0 / 3, todo),
			ui.NewCol(1.0 / 3, inprog),
			ui.NewCol(1.0 / 3, done),
		),
	)

	ui.Render(grid)

	lists := []*widgets.List{todo, inprog, done}
	activeList := 0

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			lists[activeList].ScrollDown()
		case "k", "<Up>":
			lists[activeList].ScrollUp()
		case "l", "<Right>":
			if activeList == 2 {
				activeList = 0
			} else {
				activeList++
			}
			reactivateList(lists, activeList)
		case "h", "<Left>":
			if activeList == 0 {
				activeList = 2
			} else {
				activeList--
			}
			reactivateList(lists, activeList)
		case "g":
			if previousKey == "g" {
				todo.ScrollTop()
			}
		case "G", "<End>":
			todo.ScrollBottom()
		case "n":
			toMove := lists[activeList].Rows[lists[activeList].SelectedRow]
			lists[activeList].Rows = remove(lists[activeList].Rows, lists[activeList].SelectedRow)
			ui.Render(lists[activeList])
			if activeList == 2 {
				activeList = 0
			} else {
				activeList++
			}
			lists[activeList].Rows = append(lists[activeList].Rows, toMove)
			reactivateList(lists, activeList)
		case "b":
			toMove := lists[activeList].Rows[lists[activeList].SelectedRow]
			lists[activeList].Rows = remove(lists[activeList].Rows, lists[activeList].SelectedRow)
			ui.Render(lists[activeList])
			if activeList == 0 {
				activeList = 2
			} else {
				activeList--
			}
			lists[activeList].Rows = append(lists[activeList].Rows, toMove)
			reactivateList(lists, activeList)
		}
		

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(lists[activeList])
	}
}
