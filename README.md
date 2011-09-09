Graphpaper
==========

Graphpaper is an experiment in writing a server metrics aggregation and
graphing application.

Right now it doesn't do very much. It's just some rough code exploring the
idea of storing every data point received with no loss in resolution, and
doing all aggregation later. Most existing tools summarize incoming data
before storing it to make storage and future analysis easier. Unfortunately
this makes it impossible to work with the raw numbers, making it easy to miss
patterns in the resulting aggregated data.

# Installing

These instructions are deliberately obtuse. Graphpaper doesn't even qualify as
alpha software yet and you really shouldn't be installing it unless you know
what you're doing.

Graphpaper is written in Go. You'll need to install release 59 of Go, along
with the following packages:

  * github.com/skelterjohn/go-gb/gb
  * github.com/droundy/goopt

From there running `$GOBIN/gb` should compile everything, but don't be
surprised if it doesn't.