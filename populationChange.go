package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xuri/excelize/v2"
	_ "github.com/xuri/excelize/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"strconv"
)

//go:embed graphics/*
var EmbeddedAssets embed.FS

var counter = 0
var demoApp GuiApp
var textWidget *widget.Text

func main() {
	ebiten.SetWindowSize(800, 500)
	ebiten.SetWindowTitle("Working Title")

	demoApp = GuiApp{AppUI: MakeUIWindow()}

	err := ebiten.RunGame(&demoApp)
	if err != nil {
		log.Fatalln("Error running User Interface Demo", err)
	}
}

func (g GuiApp) Update() error {
	//TODO finish me
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	//TODO finish me
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("Select a State", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Press Me", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(FunctionNameHere),
	)
	rootContainer.AddChild(button)
	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}

	//allStudents := loadStudents()
	allChanges := LoadPopChange()
	//dataAsGeneric := make([]interface{}, len(allStudents))
	stateDataGeneric := make([]interface{}, len(allChanges))
	//for position, student := range allStudents {
	//	dataAsGeneric[position] = student
	//}
	for position, state := range allChanges {
		stateDataGeneric[position] = state
	}

	listWidget := widget.NewList(
		//widget.ListOpts.Entries(dataAsGeneric),
		widget.ListOpts.Entries(stateDataGeneric),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			//fullName := "%s %s"
			//fmt.Sprintf(fullName, e.(Student).FirstName, e.(Student).LastName)
			//return e.(Student).LastName
			Name := "%s"
			fmt.Sprintf(Name, e.(PopChange).StateName)
			return e.(PopChange).StateName + " Population Change 2020: " + e.(PopChange).Change2020 +
				" Population Change 2021: " + e.(PopChange).Change2021
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			//do something when a list item changes
			e := args.Entry
			fmt.Println(e)
			pop2021, _ := strconv.ParseFloat(e.(PopChange).Population2021, 64)
			change2021, _ := strconv.ParseFloat(e.(PopChange).Change2021, 64)
			//textWidget.Label = "entry selected is " + e.(Student).LastName
			percentS := fmt.Sprintf(e.(PopChange).StateName+"'s  percentage change in population was: %f%s",
				loadPercentChange(pop2021, change2021), "%")
			textWidget.Label = percentS
		}))
	rootContainer.AddChild(listWidget)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)

	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func FunctionNameHere(args *widget.ButtonClickedEventArgs) {
	counter++
	presses := fmt.Sprintf("You have pressed the button %d times", counter)
	messages := []string{"reiufho", "woierhfowur", "uerhgiubb", "iroijjrgoijr", "98897983298",
		"843893298767&^%&GI^(", "587y935iuhh8tqg3t", "+_)+_((*&^%$#$#@%&", ":<L<?<L:{{PI)(&*^&^^%$^^&*G*",
		"ygyugyugug8yg8yg8ygygy8o", "87y*&T*YG*87g867&T^R^%E%&$E&RE%^#%$%$", "_+_+-=-00", "::::R^D^^RDR^D",
		"%EIGTI*J&^NR66b5rb865r5eg754e", "QQQQQQQQQQQQQQQQQQQQQ", "<>?<?><L:<:L<:L<:LMOJ",
		"and back to Working Title", "and back to Working Title", "and back to Working Title",
		"and back to Working Title", "and back to Working Title", "and back to Working Title",
		presses, "Having fun there?", "You should know this message is quite rare", "87878ogo8yuv",
		"8uo8y8ytythlslslslslhyguygyujjhbbbbkbbi",
		"LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG",
		"988989888u98988998798798798798798", "777", "BAR", "PICTURE 3 CHERRIES", "uguuygg76t7",
		"again again!", "still at it?", "gyuuygyguyugo8ug8y8oyy87o", "['{:;'}']"}
	message := "Don't listen to the button, he has no power here."
	textWidget.Label = message
	ebiten.SetWindowTitle(messages[rand.Intn(len(messages)-1)])

}

func LoadPopChange() []PopChange {
	excelFile2, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	all_rows2, err := excelFile2.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatal(err)
	}
	allState := make([]PopChange, 51)
	count := 0
	for number2, row2 := range all_rows2 {

		if number2 < 50 {
			continue
		}
		if row2[5] == row2[6] {
			allState[count] = PopChange{StateName: row2[6],
				Population2020: row2[8],
				Population2021: row2[9],
				Change2020:     row2[10],
				Change2021:     row2[11]}
			count++
			continue
		}

	}
	return allState
}

func loadPercentChange(old float64, change float64) float64 {
	newPop := old + change
	percent := ((old - newPop) / old) * 100
	return percent
}
