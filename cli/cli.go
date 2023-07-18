package main

import (
    "io"
    "io/ioutil"
    "os"
    "strings"
    "kimchi/tokenizer"
    "kimchi/parser"
    "kimchi/evaluator"
    "kimchi/object"
)

const EXTENSION = ".chi"

func main() {
    out := os.Stdout
    filename := os.Args[1]
    if len(os.Args) != 2 {
        io.WriteString(out, "Usage: kimchi run <filename>.chi\n")
        return
    }
    if !strings.HasSuffix(filename, EXTENSION) {
        io.WriteString(out, "Usage: kimchi run <filename>.chi\n")
        return
    }

    // read .chi file
    content, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        panic(err)
    }

    env := object.NewEnvironment()
    tokenizer := tokenizer.New(string(content))
    parser := parser.New(tokenizer)

    program := parser.ParseProgram()
    if len(parser.Errors) != 0 {
        printParserErrors(out, parser.Errors)
        return 
    }

    evaluated := evaluator.Eval(program, env)
    if evaluated.Type() == object.ERROR_OBJ {
        io.WriteString(out, evaluated.Inspect())
        io.WriteString(out, "\n")
    }
}

func printParserErrors(out io.Writer, errors []string) {
    io.WriteString(out, "Parser panicked! Errors: \n")
    for _, msg := range errors {
        io.WriteString(out, "\t"+msg+"\n")
    }
}
