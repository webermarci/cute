# (exe)Cute

## Sync mode

```bash
cute "sleep 1; echo hello" "echo world"
```

```json
{"command":"sleep 1; echo hello","output":"hello","duration":1005025792}
{"command":"echo world","output":"world","duration":6070208}
```

## Verbose

```bash
cute -v "sleep 1; echo hello" "echo world"
```

```
2022/06/07 09:58:23 "sleep 1; echo hello" returned hello in 1.010836541s
2022/06/07 09:58:23 "echo world" returned world in 6.791042ms
2022/06/07 09:58:23 All finished in 1.018446958s
{"command":"sleep 1; echo hello","output":"hello","duration":1010836541}
{"command":"echo world","output":"world","duration":6791042}
```

## Parallel mode

```bash
cute -p -v "sleep 1; echo hello" "echo world"
```

```
2022/06/07 09:59:28 "echo world" returned world in 3.826458ms
2022/06/07 09:59:29 "sleep 1; echo hello" returned hello in 1.012250209s
2022/06/07 09:59:29 All finished in 1.012380959s
{"command":"sleep 1; echo hello","output":"hello","duration":1012250209}
{"command":"echo world","output":"world","duration":3826458}
```

## Formats

### json

```bash
cute -f json "sleep 1; echo hello" "echo world"
```

```json
{"command":"sleep 1; echo hello","output":"hello","duration":1005025792}
{"command":"echo world","output":"world","duration":6070208}
```

### csv

```bash
cute -f csv "sleep 1; echo hello" "echo world"
```

```
command,output,duration
sleep 1; echo hello,hello,1004298084
echo world,world,6046166
```

### logfmt

```bash
cute -f logfmt "sleep 1; echo hello" "echo world"
```

```
command="sleep 1; echo hello" output="hello" duration=1009989500
command="echo world" output="world" duration=6490458
```
