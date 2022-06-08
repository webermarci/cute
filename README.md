# (exe)Cute
It's a simple terminal command runner tool that prints the results in json format.

## Usage

The commands file should contain lines of command groups like so:
```
first="sleep 1; echo 1" second="sleep 2; echo 2"
third="sleep 3; echo 3" fourth="echo 4"
fifth="echo 5"
```
Command groups may contain one or more commands. Command groups are executed in parallel, but commands in a group are sync.

### Basic
```bash
cute < commands.txt
```

```bash
commands.txt | cute
```

```json
{"first":1,"first_duration":1013498000,"second":2,"second_duration":2009145208}
{"fourth":4,"fourth_duration":5051916,"third":3,"third_duration":3009790541}
{"fifth":5,"fifth_duration":1741167}
```

### Verbose

```bash
cute -v < commands.txt
```

```
2022/06/08 21:28:00 fifth="echo 5" returned 5 in 1.939791ms
2022/06/08 21:28:00 fourth="echo 4" returned 4 in 2.288834ms
2022/06/08 21:28:01 first="sleep 1; echo 1" returned 1 in 1.009407042s
2022/06/08 21:28:03 third="sleep 3; echo 3" returned 3 in 3.0085735s
2022/06/08 21:28:03 second="sleep 2; echo 2" returned 2 in 2.009577041s
2022/06/08 21:28:03 All finished in 3.019255584s
```
```json
{"first":1,"first_duration":1009407042,"second":2,"second_duration":2009577041}
{"fourth":4,"fourth_duration":2288834,"third":3,"third_duration":3008573500}
{"fifth":5,"fifth_duration":1939791}
```

### Sync
You can limit the number of workers to 1, so it would run the commands one by one.

```bash
cute -v -n 1 < commands.txt
```

```
2022/06/08 21:27:31 first="sleep 1; echo 1" returned 1 in 1.005350917s
2022/06/08 21:27:33 second="sleep 2; echo 2" returned 2 in 2.013512875s
2022/06/08 21:27:33 fourth="echo 4" returned 4 in 4.173166ms
2022/06/08 21:27:36 third="sleep 3; echo 3" returned 3 in 3.011984792s
2022/06/08 21:27:36 fifth="echo 5" returned 5 in 6.530042ms
2022/06/08 21:27:36 All finished in 6.042349625s
```
```json
{"first":1,"first_duration":1005350917,"second":2,"second_duration":2013512875}
{"fourth":4,"fourth_duration":4173166,"third":3,"third_duration":3011984792}
{"fifth":5,"fifth_duration":6530042}
```
