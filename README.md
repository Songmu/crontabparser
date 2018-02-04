crontabparser
=======

[![Build Status](https://travis-ci.org/Songmu/crontabparser.png?branch=master)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/crontabparser/badge.png?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/crontabparser?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/crontabparser
[coveralls]: https://coveralls.io/r/Songmu/crontabparser?branch=master
[license]: https://github.com/Songmu/crontabparser/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/crontabparser

## Description

check crontabs and test time matches

## Synopsis

### Parse and Check crontab

```go
f, _ := os.Open("/path/to/crontab")
crontab, _ := crontabparser.Parse(f, false)
for _, ent := range crontab.Enties() {
    ...
}
for _, job := range crontab.Jobs() {
    ...
}
```

### Parse job and check if the shedule mathes the time or not

```go
job, _ := crontabparser.ParseJob("0 0 25 12 * echo 'Happy Holidays!'", false, nil)
if job.Schedule().Match(time.Date(2018, 12, 25, 0, 0, 0, 0, time.Local)) {
    exec.Command("sh", "-c", job.Command())
}
```

### Parse schedule

```go
sche, _ := crontabparser.ParseSchedule("0 0 25 12 *")
if sche.Match(time.Date(2018, 12, 25, 0, 0, 0, 0, time.Local)) {
    fmt.Println(":tada:")
}
```

## Author

[Songmu](https://github.com/Songmu)
