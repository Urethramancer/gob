package main

import (
	"os"
	"os/exec"
	"path"

	"github.com/Urethramancer/signor/log"
	"github.com/rivo/tview"
)

const (
	buildSub      = "Set architecture and flags"
	releaseFlags  = "-ldflags \"-w -s\""
	debugFlags    = "No flags"
	x86           = "x86 32-bit"
	x8664         = "x86 64-bit"
	arm32         = "ARM 32-bit"
	arm64         = "ARM 64-bit"
	mips32        = "MIPS 32-bit"
	mips64        = "MIPS 64-bit"
	mips64le      = "MIPS 64-bit LE"
	os_darwin     = "darwin"
	os_linux      = "linux"
	os_android    = "android"
	os_windows    = "windows"
	arch_386      = "386"
	arch_amd64    = "amd64"
	arch_arm32    = "arm"
	arch_arm64    = "arm64"
	arch_mips     = "mips"
	arch_mipsle   = "mipsle"
	arch_mips64   = "mips64"
	arch_mips64le = "mips64le"
)

func main() {
	var menu, osmenu, buildmenu, linuxcpumenu *tview.List
	app := tview.NewApplication()

	menu = tview.NewList().
		AddItem("Build default", "Release build for current platform", 'b', func() {
			build(nil, output(), true)
		}).
		AddItem("Build…", "Build with custom options", 'B', func() {
			app.SetRoot(osmenu, true)
		}).
		AddItem("Clean", "Remove built executable", 'c', func() {
			rungo(nil, "clean")
		}).
		AddItem("Quit", "", 'q', func() {
			app.Stop()
		})

	var env []string
	var out string
	osmenu = tview.NewList().
		AddItem("macOS", buildSub, 'm', func() {
			env = []string{"GOOS=" + os_darwin, "GOARCH=" + arch_amd64}
			out = output() + ".macos"
			app.SetRoot(buildmenu, true)
		}).
		AddItem("Linux", buildSub, 'l', func() {
			app.SetRoot(linuxcpumenu, true)
		}).
		AddItem("Back", "Return to main menu", 'q', func() {
			app.SetRoot(menu, true)
		})

	linuxcpumenu = tview.NewList().
		AddItem(x86, "", '1', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_386}
			out = output() + ".linux.x86"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(x8664, "", '2', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_amd64}
			out = output() + ".linux.amd64"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(arm32, "", '3', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_arm32}
			out = output() + ".linux.arm32"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(arm64, "", '4', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_arm64}
			out = output() + ".linux.arm64"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(mips32, "", '5', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_mips}
			out = output() + ".linux.mips32"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(mips64, "", '6', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_mips64}
			out = output() + ".linux.mips64"
			app.SetRoot(buildmenu, true)
		}).
		AddItem(mips64le, "", '7', func() {
			env = []string{"GOOS=" + os_linux, "GOARCH=" + arch_mips64le}
			out = output() + ".linux.mips64le"
			app.SetRoot(buildmenu, true)
		}).
		AddItem("Back", "Return to OS menu", 'q', func() {
			app.SetRoot(osmenu, true)
		})

	buildmenu = tview.NewList().
		AddItem("Release", releaseFlags, 'r', func() {
			build(env, out, true)
		}).
		AddItem("Debug", debugFlags, 'd', func() {
			build(env, out, false)
		}).
		AddItem("Back", "Return to OS menu", 'q', func() {
			app.SetRoot(osmenu, true)
		})

	err := app.SetRoot(menu, true).Run()
	if err != nil {
		log.Default.Err("Could't run: %s", err.Error())
	}
}

func build(env []string, name string, release bool) {
	if release {
		rungo(env, "build", "-ldflags", "-w -s", "-o", name)
	} else {
		rungo(env, "build", "-o", name)
	}
}

func rungo(env []string, opt ...string) error {
	cmd := exec.Command("go", opt...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for _, v := range env {
		cmd.Env = append(cmd.Env, v)
	}
	return cmd.Run()
}

func output() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Default.Err("Couldn't get working directory: %s", err.Error())
		os.Exit(2)
	}

	return path.Base(cwd)
}
