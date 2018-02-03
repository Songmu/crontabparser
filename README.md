checron
=======

[![Build Status](https://travis-ci.org/Songmu/checron.png?branch=master)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/checron?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/checron
[license]: https://github.com/Songmu/checron/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/checron

## Description

check crontabs and test time matches

## Synopsis

### Parse and Check crontab

```go
f, _ := os.Open("/path/to/crontab")
crontab, _ := checron.Parse(f, false)
for _, ent := range crontab.Enties() {
    ...
}
for _, job := range crontab.Jobs() {
    ...
}
```

### Parse job and check if the shedule mathes the time or not

```go
job, _ := checron.ParseJob("0 0 25 12 * echo 'Happy Holidays!'", false, nil)
if job.Schedule().Match(time.Date(2018, 12, 25, 0, 0, 0, 0, time.Local)) {
    exec.Command("sh", "-c", job.Command())
}
```

### Parse schedule

```go
sche, _ := checron.ParseSchedule("0 0 25 12 *")
if sche.Match(time.Date(2018, 12, 25, 0, 0, 0, 0, time.Local)) {
    fmt.Println(":tada:")
}
```

## Author

[Songmu](https://github.com/Songmu)
