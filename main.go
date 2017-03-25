package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	c "github.com/skilstak/go-colors"
)

//--------------------------------------------

const nw = 90
const nh = 22
const nano1 = `
[0;31m                                        [1;31m NNNNNNNN        NNNNNNNN     OOOOOOOOO      !!! 
[0;31m                                        [1;31m N:::::::N       N::::::N   OO:::::::::OO   !!:!!
[0;31m                                        [1;31m N::::::::N      N::::::N OO:::::::::::::OO !:::!
[0;31m                                        [1;31m N:::::::::N     N::::::NO:::::::OOO:::::::O!:::!
[0;31mnnnn  nnnnnnnn      aaaaaaaaaaaaa       [1;31m N::::::::::N    N::::::NO::::::O   O::::::O!:::!
[0;31mn:::nn::::::::nn    a::::::::::::a      [1;31m N:::::::::::N   N::::::NO:::::O     O:::::O!:::!
[0;31mn::::::::::::::nn   aaaaaaaaa:::::a     [1;31m N:::::::N::::N  N::::::NO:::::O     O:::::O!:::!
[0;31mnn:::::::::::::::n           a::::a     [1;31m N::::::N N::::N N::::::NO:::::O     O:::::O!:::!
[0;31m  n:::::nnnn:::::n    aaaaaaa:::::a     [1;31m N::::::N  N::::N:::::::NO:::::O     O:::::O!:::!
[0;31m  n::::n    n::::n  aa::::::::::::a     [1;31m N::::::N   N:::::::::::NO:::::O     O:::::O!:::!
[0;31m  n::::n    n::::n a::::aaaa::::::a     [1;31m N::::::N    N::::::::::NO:::::O     O:::::O!!:!!
[0;31m  n::::n    n::::na::::a    a:::::a     [1;31m N::::::N     N:::::::::NO::::::O   O::::::O !!! 
[0;31m  n::::n    n::::na::::a    a:::::a     [1;31m N::::::N      N::::::::NO:::::::OOO:::::::O     
[0;31m  n::::n    n::::na:::::aaaa::::::a     [1;31m N::::::N       N:::::::N OO:::::::::::::OO  !!! 
[0;31m  n::::n    n::::n a::::::::::aa:::a    [1;31m N::::::N        N::::::N   OO:::::::::OO   !!:!!
[0;31m  nnnnnn    nnnnnn  aaaaaaaaaa  aaaa    [1;31m NNNNNNNN         NNNNNNN     OOOOOOOOO      !!! 

                      [0;36mBut seriously, learn vi (not even vim).
                         You will be a better developer
                             and systems engineer.[0m
`

const nano2 = `
[0;31m                                        [0;31m NNNNNNNN        NNNNNNNN     OOOOOOOOO      !!! 
[0;31m                                        [0;31m N:::::::N       N::::::N   OO:::::::::OO   !!:!!
[0;31m                                        [0;31m N::::::::N      N::::::N OO:::::::::::::OO !:::!
[0;31m                                        [0;31m N:::::::::N     N::::::NO:::::::OOO:::::::O!:::!
[0;31mnnnn  nnnnnnnn      aaaaaaaaaaaaa       [0;31m N::::::::::N    N::::::NO::::::O   O::::::O!:::!
[0;31mn:::nn::::::::nn    a::::::::::::a      [0;31m N:::::::::::N   N::::::NO:::::O     O:::::O!:::!
[0;31mn::::::::::::::nn   aaaaaaaaa:::::a     [0;31m N:::::::N::::N  N::::::NO:::::O     O:::::O!:::!
[0;31mnn:::::::::::::::n           a::::a     [0;31m N::::::N N::::N N::::::NO:::::O     O:::::O!:::!
[0;31m  n:::::nnnn:::::n    aaaaaaa:::::a     [0;31m N::::::N  N::::N:::::::NO:::::O     O:::::O!:::!
[0;31m  n::::n    n::::n  aa::::::::::::a     [0;31m N::::::N   N:::::::::::NO:::::O     O:::::O!:::!
[0;31m  n::::n    n::::n a::::aaaa::::::a     [0;31m N::::::N    N::::::::::NO:::::O     O:::::O!!:!!
[0;31m  n::::n    n::::na::::a    a:::::a     [0;31m N::::::N     N:::::::::NO::::::O   O::::::O !!! 
[0;31m  n::::n    n::::na::::a    a:::::a     [0;31m N::::::N      N::::::::NO:::::::OOO:::::::O     
[0;31m  n::::n    n::::na:::::aaaa::::::a     [0;31m N::::::N       N:::::::N OO:::::::::::::OO  !!! 
[0;31m  n::::n    n::::n a::::::::::aa:::a    [0;31m N::::::N        N::::::N   OO:::::::::OO   !!:!!
  nnnnnn    nnnnnn  aaaaaaaaaa  aaaa    [0;31m NNNNNNNN         NNNNNNN     OOOOOOOOO      !!! 

                      [0;36mBut seriously, learn vi (not even vim).
                         You will be a better developer
                             and systems engineer.[0m
`

//--------------------------------------------

const sw = 24
const sh = 8
const nano1s = `
       [0;31mna [0;31mNOOOO!

     [0;36mBut seriously,
learn vi (not even vim).
  You will be a better
 developer and engineer.[0m
`

const nano2s = `
       [0;31mna [1;31mNOOOO!

     [0;36mBut seriously,
learn vi (not even vim).
  You will be a better
 developer and engineer.[0m
`

//--------------------------------------------

const xsw = 9
const xsh = 3
const nano1xs = "\n[0;31mna [0;31mNOOOO![0m\n"
const nano2xs = "\n[0;31mna [1;31mNOOOO![0m\n"

//--------------------------------------------

func display() {
	ws := getwinsize()
	n1 := nano1xs
	n2 := nano2xs
	w := xsw
	h := xsh
	switch {
	case ws.Col >= 88 && ws.Row >= 22:
		n1 = nano1
		n2 = nano2
		w = nw
		h = nh
	case ws.Col > 24 && ws.Row > 6:
		n1 = nano1s
		n2 = nano2s
		w = sw
		h = sh
	}
	cws := ws.Col - uint16(w)
	rws := ws.Row - uint16(h)
	if cws > 0 {
		cws /= 2
	}
	if rws > 0 {
		rws /= 2
	}
	cpad := strings.Repeat(" ", int(cws))
	rpad := strings.Repeat("\n", int(rws))

	var buf string

	fmt.Print(c.Clear + c.CurOff + rpad)
	buf = strings.Replace(n1, "\n", "\n"+cpad, -1)
	fmt.Print(buf)
	time.Sleep(800 * time.Millisecond)

	fmt.Print(c.Clear + c.CurOff + rpad)
	buf = strings.Replace(n2, "\n", "\n"+cpad, -1)
	fmt.Print(buf)
	time.Sleep(800 * time.Millisecond)

}

type winsize struct {
	Row, Col       uint16
	Xpixel, Ypixel uint16
}

func getwinsize() winsize {
	ws := winsize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	return ws
}

func getsize() *winsize {
	ws := new(winsize)
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			fmt.Print(c.CurOn + c.Clear)
			os.Exit(0)
		}
	}()
	for {
		display()
	}
}
