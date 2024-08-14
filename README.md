# Kabsa

A simple and lightweight tool to log key presses and WPM (Words Per Minute) in real-time.

## Table of Contents

- [Installation](#installation)
- [Motivation](#motivation)
- [How it works](#how-it-works)
- [Features](#features)

## Installation

You have 2 options to use Kabsa:

1. Use it locally with go installed on your machine, by running the following command:

```bash
sudo go run main.go
```

2. Use Docker to run the containerized version of Kabsa:

```bash
sudo docker run --rm --privileged -v /dev/input:/dev/input kabsa
```

## Motivation

I discovered MoneyType a few monthes ago, and I liked the idea and how smooth it is. I wanted to create a similar tool but with a different approach. I wanted to create a tool that logs key presses and calculates the WPM in real-time (not just in a controlled environment). I also wanted to know my way around Go, so I decided to create Kabsa

## How it works

I wanted to have a reference for my numbers so I used the same formula that MoneyType uses to calculate the WPM. You can find the formula on their [about](https://monkeytype.com/about) page:

>  wpm - total number of characters in the correctly typed words 

> (including spaces), divided by 5 and normalised to 60 seconds. 
 raw wpm - calculated just like wpm, but also includes incorrect words. 

I only can calculate the raw WPM, as I don't have a way to know if the word is correct or not. I also don't have a way to know the total number of characters in the correctly typed words. So I just calculate the raw WPM

## Features

The goal of Kabsa is to be a simple and lightweight tool to log key presses and calculate the WPM in real-time. Here are some of the features that Kabsa is intended to have:

- [x] Log key presses
- [x] Calculate the WPM in real-time
- [x] Display the WPM in real-time
- [ ] Save the WPM and cumulative characters count to a database (sqlite)
- [ ] Do some analysis on the WPM
- [ ] Display the analysis in a dashboard (Somehow XD)

## Benchmarks

I will try to do some benchmarks to see how Kabsa performs compared to monkeyType and update this section with the results
