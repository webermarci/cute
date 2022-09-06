# (exe)Cute
It's a simple terminal command runner tool that prints the results in json format.

## Usage

The commands file should contain lines of commands like so:
```
echo 1
echo 2
echo 3
```

### Basic
```bash
cute -f commands.txt
```

```json
{"command":"echo 1","output":"1","duration":2423291}
{"command":"echo 2","output":"2","duration":2118583}
{"command":"echo 3","output":"3","duration":2256458}
```

### Verbose

```bash
cute -v -f commands.txt
```

```
2022/09/06 14:44:02 echo 3 returned 3 in 2.256458ms
2022/09/06 14:44:02 echo 2 returned 2 in 2.118583ms
2022/09/06 14:44:02 echo 1 returned 1 in 2.423291ms
2022/09/06 14:44:02 All finished in 2.463834ms
```
```json
{"command":"echo 1","output":"1","duration":2423291}
{"command":"echo 2","output":"2","duration":2118583}
{"command":"echo 3","output":"3","duration":2256458}
```

### Sync
You can limit the number of workers to 1, so it would run the commands one by one.

```bash
cute -v -n 1 -f commands.txt
```

```
2022/09/06 14:44:54 echo 1 returned 1 in 1.858042ms
2022/09/06 14:44:54 echo 2 returned 2 in 1.833958ms
2022/09/06 14:44:54 echo 3 returned 3 in 1.95675ms
2022/09/06 14:44:54 All finished in 5.78725ms
```
```json
{"command":"echo 1","output":"1","duration":1858042}
{"command":"echo 2","output":"2","duration":1833958}
{"command":"echo 3","output":"3","duration":1956750}
```
