package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/printer"
	"github.com/z7zmey/php-parser/walker"
)

var verbose bool

func init() {
    flag.BoolVar(&verbose, "verbose", false, "Show all occurrences of timed lines")
}

type instrumentingVisitor struct {
    output        *strings.Builder
    indentLevel   int
    currentLine   int
    skipChildren  bool
}

func (v *instrumentingVisitor) indent() string {
    return strings.Repeat("    ", v.indentLevel)
}

func nodeToString(n node.Node) string {
    var buf strings.Builder
    p := printer.NewPrinter(&buf)
    p.Print(n)
    return buf.String()
}

func (v *instrumentingVisitor) wrapWithTimer(code string, n node.Node) string {
    pos := n.GetPosition()
    if pos == nil {
        return code
    }

    v.currentLine = pos.StartLine
    return fmt.Sprintf(`
%s$____start_%d = microtime(true);
%s%s
%s$____end_%d = microtime(true);
%s____timerPush(%d, %q, $____start_%d, $____end_%d);`,
        v.indent(), pos.StartLine,
        v.indent(), code,
        v.indent(), pos.StartLine,
        v.indent(), pos.StartLine, strings.TrimSpace(code), pos.StartLine, pos.StartLine)
}

func (v *instrumentingVisitor) EnterNode(w walker.Walkable) bool {
    if v.skipChildren {
        return false
    }

    n, ok := w.(node.Node)
    if !ok {
        return true
    }

    switch n := n.(type) {
    case *stmt.Expression:
        code := v.wrapWithTimer(nodeToString(n.Expr), n)
        v.output.WriteString(code)
        v.skipChildren = true
        return false

    case *stmt.If:
        v.output.WriteString(fmt.Sprintf("%sif (%s) {\n", v.indent(), nodeToString(n.Cond)))
        v.indentLevel++
        return true

    case *stmt.ElseIf:
        v.output.WriteString(fmt.Sprintf("%s} elseif (%s) {\n", v.indent(), nodeToString(n.Cond)))
        return true

    case *stmt.Else:
        v.output.WriteString(fmt.Sprintf("%s} else {\n", v.indent()))
        return true

    case *stmt.For:
        forStr := fmt.Sprintf("for (%s; %s; %s)",
            expressionListToString(n.Init),
            expressionListToString(n.Cond),
            expressionListToString(n.Loop))
        v.output.WriteString(fmt.Sprintf("%s%s {\n", v.indent(), forStr))
        v.indentLevel++
        return true

    case *stmt.Foreach:
        foreachStr := fmt.Sprintf("foreach (%s as %s)",
            nodeToString(n.Expr),
            nodeToString(n.Variable))
        v.output.WriteString(fmt.Sprintf("%s%s {\n", v.indent(), foreachStr))
        v.indentLevel++
        return true

    case *stmt.While:
        v.output.WriteString(fmt.Sprintf("%swhile (%s) {\n", v.indent(), nodeToString(n.Cond)))
        v.indentLevel++
        return true

    case *stmt.Function:
        params := make([]string, len(n.Params))
        for i, param := range n.Params {
            params[i] = nodeToString(param)
        }
        funcStr := fmt.Sprintf("function %s(%s)",
            nodeToString(n.FunctionName),
            strings.Join(params, ", "))
        v.output.WriteString(fmt.Sprintf("%s%s {\n", v.indent(), funcStr))
        v.indentLevel++
        return true

    case *stmt.Return:
        code := ""
        if n.Expr != nil {
            code = v.wrapWithTimer(fmt.Sprintf("return %s;", nodeToString(n.Expr)), n)
        } else {
            code = v.wrapWithTimer("return;", n)
        }
        v.output.WriteString(code)
        v.skipChildren = true
        return false

    case *stmt.Class:
        className := nodeToString(n.ClassName)
        v.output.WriteString(fmt.Sprintf("%sclass %s {\n", v.indent(), className))
        v.indentLevel++
        return true

    case *stmt.ClassMethod:
        methodName := nodeToString(n.MethodName)
        v.output.WriteString(fmt.Sprintf("%spublic function %s() {\n", v.indent(), methodName))
        v.indentLevel++
        return true
    }

    return true
}

func (v *instrumentingVisitor) LeaveNode(w walker.Walkable) {
    if v.skipChildren {
        v.skipChildren = false
        return
    }

    n, ok := w.(node.Node)
    if !ok {
        return
    }

    switch n.(type) {
    case *stmt.If, *stmt.For, *stmt.Foreach, *stmt.While, *stmt.Function, *stmt.Class, *stmt.ClassMethod:
        v.indentLevel--
        v.output.WriteString(fmt.Sprintf("%s}\n", v.indent()))
    }
}

func (v *instrumentingVisitor) EnterChildNode(key string, w walker.Walkable) {
}

func (v *instrumentingVisitor) LeaveChildNode(key string, w walker.Walkable) {
}

func (v *instrumentingVisitor) EnterChildList(key string, w walker.Walkable) {
    return
}

func (v *instrumentingVisitor) LeaveChildList(key string, w walker.Walkable) {
}

func expressionListToString(exprs []node.Node) string {
    if exprs == nil {
        return ""
    }
    
    strs := make([]string, len(exprs))
    for i, expr := range exprs {
        strs[i] = nodeToString(expr)
    }
    return strings.Join(strs, ", ")
}

func main() {
    flag.Parse()
    args := flag.Args()

    if len(args) < 1 {
        fmt.Println("Usage: php-timers [--verbose] <file.php | --restore <dir>>")
        os.Exit(1)
    }

    if args[0] == "--restore" {
        if len(args) < 2 {
            fmt.Println("Please provide directory path for restore")
            os.Exit(1)
        }
        restoreFiles(args[1])
        return
    }

    processFile(args[0])
}

func processFile(filename string) {
    // Check if org file exists
    orgFile := strings.Replace(filename, ".php", ".__org__.php", 1)
    sourceFile := filename
    
    if _, err := os.Stat(orgFile); err == nil {
        sourceFile = orgFile
    } else {
        content, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }
        err = ioutil.WriteFile(orgFile, content, 0644)
        if err != nil {
            fmt.Printf("Error creating org file: %v\n", err)
            return
        }
    }

    src, err := ioutil.ReadFile(sourceFile)
    if err != nil {
        fmt.Printf("Error reading source file: %v\n", err)
        return
    }

    parser := php7.NewParser(src, sourceFile)
    parser.Parse()

    var output strings.Builder
    visitor := &instrumentingVisitor{
        output:      &output,
        indentLevel: 0,
    }
    
    rootNode := parser.GetRootNode()
    rootNode.Walk(visitor)

    preamble := `<?php
$____ptimers = array();
function ____timerPush($line, $code, $start, $end) {
    global $____ptimers;
    $diff = ($end - $start) * 1000; // convert to milliseconds
    
    if (!isset($____ptimers[$line]) || $GLOBALS["verbose"]) {
        $____ptimers[] = array(
            "line" => $line,
            "code" => $code,
            "time" => $diff
        );
    } else {
        $____ptimers[$line]["time"] += $diff;
    }
}
$GLOBALS["verbose"] = ` + fmt.Sprintf("%v", verbose) + `;
`

    postamble := `
echo "<!-- Results: \n";
print_r($____ptimers);
echo "\n-->\n";
`

    finalCode := preamble + output.String() + postamble

    err = ioutil.WriteFile(filename, []byte(finalCode), 0644)
    if err != nil {
        fmt.Printf("Error writing instrumented file: %v\n", err)
        return
    }
}

func restoreFiles(dir string) {
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if strings.HasSuffix(path, ".__org__.php") {
            originalPath := strings.TrimSuffix(path, ".__org__.php") + ".php"
            
            content, err := ioutil.ReadFile(path)
            if err != nil {
                fmt.Printf("Error reading org file %s: %v\n", path, err)
                return nil
            }
            
            err = ioutil.WriteFile(originalPath, content, 0644)
            if err != nil {
                fmt.Printf("Error restoring file %s: %v\n", originalPath, err)
                return nil
            }
            
            os.Remove(path)
            fmt.Printf("Restored %s\n", originalPath)
        }
        return nil
    })
}
