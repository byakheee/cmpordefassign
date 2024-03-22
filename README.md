# cmpordefassign

cmpordefassignは、変数の再代入を減らし、cmp.Or関数を利用してデフォルト値の代入を最適化することを目的としたlinterです。

# 効果

これはエラー
```go
hoge := "init"
if v := fuga(); v != nil {
    hoge := *v
}
```

これはOK
```go
hoge := cmp.Or(fuga(), "init")
```

# Install

```
go install github.com/byakheee/cmpordefassign
```

# Usage

対象のファイルにエラーがなければ Exitcode: 0.

エラーがあれば Exitcode: 1.

入力に問題があれば Exitcode: 2.

```sh
cmpordefassign ./...
```

# Test

```
go run main.go ./examples/...
```
