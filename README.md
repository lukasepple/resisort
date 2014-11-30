# resisort – A resistor sorting helper tool
You bought some huge collection of resistors? Want to sort them? Don't want to spend hours? `resisort` can help you with that!

## Installation

First install [Go](http://golang.org). After setting up your `GOPATH` execute:

	go get github.com/lukasepple/resisort

## Usage

First you have to create a file of the resistor sizes you have. The sizes just have to be separated with a newline:

	100
	120Ω
	33R
	12K
	4K7
	1M

Now `resisort` can help you.

* If you have some specific count of containers you can put the resistors into, specify the count with `--containers`
* If you want to have a specific maximal count of resistor in one container, specify that with `--resistors-per-container`
* You can specify the filename with `--file`. If not stdin is used

Like this: `resisort --containers 10 --file foo.txt` or `resisort --resistors-per-container 3 --file bar.txt`

It will print out an nice scheme like this:

	In order to sort 61 resistors you can use 9 container(s) with up to 7 resistor(s) each!
	  1st Container:        10Ω -        33Ω
	  2nd Container:        39Ω -       120Ω
	  3rd Container:       150Ω -       470Ω
	  4th Container:       560Ω -      1.8KΩ
	  5th Container:      2.2KΩ -      6.8KΩ
	  6th Container:      8.2KΩ -       27KΩ
	  7th Container:       33KΩ -      100KΩ
	  8th Container:      120KΩ -      390KΩ
	  9th Container:      470KΩ -        1MΩ

Enjoy!

## Wanted Features
* support for color codes
* systemd-integration [not interely seriously]

## Fun Facts
* resisort is an anagram of resistor
* I spent more time on writing this program than I would have spent if I had tried to sort them by hand.

## Free Software now sorts your resistors!
resisort is published under the [GPLv3](./LICENSE)
