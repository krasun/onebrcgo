# The One Billion Row Challenge with Go

The goal of [the challenge](https://scalabledeveloper.com/posts/one-billion-row-challenge/) is to write a program for retrieving temperature measurement values from a text file and calculating the min, mean, and max temperature per weather station. And the file has 1,000,000,000 rows.

The text file has a simple structure with one measurement value per row:

```
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
St. John's;15.2
Cracow;12.6
```

The program should print out the min, mean, and max values per station, alphabetically ordered like so:

```
{Abha=5.0/18.0/27.4, Abidjan=15.7/26.0/34.1, Abéché=12.1/29.4/35.6, Accra=14.7/26.4/33.1, Addis Ababa=2.1/16.0/24.3, Adelaide=4.1/17.3/29.7, ...}
```

The goal of the challenge is to create the fastest implementation for this task.

`onebrcgo` is a solution to [the one billion row challenge](https://scalabledeveloper.com/posts/one-billion-row-challenge/) with Go.

## Profiling

Run:

```shell
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
```

And then:

```shell
go tool pprof cpu.prof
```

Or:

```shell
go tool pprof mem.prof
```

## License

`onebrcgo` is released under [the MIT license](./LICENSE).