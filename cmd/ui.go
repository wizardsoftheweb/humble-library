package wotwhb

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

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

func downloadIcon(iconUrl string) image.Image {
	response, err := http.Get(iconUrl)
	fatalCheck(err)
	defer (func() { _ = response.Body.Close() })()
	img, _, err := image.Decode(response.Body)
	fatalCheck(err)
	return img
}

func RunUi() {
	mungeData, orderLists, titleLists := loadUiData()
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
		platformPanes[index].Rows = titleLists[index]
		platformPanes[index].SelectedRowStyle.Bg = ui.ColorBlue
		platformPanes[index].SetRect(0, 4, x/2, y)
	}
	refreshUi := func() {
		listIndex := platformPanes[navPane.ActiveTabIndex].SelectedRow
		iconUrl := orderLists[navPane.ActiveTabIndex][listIndex].SubProducts[0].Icon
		if "" != iconUrl {
			icon := widgets.NewImage(downloadIcon(iconUrl))
			icon.SetRect(x/2, 4, x, y/2)
			ui.Render(icon)
		} else {
			block := ui.NewBlock()
			block.SetRect(x/2, 4, x, y/2)
			ui.Render(block)
		}
		para := widgets.NewParagraph()
		para.Text = fmt.Sprint(orderLists[navPane.ActiveTabIndex][listIndex].SubProducts[0].Downloads[0].DownloadStruct)
		para.WrapText = true
		para.SetRect(x/2, y/2, x, y)
		ui.Render(
			platformPanes[navPane.ActiveTabIndex],
			navPane,
			para,
		)
	}

	renderTab := func() {
		refreshUi()
	}

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
			refreshUi()
		case "k", "<Up>":
			platformPanes[navPane.ActiveTabIndex].ScrollUp()
			refreshUi()
		case "<C-d>":
			platformPanes[navPane.ActiveTabIndex].ScrollHalfPageDown()
			refreshUi()
		case "<C-u>":
			platformPanes[navPane.ActiveTabIndex].ScrollHalfPageUp()
			refreshUi()
		case "<C-f>":
			platformPanes[navPane.ActiveTabIndex].ScrollPageDown()
			refreshUi()
		case "<C-b>":
			platformPanes[navPane.ActiveTabIndex].ScrollPageUp()
			refreshUi()
		case "g":
			if previousKey == "g" {
				platformPanes[navPane.ActiveTabIndex].ScrollTop()
			}
			refreshUi()
		case "<Home>":
			platformPanes[navPane.ActiveTabIndex].ScrollTop()
			refreshUi()
		case "G", "<End>":
			platformPanes[navPane.ActiveTabIndex].ScrollBottom()
			refreshUi()
		}
		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
	}
}
