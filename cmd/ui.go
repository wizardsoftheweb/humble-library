package wotwhb

import (
	"log"
	"math"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// func RunUi() {
// 	if err := ui.Init(); err != nil {
// 		log.Fatalf("failed to initialize termui: %v", err)
// 	}
// 	defer ui.Close()
//
// 	p := widgets.NewParagraph()
// 	p.Text = "Hello World!"
// 	p.SetRect(0, 0, 25, 5)
//
// 	ui.Render(p)
//
// 	for e := range ui.PollEvents() {
// 		if e.Type == ui.KeyboardEvent {
// 			break
// 		}
// 	}
// }
//
// type nodeValue string
//
// func (nv nodeValue) String() string {
// 	return string(nv)
// }
//
// func RunUi() {
// 	if err := ui.Init(); err != nil {
// 		log.Fatalf("failed to initialize termui: %v", err)
// 	}
// 	defer ui.Close()
//
// 	nodes := []*widgets.TreeNode{
// 		{
// 			Value: nodeValue("Key 1"),
// 			Nodes: []*widgets.TreeNode{
// 				{
// 					Value: nodeValue("Key 1.1"),
// 					Nodes: []*widgets.TreeNode{
// 						{
// 							Value: nodeValue("Key 1.1.1"),
// 							Nodes: nil,
// 						},
// 						{
// 							Value: nodeValue("Key 1.1.2"),
// 							Nodes: nil,
// 						},
// 					},
// 				},
// 				{
// 					Value: nodeValue("Key 1.2"),
// 					Nodes: nil,
// 				},
// 			},
// 		},
// 		{
// 			Value: nodeValue("Key 2"),
// 			Nodes: []*widgets.TreeNode{
// 				{
// 					Value: nodeValue("Key 2.1"),
// 					Nodes: nil,
// 				},
// 				{
// 					Value: nodeValue("Key 2.2"),
// 					Nodes: nil,
// 				},
// 				{
// 					Value: nodeValue("Key 2.3"),
// 					Nodes: nil,
// 				},
// 			},
// 		},
// 		{
// 			Value: nodeValue("Key 3"),
// 			Nodes: nil,
// 		},
// 	}
//
// 	l := widgets.NewTree()
// 	l.TextStyle = ui.NewStyle(ui.ColorYellow)
// 	l.WrapText = false
// 	l.SetNodes(nodes)
//
// 	x, y := ui.TerminalDimensions()
//
// 	l.SetRect(0, 0, x, y)
//
// 	ui.Render(l)
//
// 	previousKey := ""
// 	uiEvents := ui.PollEvents()
// 	for {
// 		e := <-uiEvents
// 		switch e.ID {
// 		case "q", "<C-c>":
// 			return
// 		case "j", "<Down>":
// 			l.ScrollDown()
// 		case "k", "<Up>":
// 			l.ScrollUp()
// 		case "<C-d>":
// 			l.ScrollHalfPageDown()
// 		case "<C-u>":
// 			l.ScrollHalfPageUp()
// 		case "<C-f>":
// 			l.ScrollPageDown()
// 		case "<C-b>":
// 			l.ScrollPageUp()
// 		case "g":
// 			if previousKey == "g" {
// 				l.ScrollTop()
// 			}
// 		case "<Home>":
// 			l.ScrollTop()
// 		case "<Enter>":
// 			l.ToggleExpand()
// 		case "G", "<End>":
// 			l.ScrollBottom()
// 		case "E":
// 			l.ExpandAll()
// 		case "C":
// 			l.CollapseAll()
// 		case "<Resize>":
// 			x, y := ui.TerminalDimensions()
// 			l.SetRect(0, 0, x, y)
// 		}
//
// 		if previousKey == "g" {
// 			previousKey = ""
// 		} else {
// 			previousKey = e.ID
// 		}
//
// 		ui.Render(l)
// 	}
// }

func loadUiData() (*MungeComponents, [][]HbOrder, [][]string) {
	mungeData := loadMungeComponents()
	allOrders := loadAllOrdersAsStruct()
	orderLists := make([][]HbOrder, len(mungeData.PlatformList))
	titleLists := make([][]string, len(mungeData.PlatformList))
	for index, platform := range mungeData.PlatformList {
		for _, order := range allOrders {
			for _, subproduct := range order.SubProducts {
				for _, download := range subproduct.Downloads {
					if platform == download.Platform {
						titleLists[index] = append(
							titleLists[index],
							subproduct.HumanName,
						)
						freshOrder := order
						freshSubproduct := subproduct
						freshDownload := download
						freshSubproduct.Downloads = []HbSubProductDownload{freshDownload}
						freshOrder.SubProducts = []HbSubProduct{freshSubproduct}
						orderLists[index] = append(
							orderLists[index],
							freshOrder,
						)
					}
				}
			}
		}
	}
	return mungeData, orderLists, titleLists
}

func RunUi() {
	mungeData, orderLists, _ := loadUiData()
	err := ui.Init()
	fatalCheck(err)
	defer ui.Close()
	x, y := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = "Press q to quit. Use 'h' and 'l' to switch tabs."
	header.SetRect(0, 0, x, 1)
	header.Border = false

	navPane := widgets.NewTabPane(mungeData.PlatformList...)
	navPane.SetRect(0, 1, x, 4)

	platformPanes := make([]*widgets.List, len(mungeData.PlatformList))
	for index, platform := range mungeData.PlatformList {
		platformPanes[index] = widgets.NewList()
		platformPanes[index].Title = platform
		platformPanes[index].Rows = orderLists[index].
		platformPanes[index].SelectedRowStyle.Bg = ui.ColorBlue
		platformPanes[index].SetRect(0, 4, x, y)
	}

	renderTab := func() {
		ui.Render(platformPanes[navPane.ActiveTabIndex])
	}

	ui.Render(
		header,
		navPane,
	)
	renderTab()
	previousKey := ""
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "h":
			navPane.FocusLeft()
			ui.Clear()
			ui.Render(header, navPane)
			renderTab()
		case "l":
			navPane.FocusRight()
			ui.Clear()
			ui.Render(header, navPane)
			renderTab()
		case "j", "<Down>":
			platformPanes[navPane.ActiveTabIndex].ScrollDown()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "k", "<Up>":
			platformPanes[navPane.ActiveTabIndex].ScrollUp()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "<C-d>":
			platformPanes[navPane.ActiveTabIndex].ScrollHalfPageDown()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "<C-u>":
			platformPanes[navPane.ActiveTabIndex].ScrollHalfPageUp()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "<C-f>":
			platformPanes[navPane.ActiveTabIndex].ScrollPageDown()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "<C-b>":
			platformPanes[navPane.ActiveTabIndex].ScrollPageUp()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "g":
			if previousKey == "g" {
				platformPanes[navPane.ActiveTabIndex].ScrollTop()
			}
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "<Home>":
			platformPanes[navPane.ActiveTabIndex].ScrollTop()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		case "G", "<End>":
			platformPanes[navPane.ActiveTabIndex].ScrollBottom()
			ui.Render(platformPanes[navPane.ActiveTabIndex])
		}
		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
	}
}

func RunUiOld() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	header := widgets.NewParagraph()
	header.Text = "Press q to quit, Press h or l to switch tabs"
	header.SetRect(0, 0, 50, 1)
	header.Border = false
	header.TextStyle.Bg = ui.ColorBlue

	p2 := widgets.NewParagraph()
	p2.Text = "Press q to quit\nPress h or l to switch tabs\n"
	p2.Title = "Keys"
	p2.SetRect(5, 5, 40, 15)
	p2.BorderStyle.Fg = ui.ColorYellow

	bc := widgets.NewBarChart()
	bc.Title = "Bar Chart"
	bc.Data = []float64{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bc.SetRect(5, 5, 35, 10)
	bc.Labels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}

	tabpane := widgets.NewTabPane("pierwszy", "drugi", "trzeci", "żółw", "four", "five")
	tabpane.SetRect(0, 1, 50, 4)
	tabpane.Border = true

	renderTab := func() {
		switch tabpane.ActiveTabIndex {
		case 0:
			ui.Render(p2)
		case 1:
			ui.Render(bc)
		}
	}

	ui.Render(header, tabpane, p2)

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "h":
			tabpane.FocusLeft()
			ui.Clear()
			ui.Render(header, tabpane)
			renderTab()
		case "l":
			tabpane.FocusRight()
			ui.Clear()
			ui.Render(header, tabpane)
			renderTab()
		}
	}
}
func RunUiNew() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	sinFloat64 := (func() []float64 {
		n := 400
		data := make([]float64, n)
		for i := range data {
			data[i] = 1 + math.Sin(float64(i)/5)
		}
		return data
	})()

	sl := widgets.NewSparkline()
	sl.Data = sinFloat64[:100]
	sl.LineColor = ui.ColorCyan
	sl.TitleStyle.Fg = ui.ColorWhite

	slg := widgets.NewSparklineGroup(sl)
	slg.Title = "Sparkline"

	lc := widgets.NewPlot()
	lc.Title = "braille-mode Line Chart"
	lc.Data = append(lc.Data, sinFloat64)
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorYellow

	gs := make([]*widgets.Gauge, 3)
	for i := range gs {
		gs[i] = widgets.NewGauge()
		gs[i].Percent = i * 10
		gs[i].BarColor = ui.ColorRed
	}

	ls := widgets.NewList()
	ls.Rows = []string{
		"[1] Downloading File 1",
		"",
		"",
		"",
		"[2] Downloading File 2",
		"",
		"",
		"",
		"[3] Uploading File 3",
	}
	ls.Border = false

	p := widgets.NewParagraph()
	p.Text = "<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget"
	p.Title = "Demonstration"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, slg),
			ui.NewCol(1.0/2, lc),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/4, ls),
			ui.NewCol(1.0/4,
				ui.NewRow(.9/3, gs[0]),
				ui.NewRow(.9/3, gs[1]),
				ui.NewRow(1.2/3, gs[2]),
			),
			ui.NewCol(1.0/2, p),
		),
	)

	ui.Render(grid)

	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			if tickerCount == 100 {
				return
			}
			for _, g := range gs {
				g.Percent = (g.Percent + 3) % 100
			}
			slg.Sparklines[0].Data = sinFloat64[tickerCount : tickerCount+100]
			lc.Data[0] = sinFloat64[2*tickerCount:]
			ui.Render(grid)
			tickerCount++
		}
	}
}
