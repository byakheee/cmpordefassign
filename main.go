package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	logger      *slog.Logger
	hasProblem  bool
	showVersion *bool
	version     = "v0.3.1"
	ignoreList  = []string{}
)

func init() {
	verbose := flag.Bool("v", false, "Enable verbose logging")
	showVersion = flag.Bool("version", false, "Prints the version of the program")
	ignoreDirs := flag.String("ignore", "", "Comma-separated list of directories to ignore (matched by prefix)")
	flag.Parse()

	var handler *slog.TextHandler
	// handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	if *verbose {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	logger = slog.New(handler)

	if *ignoreDirs != "" {
		logger.Debug("Set ignoring flag.", "dirs", *ignoreDirs)
		ignoreList = strings.Split(*ignoreDirs, ",")
	}
}

func main() {
	paths := flag.Args()
	if *showVersion {
		fmt.Println("cmpordefassign ", version)
		os.Exit(0)
	}
	if len(paths) < 1 {
		logger.Error("Usage: cmporlinter <path> [<path> ...]")
		os.Exit(2)
	}

	for _, path := range paths {
		analyzePath(path)
	}
	if hasProblem {
		os.Exit(1)
	}
}

func analyzePath(path string) {
	logger.Debug("Analyzing path.", "path", path)
	if strings.HasSuffix(path, "/...") {
		basePath := strings.TrimSuffix(path, "/...")
		logger.Debug("Found './...'. Walking the path.", "path", basePath)
		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				analyzeFile(path)
			}
			return nil
		})
		if err != nil {
			logger.Error("Error walking the path.", "path", basePath, "error", err)
		}
	} else {
		info, err := os.Stat(path)
		if err != nil {
			logger.Error("Error accessing path.", "path", path, "error", err)
			return
		}
		if info.IsDir() {
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(path, ".go") {
					analyzeFile(path)
				}
				return nil
			})
			if err != nil {
				logger.Error("Error walking the directory.", "path", path, "error", err)
			}
		} else {
			analyzeFile(path)
		}
	}
}

func analyzeFile(filePath string) {
	for _, ignore := range ignoreList {
		if strings.HasPrefix(filePath, ignore) {
			logger.Debug("Ignoring file.", "file", filePath)
			return
		}
	}
	logger.Debug("Analyzing file.", "file", filePath)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", filePath, err)
		return
	}

	// ゼロ値で初期化された変数を記録するマップ
	zeroValueVars := make(map[string]struct{})

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt:
			logger.Debug("Found if statement.", "line", fset.Position(n.Pos()).Line)

			// if 文が二項式か、さらにオペレーターが '!=' かどうかを確認
			binExpr, ok := n.Cond.(*ast.BinaryExpr)
			if !ok || binExpr.Op != token.NEQ {
				return true
			}
			logger.Debug("Found '!=' in if statement.")

			// `!=` 演算子の右辺がnilもしくはゼロ値かどうかをチェック
			if !(isNil(binExpr.Y) || isZeroValue(binExpr.Y)) {
				return true
			}
			logger.Debug("Found nil or zero value in if statement.")

			// if文のスコープ内で宣言された変数の名前を記録するマップ
			declaredVars := make(map[string]struct{})

			// 変数の再代入をチェック
			for _, stmt := range n.Body.List {
				// 変数宣言をチェック
				declStmt, ok := stmt.(*ast.DeclStmt)
				if ok {
					genDecl, ok := declStmt.Decl.(*ast.GenDecl)
					if ok && genDecl.Tok == token.VAR {
						for _, spec := range genDecl.Specs {
							valueSpec := spec.(*ast.ValueSpec)
							for _, name := range valueSpec.Names {
								declaredVars[name.Name] = struct{}{}
							}
						}
					}
				}
				// 変数の再代入をチェック
				assignStmt, ok := stmt.(*ast.AssignStmt)
				if !ok || assignStmt.Tok != token.ASSIGN {
					// 代入文でない場合はスキップ
					continue
				}
				for _, lhs := range assignStmt.Lhs {
					ident, ok := lhs.(*ast.Ident)
					if !ok {
						// 代入先が変数でない場合はスキップ
						continue
					}
					if _, ok := declaredVars[ident.Name]; ok {
						// この変数はif文のスコープ内で宣言されているため、エラー条件から除外する
						continue
					}
					if _, ok := zeroValueVars[ident.Name]; ok {
						// ゼロ値で初期化された変数はエラー条件から除外する
						continue
					}

					// 変数が再代入されているため、cmp.Orを使用することが推奨される
					hasProblem = true
					pos := fset.Position(n.Pos())
					fmt.Printf(
						"%s:%d:%d: Variable: %s can be assigned in cmp.Or (cmpordefassign)\n",
						pos.Filename,
						pos.Line,
						pos.Column,
						ident.Name)
				}

			}
		case *ast.GenDecl:
			// ゼロ値で初期化されている変数を記録
			logger.Debug("Found general declaration.", "line", fset.Position(n.Pos()).Line)
			if n.Tok == token.VAR {
				for _, spec := range n.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					for _, name := range valueSpec.Names {
						if valueSpec.Values == nil {
							// varでゼロ値に初期化されている変数を記録
							zeroValueVars[name.Name] = struct{}{}
						}
					}
				}
			}
		case *ast.AssignStmt:
			// 代入文でゼロ値に初期化されている変数を記録
			logger.Debug("Found assignment statement.", "line", fset.Position(n.Pos()).Line)
			if n.Tok == token.DEFINE {
				for i, rhs := range n.Rhs {
					if isNil(rhs) || isZeroValue(rhs) {
						ident, ok := n.Lhs[i].(*ast.Ident)
						if ok {
							zeroValueVars[ident.Name] = struct{}{}
						}
					}
				}
			}
		}
		return true
	})
}

func isNil(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

func isZeroValue(expr ast.Expr) bool {
	basicLit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	isZeroValue := false
	switch basicLit.Kind {
	case token.INT:
		isZeroValue = basicLit.Value == "0"
	case token.STRING:
		isZeroValue = basicLit.Value == `""`
	case token.CHAR:
		isZeroValue = basicLit.Value == `'\x00'` // Goでは'\x00'はcharのゼロ値
	case token.FLOAT:
		isZeroValue = basicLit.Value == "0" || basicLit.Value == "0.0"
	case token.IMAG:
		isZeroValue = basicLit.Value == "0i"
	case token.IDENT:
		isZeroValue = basicLit.Value == "false"
	}
	return isZeroValue
}
