package tokenizer

import (
    "kimchi/token"
)

type Tokenizer struct {
    input string
    position int
    peekPosition int
    char byte
}

// ==============
// Public methods
// ==============
func New(input string) *Tokenizer {
    tokenizer := &Tokenizer{input: input}
    tokenizer.readChar()

    return tokenizer
}
func (self *Tokenizer) GetToken() token.Token {
    self.skipWhitespace()
    
    // Comments
    if self.char == '#' {
        self.skipComment()
    } 
    // Identifiers and keywords
    if isLetter(self.char) && self.char != '_' {
        return token.NewIdentifier(self.readIdentifier())
    }
    // Numbers
    if isNumber(self.char) && self.char != '.' {
        return token.NewNumber(self.readNumber())
    }
    // Strings
    if self.char == '"' {
        return token.NewString(self.readString())
    }
    // Two char operators
    if (self.currentCharIs('<') && self.peekCharIs('=')) || (self.currentCharIs('>') && self.peekCharIs('=')) {
        char1 := self.char
        char2 := self.input[self.peekPosition]
        self.readChar()
        self.readChar()
        return token.NewTwoChar(char1, char2)
    }
    // Delimiters and operators
    token := token.NewChar(self.char)
    self.readChar()

    return token
}

// ===============
// Private methods
// ===============
func (self *Tokenizer) readChar() {
    if self.peekPosition >= len(self.input) {
        self.char = 0
    } else {
        self.char = self.input[self.peekPosition]
    }
    self.position = self.peekPosition
    self.peekPosition += 1
}
func (self *Tokenizer) currentCharIs(char byte) bool {
    return self.char == char
}
func (self *Tokenizer) peekCharIs(char byte) bool {
    return self.input[self.peekPosition] == char
}

func (self *Tokenizer) skipWhitespace() {
    for isWhitespace(self.char) {
        self.readChar()
    }
}
func (self *Tokenizer) skipComment() {
    for self.char != '\n' {
        self.readChar()
    }
    self.skipWhitespace()
    if self.char == '#' {
        self.skipComment()
    }
}
func (self *Tokenizer) readIdentifier() string {
    position := self.position
    for isLetter(self.char) || isNumber(self.char) {
        self.readChar()
    }
    return self.input[position:self.position]
}
func (self *Tokenizer) readNumber() string {
    position := self.position
    for isNumber(self.char) {
        self.readChar()
    }
    return self.input[position:self.position]
}
func (self *Tokenizer) readString() string {
    self.readChar()
    position := self.position
    for self.char != '"' {
        self.readChar()
    }
    self.readChar()
    return self.input[position:self.position-1]
}

// ==============
// Helper methods
// ==============
func isLetter(char byte) bool {
    return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}
func isNumber(char byte) bool {
    return ('0' <= char && char <= '9') || char == '.'
}
func isWhitespace(char byte) bool {
    return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
