package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"time"
)

const (
	fallacyScale = 3
	radioScale   = 2
	radioSpace   = 30
	radioSize    = 10
	origX        = 1024.0
	origY        = 768.0
)

var background = colornames.Firebrick

var cfg = pixelgl.WindowConfig{
	Title:     "Fallacy Quest",
	Bounds:    pixel.R(0, 0, origX, origY),
	VSync:     true,
	Resizable: true,
}

var winX = origX
var winY = origY

type fallacyProps struct {
	fullName string
	argCount int
	related  []string
}

var fallacyNames = map[string]fallacyProps{
	"hominem":       {"Argumentum ad Hominem", 1, []string{"straw", "emotion"}},
	"straw":         {"Straw Man", 2, []string{"hominem", "emotion"}},
	"emotion":       {"Appeal to Emotion", 1, []string{"hominem", "straw"}},
	"analogy":       {"Weak Analogy", 2, []string{"slippery", "accident", "authority", "popularity", "hasty"}},
	"hasty":         {"Hasty Generalization", 2, []string{"accident", "analogy", "authority", "popularity"}},
	"accident":      {"Accident", 2, []string{"hasty", "analogy", "authority", "popularity"}},
	"post":          {"Post Hoc Ergo Propter Hoc", 2, []string{"cum", "slippery", "popularity", "authority"}},
	"cum":           {"Cum Hoc Ergo Propter Hoc", 2, []string{"post", "slippery", "popularity", "authority"}},
	"slippery":      {"Slippery Slope", 2, []string{"cum", "post", "accident", "authority", "popularity"}},
	"authority":     {"Fallacious Appeal to Authority", 1, []string{"popularity", "cum", "post", "accident"}},
	"popularity":    {"Fallacious Appeal to Popularity", 1, []string{"authority", "cum", "post", "accident"}},
	"affirming":     {"Affirming the Consequent", 2, []string{"denying", "undistributed", "composition", "division"}},
	"denying":       {"Denying the Antecedent", 2, []string{"affirming", "undistributed", "composition", "division"}},
	"undistributed": {"Undistributed Middle", 2, []string{"denying", "affirming", "composition", "division"}},
	"equivocation":  {"Equivocation", 1, []string{"composition", "division"}},
	"composition":   {"Composition", 3, []string{"division", "equivocation"}},
	"division":      {"Division", 3, []string{"composition", "equivocation"}},
}

var fallacies = []fallacy{
	{name: "straw", phrases: []string{"Jane complains about", "the way I clean.", "She must want", "to be able", "to eat off the floor"}, ans: []int{1, 4}},
	{name: "hominem", phrases: []string{"Don't listen", "to Al Gore.", "He spews", "liberal propaganda."}, ans: []int{3}},
	{name: "hominem", phrases: []string{"Rush Limbaugh", "is a pompous windbag.", "Don't listen", "to him."}, ans: []int{1}},
	{name: "emotion", phrases: []string{"All guns", "need to be banned.", "Won't anyone", "think of the children?"}, ans: []int{3}},
	{name: "hominem", phrases: []string{"People", "who don't believe", "in gay marriage", "are absolute sickos,", "and shouldn't be", "taken seriously"}, ans: []int{3}},
	{name: "straw", phrases: []string{"How could", "global warming", "exist,", "it snowed", "just yesterday?"}, ans: []int{1, 3}},
	{name: "straw", phrases: []string{"Curbing violence", "in movies", "doesn't make sense.", "Do you think", "they should just make", "movies for kids?"}, ans: []int{0, 5}},
	{name: "straw", phrases: []string{"A cruise", "would be nice", "but we can't", "spend all our money", "on vacations!"}, ans: []int{0, 3}},
	{name: "straw", phrases: []string{"Why do you", "want more shoes?", "Nobody needs", "a thousand pairs of shoes!"}, ans: []int{1, 3}},
	{name: "slippery", phrases: []string{"Once", "I eat this", "chocolate,", "I will keep eating", "and won't stop."}, ans: []int{1, 3}},
	{name: "authority", phrases: []string{"The prayer", "cured", "her rheumatism.", "She said", "it did", "and who would know better than she?"}, ans: []int{5}},
	{name: "hasty", phrases: []string{"If they", "messed up your order", "you should", "stop doing business", "with them."}, ans: []int{1, 3}},
	{name: "cum", phrases: []string{"Countries that don't eat meat", "have", "less prostate cancer.", "Therefore,", "eating meat", "leads to", "prostate cancer"}, ans: []int{4, 6}},
	{name: "accident", phrases: []string{"I saw", "a teacher with their phone", "even though", "school policy", "says", "no phones in school.", "What gives?"}, ans: []int{1, 5}},
	{name: "cum", phrases: []string{"Smokers tend", "to come from", "low income", "areas.", "What is it", "about low income", "that", "makes people smoke?"}, ans: []int{5, 7}},
	{name: "hasty", phrases: []string{"American air is so polluted.", "I saw", "this one place", "in Houston", "with so much pollution."}, ans: []int{0, 2}},
	{name: "post", phrases: []string{"After my granddad", "had his", "heart attack", "his hair turned", "completely white.", "I didn't know", "a heart attack", "could cause that."}, ans: []int{2, 4}},
	{name: "authority", phrases: []string{"Gay parents", "cannot raise", "babies correctly.", "Reverend Jacob", "says that."}, ans: []int{3}},
	{name: "popularity", phrases: []string{"Being overweight", "can't be bad.", "85% of people", "are overweight,", "as a matter", "of fact."}, ans: []int{2}},
	{name: "popularity", phrases: []string{"Yawns are", "contagious.", "Ask anyone."}, ans: []int{2}},
	{name: "popularity", phrases: []string{"Caesar was", "a great dictator.", "After all", "everyone loved him."}, ans: []int{3}},
	{name: "popularity", phrases: []string{"Donald Trump", "must be", "the worst president.", "I mean,", "just look", "at how many people", "hate him."}, ans: []int{5}},
	{name: "accident", phrases: []string{"I can", "burn tires", "in my backyard", "if I want to.", "After all,", "it's a free country."}, ans: []int{1, 5}},
	{name: "slippery", phrases: []string{"They want", "to make", "it illegal to", "hit someone with his helmet?", "What's next,", "making tackling illegal?"}, ans: []int{2, 5}},
	{name: "authority", phrases: []string{"When I", "retake this", "stupid", "physiology course,", "I'll get", "an athlete", "to teach", "it to me.", "They're bound", "to know", "it."}, ans: []int{5}},
	{name: "authority", phrases: []string{"Alicia", "doesn't think", "it would be", "illegal,", "and I", "trust her."}, ans: []int{0}},
	{name: "slippery", phrases: []string{"If something", "isn't done", "soon,", "all English people", "will", "turn Muslim."}, ans: []int{3, 5}},
	{name: "equivocation", phrases: []string{"Professor Park", "can tell you", "if you are sick.", "After all,", "he is", "a doctor."}, ans: []int{5}},
	{name: "composition", phrases: []string{"Sodium", "is toxic", "and so is", "chlorine.", "Therefore,", "I refuse", "to eat", "salt,", "which is", "made of", "the two."}, ans: []int{0, 4, 8}},
	{name: "affirming", phrases: []string{"Rich people", "buy a car", "like a Mercedes or Bentley.", "You have", "a Bentley", "therefore", "you must be rich."}, ans: []int{4, 6}},
	{name: "undistributed", phrases: []string{"All hotels", "in the Southwest chain", "have elaborate lobbies.", "The Arlington", "also has a great lobby", "therefore", "it is a Southwest hotel."}, ans: []int{3, 6}},
	{name: "division", phrases: []string{"Water", "is wet.", "Therefore,", "both hydrogen", "and oxygen", "must be wet."}, ans: []int{0, 3, 4}},
	{name: "denying", phrases: []string{"If you", "are not", "21 or older", "you cannot drink.", "You are 21", "therefore", "you can drink."}, ans: []int{4, 6}},
	{name: "affirming", phrases: []string{"If Sally", "is 21 or older", "she can legally drink.", "Sally can legally drink", "therefore", "she is 21 or older"}, ans: []int{3, 5}},
	{name: "denying", phrases: []string{"If it is legal", "for Sally to drink", "then", "she is 21 or older.", "Sally cannot legally drink", "therefore", "she is under 21."}, ans: []int{4, 6}},
	{name: "affirming", phrases: []string{"If Sally", "is 21 or older", "she can legally drink.", "Sally is not 21 or older", "therefore", "she cannot legally drink."}, ans: []int{3, 5}},
	{name: "affirming", phrases: []string{"If you", "dropped out", "of college,", "you wouldn't", "make much", "money.", "Chris doesn't make much money,", "therefore", "he dropped out."}, ans: []int{6, 8}},
	{name: "equivocation", phrases: []string{"Of course", "he couldn't", "see your point.", "Dude's blind."}, ans: []int{2}},
	{name: "affirming", phrases: []string{"When James", "gets the paper", "Mr. Fields", "gives him", "a tip.", "Yesterday,", "Mr. Fields", "gave him a tip", "so he must've", "gotten the paper."}, ans: []int{7, 9}},
	{name: "equivocation", phrases: []string{"I'll tell you", "right now", "Mr. Horace,", "no daughter", "of mine", "is going to work", "at a", "strip mall."}, ans: []int{7}},
}

func resized(win *pixelgl.Window) bool {
	if win.Bounds().W() != winX || win.Bounds().H() != winY {
		winX = win.Bounds().W()
		winY = win.Bounds().H()
		return true
	} else {
		return false
	}
}

type textProps struct {
	line   float64
	deltaX float64
	deltaY float64
	txt    *text.Text
	bounds pixel.Rect
	active bool
}

type fallacy struct {
	win     *pixelgl.Window
	name    string
	ans     []int
	phrases []string
	texts   []textProps
	mask    []int
}

func (f *fallacy) calcTexts() {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	curLine := 0.0
	lineX := 0.0
	totalY := 0.0
	texts := make([]textProps, 0)
	maxY := 0.0
	space := text.New(pixel.ZV, atlas)
	fmt.Fprint(space, " ")
	spacing := space.Bounds().W()
	for i, phrase := range f.phrases {
		txt := text.New(pixel.ZV, atlas)
		if lineX > 0 {
			lineX += spacing
		}
		newLine := false
		if phrase[0] == '\n' {
			phrase = phrase[1:]
			newLine = true
		}
		fmt.Fprint(txt, phrase)
		if lineX+txt.Bounds().W() >= f.win.Bounds().W()/2-winX*200/origX {
			newLine = true
		}
		if newLine {
			curLine += 1
			lineX = 0
			for j := range texts[:i] {
				texts[j].deltaY += maxY * 1.5
			}
			totalY -= maxY
			maxY = 0
		}
		maxY = math.Max(maxY, txt.Bounds().H())
		for j := range texts[:i] {
			if texts[j].line == curLine {
				texts[j].deltaX -= (txt.Bounds().W() + spacing) * 1.5
			}
		}
		texts = append(texts, textProps{line: curLine, deltaX: lineX * 1.5, deltaY: totalY * 1.5, txt: txt})
		lineX += txt.Bounds().W()
	}
	for i, t := range texts {
		r := t.txt.Bounds()
		r.Max = r.Max.Scaled(fallacyScale)
		r = r.Moved(f.win.Bounds().Center().Sub(t.txt.Bounds().Center()).Sub(t.txt.Bounds().Size()).Add(pixel.V(t.deltaX, t.deltaY)))
		texts[i].bounds = r
	}
	f.texts = texts
}

func (f *fallacy) draw() {
	if resized(f.win) {
		f.calcTexts()
	}
	for i, v := range f.texts {
		mat := pixel.IM.Scaled(v.txt.Bounds().Center(), fallacyScale).Moved(f.win.Bounds().Center().Sub(v.txt.Bounds().Center()))
		found := false
		v.txt.Draw(f.win, mat.Moved(pixel.V(v.deltaX, v.deltaY)))
		for _, k := range f.mask {
			if i == k {
				found = true
				break
			}
		}
		if found {
			continue
		}
		var edgeColor color.Color
		if v.active { // Selected
			edgeColor = colornames.Lightblue
		}
		if v.bounds.Contains(f.win.MousePosition()) && edgeColor == nil { // Hover
			edgeColor = colornames.Blue
		} else if v.bounds.Contains(f.win.MousePosition()) { // Hover and Selected
			edgeColor = colornames.Darkblue
		}
		if edgeColor != nil {
			im := imdraw.New(nil)
			im.Color = edgeColor
			im.Push(v.bounds.Min, v.bounds.Max)
			im.Rectangle(3)
			im.Draw(f.win)
		}
	}
}

func (f *fallacy) check() {
	f.draw()
	for i, v := range f.texts {
		if v.bounds.Contains(f.win.MousePosition()) && f.win.JustPressed(pixelgl.MouseButton1) {
			f.texts[i].active = !v.active
		}
	}
}

func randFallacy(win *pixelgl.Window) fallacy {
	result := fallacies
	f := result[rand.Intn(len(result))]
	f.win = win
	f.calcTexts()
	return f
}

type button struct {
	win            *pixelgl.Window
	text           string
	edgeColor      color.Color
	unpressedColor color.Color
	pressedColor   color.Color
	rect           pixel.Rect
	pressed        bool
	justUnpressed  bool
}

func (b *button) draw() {
	if b.justUnpressed {
		b.justUnpressed = false
	}
	im := imdraw.New(nil)
	var buttonColor color.Color
	if b.win.JustPressed(pixelgl.MouseButtonLeft) && b.rect.Contains(b.win.MousePosition()) {
		b.pressed = true
	}
	if b.pressed && !b.win.Pressed(pixelgl.MouseButtonLeft) {
		b.pressed = false
		b.justUnpressed = true
	}
	if b.pressed {
		buttonColor = b.pressedColor
	} else {
		buttonColor = b.unpressedColor
	}
	im.Color = buttonColor
	im.Push(b.rect.Min, b.rect.Max)
	im.Rectangle(0)
	im.Draw(b.win)
	edge := imdraw.New(nil)
	if b.edgeColor == nil {
		return
	}
	edge.Color = b.edgeColor
	edge.Push(b.rect.Min, b.rect.Max)
	edge.Rectangle(3)
	edge.Draw(b.win)
}

func (b *button) check() bool {
	b.draw()
	return b.justUnpressed && b.rect.Contains(b.win.MousePosition()) && b.win.JustReleased(pixelgl.MouseButtonLeft)
}

func newButton(win *pixelgl.Window, rect pixel.Rect, unpressedColor color.Color, pressedColor color.Color) button {
	return button{
		win:            win,
		unpressedColor: unpressedColor,
		pressedColor:   pressedColor,
		rect:           rect,
	}
}

type radioButton struct {
	name    string
	display string
	deltaY  float64
	b       button
	pressed bool
}

type choice struct {
	win     *pixelgl.Window
	centerX float64
	centerY float64
	buttons []radioButton
}

func newChoice(win *pixelgl.Window, displays []string, names []string) choice {
	var buttons []radioButton
	for i := range displays {
		buttons = append(buttons, radioButton{name: names[i], display: displays[i], pressed: false})
	}
	c := choice{win: win, buttons: buttons}
	c.calcChoice()
	return c
}

func (c *choice) setCenter(centerX float64, centerY float64) {
	c.centerX = centerX
	c.centerY = centerY
	c.calcChoice()
}

func (c *choice) calcChoice() {
	for i := range c.buttons {
		deltaY := winY * radioSpace / origY * ((float64(len(c.buttons) - 1))/2 - float64(i))
		center := pixel.V(c.centerX, c.centerY+deltaY)
		c.buttons[i].deltaY = deltaY
		b := newButton(c.win, pixel.R(center.X-winX*radioSize/origX, center.Y-winY*radioSize/origY, center.X+winX*radioSize/origX, center.Y+winY*radioSize/origY), color.Transparent, color.Transparent)
		c.buttons[i].b = b
	}
}

func (c *choice) draw() {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	curPressed := -1
	newPressed := -1
	for i := range c.buttons {
		if c.buttons[i].pressed {
			curPressed = i
		} else if c.buttons[i].b.check() {
			newPressed = i
		}
	}
	pressed := curPressed
	if newPressed > -1 {
		pressed = newPressed
	}
	im := imdraw.New(nil)
	for i := range c.buttons {
		if i == pressed {
			c.buttons[i].pressed = true
		} else {
			c.buttons[i].pressed = false
		}
		if c.buttons[i].pressed {
			im.Color = colornames.Blue
		} else if c.buttons[i].b.pressed {
			im.Color = colornames.Lightgray
		} else {
			im.Color = color.Transparent
		}
		center := pixel.V(c.centerX, c.centerY+c.buttons[i].deltaY)
		im.Push(center)
		im.Circle(winY*radioSize/origY, 0)
		im.Color = colornames.Gray
		im.Push(center)
		im.Circle(winY*radioSize/origY, 3)
		txt := text.New(pixel.ZV, atlas)
		fmt.Fprint(txt, c.buttons[i].display)
		txt.Draw(c.win, pixel.IM.ScaledXY(txt.Bounds().Center(), pixel.V(winX*radioScale/origX, winY*radioScale/origY)).Moved(center.Add(pixel.V(txt.Bounds().W()*winX*(radioScale/2)/origX+winX*(radioSize+10)/origX, 0)).Sub(txt.Bounds().Center())))
	}
	im.Draw(c.win)
}

func shuffle(fallacySlice *[]string) {
	for i := range *fallacySlice {
		j := rand.Intn(i + 1)
		(*fallacySlice)[i], (*fallacySlice)[j] = (*fallacySlice)[j], (*fallacySlice)[i]
	}
}

func getFallacyChoices(name string) []string {
	result := []string{name}
	related := fallacyNames[name].related
	shuffle(&related)
	result = append(result, related[0], related[1])
	fallacyKeys := make([]string, len(fallacyNames))
	i := 0
	for k := range fallacyNames {
		fallacyKeys[i] = k
		i += 1
	}
	shuffle(&fallacyKeys)
	for j := range fallacyKeys {
		found := false
		for k := range result {
			if fallacyKeys[j] == result[k] {
				found = true
				break
			}
		}
		if !found {
			result = append(result, fallacyKeys[j])
			break
		}
	}
	shuffle(&result)
	return result
}

func menu(win *pixelgl.Window) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
resize:
	r := win.Bounds()
	// Title
	titlePos := pixel.V(r.W()/2, 8.5*r.H()/11)
	titleTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(titleTxt, "Fallacy Quest")
	// Start
	startPos := pixel.V(r.W()/2, 6*r.H()/11)
	startButton := newButton(win, pixel.R(startPos.X-winX*100/origX, startPos.Y-winY*50/origY, startPos.X+winX*100/origX, startPos.Y+winY*50/origY), colornames.Sandybrown, colornames.Rosybrown)
	startTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(startTxt, "Start")
	// Tutorial
	tutorialPos := pixel.V(r.W()/2, 4*r.H()/11)
	tutorialButton := newButton(win, pixel.R(tutorialPos.X-winX*100/origX, tutorialPos.Y-winY*50/origY, tutorialPos.X+winX*100/origX, tutorialPos.Y+winY*50/origY), colornames.Sandybrown, colornames.Rosybrown)
	tutorialTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(tutorialTxt, "Tutorial")
	// Quit
	quitPos := pixel.V(r.W()/2, 2*r.H()/11)
	quitButton := newButton(win, pixel.R(quitPos.X-winX*100/origX, quitPos.Y-winY*50/origY, quitPos.X+winX*100/origX, quitPos.Y+winY*50/origY), colornames.Sandybrown, colornames.Rosybrown)
	quitTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(quitTxt, "Quit")
	for !win.Closed() {
		win.Clear(background)
		// Title Draw
		titleTxt.Draw(win, pixel.IM.ScaledXY(titleTxt.Bounds().Center(), pixel.V(winX*5/origX, winY*5/origY)).Moved(titlePos.Sub(titleTxt.Bounds().Center())))
		// Start Check
		if startButton.check() {
			start(win, false)
			goto resize
		}
		startTxt.Draw(win, pixel.IM.ScaledXY(startTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(startPos.Sub(startTxt.Bounds().Center())))
		// Tutorial Check
		if tutorialButton.check() {
			start(win, true)
			goto resize
		}
		tutorialTxt.Draw(win, pixel.IM.ScaledXY(tutorialTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(tutorialPos.Sub(tutorialTxt.Bounds().Center())))
		// Quit Check
		if quitButton.check() {
			return
		}
		quitTxt.Draw(win, pixel.IM.ScaledXY(quitTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(quitPos.Sub(quitTxt.Bounds().Center())))
		win.Update()
		if resized(win) {
			goto resize
		}
	}
}

func start(win *pixelgl.Window, tutorial bool) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
reset:
	score := 0.0
	combo := 0
	count := 1
	total := 10
	if tutorial {
		total = 1
	}
	tutStep := 0
reload:
	timer := 0.0
	possibleGain := 10 + float64(combo)*score/2
	pointsScored := 0.0
	var f fallacy
	f = randFallacy(win)
	if tutorial {
		f.name = "analogy"
		f.ans = []int{0, 4}
		f.phrases = []string{"Mice", "are afraid", "of cats", "therefore", "humans", "are afraid", "of cats."}
		f.calcTexts()
	}
	choices := getFallacyChoices(f.name)
	var fallacyList []string
	for _, choice := range choices {
		fallacyList = append(fallacyList, fmt.Sprintf("%v (%d)", fallacyNames[choice].fullName, fallacyNames[choice].argCount))
	}
	c := newChoice(win, fallacyList, choices)
resize:
	correct := false
	c.setCenter(winX/3, 4*winY/5)
	f.calcTexts()
	r := win.Bounds()
	// Back Button
	back := newButton(win, pixel.R(r.Min.X, r.Max.Y-winY*100/origY, r.Min.X+winX*100/origX, r.Max.Y), colornames.Sandybrown, colornames.Rosybrown)
	backIcon := imdraw.New(nil)
	backIcon.Push(pixel.V(r.Min.X+winX*10/origX, r.Max.Y-winY*50/origY), pixel.V(r.Min.X+winX*90/origX, r.Max.Y-winY*10/origY), pixel.V(r.Min.X+winX*90/origX, r.Max.Y-winY*90/origY))
	backIcon.Polygon(0)
	// Check Button
	check := newButton(win, pixel.R(winX/2-winX*210/origX, winY/4-winY*50/origY, winX/2-winX*10/origX, winY/4+winY*50/origY), colornames.Green, colornames.Darkgreen)
	checkTxt := text.New(pixel.ZV, atlas)
	checkTxt.Clear()
	fmt.Fprint(checkTxt, "Check")
	// Skip Button
	skip := newButton(win, pixel.R(winX/2+winX*10/origX, winY/4-winY*50/origY, winX/2+winX*210/origX, winY/4+winY*50/origY), colornames.Red, colornames.Darkred)
	skipTxt := text.New(pixel.ZV, atlas)
	skipTxt.Clear()
	fmt.Fprint(skipTxt, "Skip")
	// Progress Text
	progressTxt := text.New(pixel.ZV, atlas)
	progressTxt.Clear()
	fmt.Fprintf(progressTxt, "%d/%d", count, total)
	// Score Text
	scoreTxt := text.New(pixel.ZV, atlas)
	scoreTxt.Clear()
	fmt.Fprintf(scoreTxt, "Score: %.2f", score)
	last := time.Now()
	// Tutorial Text
	tutTxt := text.New(pixel.ZV, atlas)
	tutTxt.Color = colornames.Black
	tutNext := newButton(win, pixel.R(winX/2-winX*60/origX, 10.5*winY/17-winY*25/origY, winX/2+winX*60/origX, 10.5*winY/17+winY*25/origY), colornames.Blue, colornames.Darkblue)
	tutNextTxt := text.New(pixel.ZV, atlas)
	tutNextTxt.Color = colornames.Green
	fmt.Fprint(tutNextTxt, "Next ->")
	for !win.Closed() {
		// Delta Time
		dt := time.Since(last).Seconds()
		last = time.Now()
		// Clear background
		win.Clear(background)
		// Update timer
		timer += dt
		// Choices
		c.draw()
		// Check
		if check.check() {
			for _, v := range c.buttons {
				if v.pressed {
					correct = v.name == f.name
					break
				}
			}
			var selected []int
			for i, v := range f.texts {
				if v.active {
					selected = append(selected, i)
				}
			}
			if len(f.ans) != len(selected) {
				correct = false
			} else {
				for _, v := range selected {
					found := false
					for _, k := range f.ans {
						if v == k {
							found = true
						}
					}
					if !found {
						correct = false
					}
				}
			}
			if correct { // Correct
				pointsScored = math.Max(possibleGain*(4/(timer+5)+.2), 0)
				checkTxt.Clear()
				fmt.Fprint(checkTxt, "Correct!")
				check.unpressedColor = color.Transparent
				check.pressedColor = color.Transparent
				skipTxt.Clear()
				fmt.Fprint(skipTxt, "Continue")
				skip.unpressedColor = colornames.Blue
				skip.pressedColor = colornames.Darkblue
			} else { // Incorrect
				combo = 0
				possibleGain -= possibleGain / 16
				score -= score / 16
				score = math.Max(score, 0)
				scoreTxt.Clear()
				fmt.Fprintf(scoreTxt, "Score: %.2f", score)
			}
		}
		checkTxt.Draw(win, pixel.IM.ScaledXY(checkTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(check.rect.Center().Sub(checkTxt.Bounds().Center())))
		// Skip
		if skip.check() {
			if correct {
				combo += 1
				score += pointsScored
			} else {
				combo = 0
				score -= score / 4
			}
			count += 1
			if count > total {
				retry := winScreen(win, score)
				if retry {
					goto reset
				} else {
					return
				}
			}
			goto reload
		}
		skipTxt.Draw(win, pixel.IM.ScaledXY(skipTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(skip.rect.Center().Sub(skipTxt.Bounds().Center())))
		// Progress
		progressTxt.Draw(win, pixel.IM.ScaledXY(progressTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(pixel.V(winX/2, 10*winY/11).Sub(progressTxt.Bounds().Center())))
		// Score
		scoreTxt.Draw(win, pixel.IM.ScaledXY(scoreTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(pixel.V(winX/2, winY/9).Sub(scoreTxt.Bounds().Center())))
		// Back
		if back.check() {
			return
		}
		backIcon.Draw(win)
		// Fallacies
		f.check()
		// Tutorial
		if tutorial {
			//tutCover.Draw(win)
			if tutStep <= 45 && tutNext.check() {
				tutStep += 1
			}
			if tutStep <= 45 {
				tutNextTxt.Draw(win, pixel.IM.ScaledXY(tutNextTxt.Bounds().Center(), pixel.V(winX*2/origX, winY*2/origY)).Moved(tutNext.rect.Center().Sub(tutNextTxt.Bounds().Center())))
			}
			tutTxt.Clear()
			switch tutStep {
			case 0:
				fmt.Fprint(tutTxt, "Welcome to Fallacy Quest!")
				break
			case 1:
				fmt.Fprint(tutTxt, "This is a fallacy, but which one?")
				break
			case 2:
				fmt.Fprint(tutTxt, "First, find the correct answer")
				break
			case 3:
				fmt.Fprint(tutTxt, `In this case, that's "Weak Analogy"`)
				break
			case 4:
				fmt.Fprint(tutTxt, "So now, click the circle next to the answer")
				break
			case 5:
				fmt.Fprint(tutTxt, "The number in parentheses next to the answer...")
				break
			case 6:
				fmt.Fprint(tutTxt, `...tells you how many "choices" it takes`)
				break
			case 7:
				fmt.Fprint(tutTxt, "The choices determine what the fallacy actually is")
				break
			case 8:
				fmt.Fprint(tutTxt, "So for a Weak Analogy, that would be...")
				break
			case 9:
				fmt.Fprint(tutTxt, "...the two things being analogized")
				break
			case 10:
				fmt.Fprint(tutTxt, "For an Accident it would be...")
				break
			case 11:
				fmt.Fprint(tutTxt, "...the generalization and the exceptional case")
				break
			case 12:
				fmt.Fprint(tutTxt, "Pretty easy right?")
				break
			case 13:
				fmt.Fprint(tutTxt, `Well, once you've figured out the "choices"...`)
				break
			case 14:
				fmt.Fprint(tutTxt, "...you can go ahead on click on them to select them")
				break
			case 15:
				fmt.Fprint(tutTxt, `In this case, the choices would be "Mice" and "humans"...`)
				break
			case 16:
				fmt.Fprint(tutTxt, "...since those are the things being analogized weakly")
				break
			case 17:
				fmt.Fprint(tutTxt, `So go on and click the words "Mice" and "humans" in the text below`)
				break
			case 18:
				fmt.Fprint(tutTxt, "Once you've bubbled in your answer above...")
				break
			case 19:
				fmt.Fprint(tutTxt, "...and selected your choices below...")
				break
			case 20:
				fmt.Fprint(tutTxt, `...you can check your answer by clicking on the green "Check" button`)
				break
			case 21:
				fmt.Fprint(tutTxt, "If your answer is correct, you'll win some points and move on")
				break
			case 22:
				fmt.Fprint(tutTxt, "If not, don't worry!")
				break
			case 23:
				fmt.Fprint(tutTxt, "You'll be given as many chances as you need to retry the question")
				break
			case 24:
				fmt.Fprint(tutTxt, "But if you're stuck, you can always skip the question")
				break
			case 25:
				fmt.Fprint(tutTxt, "Have fun!")
				break
			case 26:
				fmt.Fprint(tutTxt, `The following is a description of fallacies and their "choices"`)
				break
			case 27:
				fmt.Fprint(tutTxt, "Argumentum ad Hominem (1): The insult or attack")
				break
			case 28:
				fmt.Fprint(tutTxt, "Straw Man (2): The actual argument and the strawman argument")
				break
			case 29:
				fmt.Fprint(tutTxt, "Appeal to Emotion (1): The appeal to emotion")
				break
			case 30:
				fmt.Fprint(tutTxt, "Weak Analogy (2): The statements being analogized")
				break
			case 31:
				fmt.Fprint(tutTxt, "Hasty Generalization (2): The actual event and the generalization")
				break
			case 32:
				fmt.Fprint(tutTxt, "Accident (2): The generalization and the exceptional case")
				break
			case 33:
				fmt.Fprint(tutTxt, "Post Hoc Ergo Propter Hoc (2): The two events being compared")
				break
			case 34:
				fmt.Fprint(tutTxt, "Cum Hoc Ergo Propter Hoc (2): The two events being compared")
				break
			case 35:
				fmt.Fprint(tutTxt, "Slippery Slope (2): The initial event, and the slippery slope")
				break
			case 36:
				fmt.Fprint(tutTxt, "Fallacious Appeal to Authority (1): The false authority")
				break
			case 37:
				fmt.Fprint(tutTxt, "Fallacious Appeal to Popularity (1): The populace")
				break
			case 38:
				fmt.Fprint(tutTxt, "Affirming the Consequent (2): ...")
				break
			case 39:
				fmt.Fprint(tutTxt, "...The affirmed consequent, and the concluded antecedent")
				break
			case 40:
				fmt.Fprint(tutTxt, "Denying the Antecedent (2): ...")
				break
			case 41:
				fmt.Fprint(tutTxt, "...The denied antecedent, and the concluded consequent")
				break
			case 42:
				fmt.Fprint(tutTxt, "Undistributed Middle (2): The fallacious elements")
				break
			case 43:
				fmt.Fprint(tutTxt, "Equivocation (1): The ambiguous phrase")
				break
			case 44:
				fmt.Fprint(tutTxt, "Combination (3): The two addends, and the resultant")
				break
			case 45:
				fmt.Fprint(tutTxt, "Division (3): The original, and the two resultants")
			}
			tutTxt.Draw(win, pixel.IM.ScaledXY(tutTxt.Bounds().Center(), pixel.V(winX*2/origX, winY*2/origY)).Moved(pixel.V(winX/2, 11.5*winY/17).Sub(tutTxt.Bounds().Center())))
		}
		// Update
		win.Update()
		// Resize
		if resized(win) {
			goto resize
		}
	}
}

func winScreen(win *pixelgl.Window, score float64) bool {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
resize:
	congratsTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(congratsTxt, "Congratulations!")
	// Score Text
	scoreTxt := text.New(pixel.ZV, atlas)
	fmt.Fprintf(scoreTxt, "Score: %.2f", score)
	// Menu Button
	menu := newButton(win, pixel.R(winX/2-winX*210/origX, winY/5-winY*50/origY, winX/2-winX*10/origX, winY/5+winY*50/origY), colornames.Sandybrown, colornames.Rosybrown)
	menuTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(menuTxt, "Menu")
	// Replay
	replay := newButton(win, pixel.R(winX/2+winX*10/origX, winY/5-winY*50/origY, winX/2+winX*210/origX, winY/5+winY*50/origY), colornames.Green, colornames.Darkgreen)
	replayTxt := text.New(pixel.ZV, atlas)
	fmt.Fprint(replayTxt, "Replay")
	for !win.Closed() {
		win.Clear(background)
		// Congrats
		congratsTxt.Draw(win, pixel.IM.ScaledXY(congratsTxt.Bounds().Center(), pixel.V(winX*5/origX, winY*5/origY)).Moved(pixel.V(winX/2, 4*winY/5).Sub(congratsTxt.Bounds().Center())))
		// Score
		scoreTxt.Draw(win, pixel.IM.ScaledXY(scoreTxt.Bounds().Center(), pixel.V(winX*5/origX, winY*5/origY)).Moved(pixel.V(winX/2, 3*winY/5).Sub(scoreTxt.Bounds().Center())))
		// Menu
		if menu.check() {
			return false
		}
		menuTxt.Draw(win, pixel.IM.ScaledXY(menuTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(menu.rect.Center().Sub(menuTxt.Bounds().Center())))
		// Replay
		if replay.check() {
			return true
		}
		replayTxt.Draw(win, pixel.IM.ScaledXY(replayTxt.Bounds().Center(), pixel.V(winX*3/origX, winY*3/origY)).Moved(replay.rect.Center().Sub(replayTxt.Bounds().Center())))
		win.Update()
		if resized(win) {
			goto resize
		}
	}
	return false
}

func run() {
	var win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	menu(win)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pixelgl.Run(run)

}
