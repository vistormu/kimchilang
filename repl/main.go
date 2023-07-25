package main

import (
    "fmt"
    "os"
    "os/user"
    "io"
    "bufio"
    "kimchi/tokenizer"
    "kimchi/parser"
    "kimchi/evaluator"
    "kimchi/object"
)

const PROMPT = ">> "

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hello %s! This is the Kimchi programming language!\n", user.Username)
    start(os.Stdin, os.Stdout)
}

func start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }

        line := scanner.Text()
        tokenizer := tokenizer.New(line)
        parser := parser.New(tokenizer)

        program := parser.ParseProgram()
        if len(parser.Errors) != 0 {
            printParserErrors(out, parser.Errors)
            continue
        }

        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
    }
}

func printParserErrors(out io.Writer, errors []string) {
    io.WriteString(out, "Parser panicked! Errors: \n")
    for _, msg := range errors {
        io.WriteString(out, "\t"+msg+"\n")
    }
}
