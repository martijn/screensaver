package main

import (
	"math/rand"
	"time"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Clock struct {
	time       time.Time
	widget     *widgets.Paragraph
	x, y       int
	xDir, yDir int
	delay      time.Duration
}

const (
	clockWidth  = 10
	clockHeight = 3
)

// Update the clock widget and render on the ui
func (c *Clock) Draw() {
	c.widget.SetRect(c.x, c.y, c.x+clockWidth, c.y+clockHeight)
	c.widget.Text = c.time.Format("15:04:05")

	ui.Render(c.widget)
}

// Update the time field every second
func (c *Clock) Ticker(ch chan bool) {
	for {
		c.time = time.Now()
		ch <- true
		time.Sleep(time.Second)
	}
}

// Bounce the clock around and change colors every step
func (c *Clock) Mover(ch chan bool) {
	for {
		time.Sleep(c.delay)

		switch rand.Intn(2) {
		case 0:
			c.x = c.x + c.xDir
		case 1:
			c.y = c.y + c.yDir
		}

		maxX, maxY := ui.TerminalDimensions()
		if (c.x + clockWidth) >= maxX {
			c.xDir = -1
		} else if c.x <= 0 {
			c.xDir = 1
		} else if (c.y + clockHeight) >= maxY {
			c.yDir = -1
		} else if c.y <= 0 {
			c.yDir = 1
		}

		c.widget.BorderStyle.Fg = termui.StandardColors[rand.Intn(len(termui.StandardColors))]

		ch <- true
	}
}

func NewClock() *Clock {
	maxX, maxY := ui.TerminalDimensions()

	return &Clock{
		time.Now(),
		widgets.NewParagraph(),
		rand.Intn(maxX - clockWidth),
		rand.Intn(maxY - clockHeight),
		1,
		1,
		800 * time.Millisecond,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Initialize termui
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// Initialize clock
	clock := NewClock()
	redrawClock := make(chan bool)

	// Start two goroutines that both send a message to the redrawClock channel
	// when they want the clock to be redrawn
	go clock.Ticker(redrawClock)
	go clock.Mover(redrawClock)

	// Main event loop
	uiEvents := ui.PollEvents()
	for {
		// Select polls channels in all cases simultaneously
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Up>":
				if clock.delay > 10*time.Millisecond {
					clock.delay = clock.delay / 2
				}
			case "<Down>":
				if clock.delay < 10*time.Second {
					clock.delay = clock.delay * 2
				}
			case "q", "<C-c>", "<Escape>":
				return
			}
		case <-redrawClock:
			clock.Draw()
		}
	}
}
