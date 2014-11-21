localize
========

.strings - .csv converter

## usage

### CSV生成

ディレクトリ内から.stringsファイルを探し出して単一のcsvファイルに変換します。

```
localize -csv output.csv
```


### .strings生成

csvから.stringsに変換します。

```
localize -strings input.csv
```

出力される.stringsファイルのパスはcsv内のFileカラムの内容によって決定されます。
`*.lproj`となっている部分が各言語用の名前に展開されます。


### 備考

- エンコーディングはutf8のみ対応
- 現在のところ対応言語はjaとenのみ


