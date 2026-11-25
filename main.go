package main

import (
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

func main() {
	a := app.New()
	w := a.NewWindow("germ")

	// TextGrid is basically a grid of characters, which is exactly what a
	// terminal needs. Just a placeholder string for now.
	ui := widget.NewTextGrid()                  // Create a new TextGrid
	ui.SetText("Hello World from the terminal") // set text to display

	// Spin up bash inside a pseudo-terminal so it thinks it's talking to a
	// real terminal (colors, line editing and so on).
	c := exec.Command("/bin/bash")
	p, err := pty.Start(c)
	if err != nil {
		fyne.LogError("Failed to start pty", err)
		os.Exit(1)
	}

	// Don't leave a zombie bash behind when we exit.
	defer c.Process.Kill()

	// Quick smoke test: type "ls" and hit enter (\r, since that's what a
	// real Enter key sends).
	p.Write([]byte("ls\r"))

	time.Sleep(1 * time.Second)
	// Only grabs the first 1KB of output, and we don't show it anywhere yet.
	b := make([]byte, 1024)
	_, err = p.Read(b)
	if err != nil {
		fyne.LogError("Failed to read pty", err)
	}

	// Fixed-size cells so the grid doesn't stretch when the window is resized.
	w.SetContent(
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
			ui,
		),
	)

	// Blocks until the window is closed, so the deferred Kill runs after this.
	w.ShowAndRun()
}
