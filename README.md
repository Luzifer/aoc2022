# Luzifer / aoc2022

My solutions to the [Advent of Code 2022](https://adventofcode.com/2022)… At least until I no longer have the motivation to think about solutions…

| Day | Language | Description |
| :---: | :--- | :--- |
| [1](https://adventofcode.com/2022/day/1) | [Go](./day01/main.go) | Quick and dirty, not much stucture, no optimizations. Just read the input from `os.Stdin` and spit out the two solutions. |
| [1](https://adventofcode.com/2022/day/1) | [Bash](./day01/bash_solution.sh) | Had some time left for the stream and added this as a kind of troll-approach. Works but probably shouldn't be taken too serious. |
| [2](https://adventofcode.com/2022/day/2) | [Go](./day02/main.go) | Possibly "slightly" overengineered but it works and has speaking names… |
| [3](https://adventofcode.com/2022/day/3) | [Go](./day03/main.go) | Again quick and dirty, not much documented but does the trick… |
| [4](https://adventofcode.com/2022/day/4) | [Go](./day04/main.go) | One moore quick and dirty solved and missing error handling but still works… |
| [5](https://adventofcode.com/2022/day/5) | [Go](./day05/main.go), [Test](./day05/main_test.go) | Few shortcuts taken by assuming valid input, otherwise proably not that bad… |
| [6](https://adventofcode.com/2022/day/6) | [Go](./day06/main.go), [Test](./day06/main_test.go) | Quite simple, possibly more easy to implement but works fine and fast. |
| [7](https://adventofcode.com/2022/day/7) | [Go](./day07/main.go), [Test](./day07/main_test.go) | Ran into trouble with path matching but finally managed to solve it. Probably not pretty but represents the filesystem and works… |
| [8](https://adventofcode.com/2022/day/8) | [Go](./day08/main.go), [Test](./day08/main_test.go) | Quote straight forward again, maybe overengineered but hey, pretty, object-oriented and tested… |
| [9](https://adventofcode.com/2022/day/9) | [Go](./day09/main.go), [Test](./day09/main_test.go) | Shortly raged after seeing part two, thought about just quitting and calling it a day with just one star and then basically rewrote the whole solution… |
| [10](https://adventofcode.com/2022/day/10) | [Go](./day10/main.go), [Test](./day10/main_test.go) | Possibly again overengineered but I fear we have to extend that CPU later on so made it more dynamic than it needs to be and added an "rendering device" on top… |
| [11](https://adventofcode.com/2022/day/11) | [Go](./day11/main.go), [Test](./day11/main_test.go) | Part one was very much straight-forward, for part two got the hint "large numbers", switched from `int` to `uint64` and it still didn't work. Before porting all the code to `big.Int` which probably would have solved the issue looked into the [AoC Subreddit](https://www.reddit.com/r/adventofcode) and ported mathematical knowledge from a [Rust solution](https://www.reddit.com/r/adventofcode/comments/zifqmh/comment/izs6tz7/) into Go to avoid using `big.Int`… |
| [12](https://adventofcode.com/2022/day/12) | [Go](./day12/main.go), [Test](./day12/main_test.go) | I hate path-finding! I really hate it. And every year there are 1+ puzzles requiring path-finding… Would probably be solvable without A\* library but whatever. Solved is solved. |
| [13](https://adventofcode.com/2022/day/13) | [Go](./day13/main.go), [Test](./day13/main_test.go) | Maybe too complicated again (seeing a trend here) but works like a charm… |
| [14](https://adventofcode.com/2022/day/14) | [Go](./day14/main.go), [Test](./day14/main_test.go) | Did I really need to render the fields? No. Was it interesting to see? Yes. |
