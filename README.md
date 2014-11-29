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
	4700
	1M
	12M2

Now `resisort` can help you.

* If you have some specific count of containers you can put the resistors into, specify the count with `--containers`
* If you want to have a specific maximal count of resistor in one container, specify that with `--resistors-per-container`
* without an option specify the filename

Like this: `resisort --containers 10 foo.txt` or `resisort --resistors-per-container 3 bar.txt`

It will print out an nice scheme like this:

	In order to sort 61 resistors you can use 8 container(s) with up to 8 resistor(s) each!
	1th Container:        10Ω -        33Ω
	2th Container:        39Ω -       150Ω
	3th Container:       180Ω -       680Ω
	4th Container:       820Ω -      3.9KΩ
	5th Container:      4.7KΩ -     18.0KΩ
	6th Container:     22.0KΩ -     82.0KΩ
	7th Container:    100.0KΩ -    390.0KΩ
	8th Container:    470.0KΩ -      1.0MΩ

Enjoy!

## Wanted Features
* support for color codes

## Fun Facts
* resisort is an anagram of resistor
* I spent more time on writing this program than I would have spent if I had tried to sort them by hand.
