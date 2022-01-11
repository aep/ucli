package ucli

import (
    "bufio"
    "fmt"
    "io"
    "unicode"
)

type token struct {
    Struct  bool
    Text    string
    Line    int
    Col     int
}

type scanner struct {
    r *bufio.Reader

    line        int
    col         int
    quoted      bool
    comment     bool
    capture     string
}

func newScanner(f io.Reader) scanner {
    return scanner {
        r : bufio.NewReader(f),
        line : 1,
        col  : 0,
        quoted : false,
    }
}

func (self *scanner) next() (token, error) {

    var capture     = ""
    var start_line  = self.line
    var start_col   = self.col

    for ;; {
        char, _, err := self.r.ReadRune();
        if err != nil {
            return token{},err
        }

        var precol = self.col
        if char  == '\n' {
            self.line += 1
            self.col = 1
        } else {
            self.col += 1
        }

        // in comment , continue until newline
        if self.comment == true {
            if char == '\n' {
                self.comment = false
            }
            continue
        }

        // in quoted, continue until unquote
        if self.quoted == true {
            if char == '"' {
                self.quoted = false
                return token {
                    Line:   start_line,
                    Col:    start_col,
                    Text:   capture,
                }, nil
            } else {
                capture += string(char)
            }
            continue
        }

        // any alphanumneric is captured
        if unicode.IsLetter(char) || unicode.IsDigit(char) {
            if len(capture) == 0 {
                start_line  = self.line
                start_col   = self.col
            }
            capture += string(char)
            continue
        }

        // non-letter terminates capture
        if len(capture) > 0 {
            self.r.UnreadRune()
            if self.col == 1 {
                self.line -= 1
            }
            self.col = precol
            return token {
                Line:   start_line,
                Col:    start_col,
                Text:   capture,
            },nil
        }

        if unicode.IsSpace(char) {
            continue
        }

        // start of comment
        if char == '#' {
            start_line  = self.line
            start_col   = self.col
            self.comment = true
            continue
        }

        // start of a quote
        if char == '"' {
            start_line  = self.line
            start_col   = self.col
            self.quoted = true
            continue
        }

        // everything else is structural
        return token {
            Struct: true,
            Line:   self.line,
            Col:    self.col,
            Text:   string(char),
        }, nil
    }
}


func Parse(f io.Reader) (r map[string]map[string]map[string]string, err error) {

    var scanner = newScanner(f)
    r = make(map[string]map[string]map[string]string)

    for ;; {

        token , err := scanner.next();
        if err != nil { if err == io.EOF { break }; return r, err }
        if token.Struct {
            return nil, fmt.Errorf("%d:%d: expected label but got %s", token.Line, token.Col, token.Text)
        }

        var root = token.Text
        var name = ""
        var rr = make(map[string]string)

        for i:=0;;i++ {
            token, err := scanner.next();
            if err != nil { return r, err }
            if token.Struct {
                if token.Text != "{" {
                    return nil, fmt.Errorf("%d:%d: expected label or { but got %s", token.Line, token.Col, token.Text)
                }
                break
            }
            if i == 0 {
                name = token.Text
            }
            rr[fmt.Sprintf("[%d]", i)] = token.Text
        }

        for ;; {
            token, err := scanner.next();
            if err != nil { return r, err }
            if token.Struct {
                if token.Text != "}" {
                    return nil, fmt.Errorf("%d:%d: expected key or } but got %s", token.Line, token.Col, token.Text)
                }
                break
            }
            key := token.Text

            token, err = scanner.next();
            if err != nil { return r, err }
            if token.Struct {
                if token.Text != "=" && token.Text != ":" {
                    return nil, fmt.Errorf("%d:%d: expected : or = or value but got %s", token.Line, token.Col, token.Text)
                }
                token, err = scanner.next();
                if err != nil { return r, err }
            }
            if token.Struct {
                return nil, fmt.Errorf("%d:%d: expected value but got %s", token.Line, token.Col, token.Text)
            }

            rr[key] = token.Text
        }

        rm := r[root]
        if rm == nil {
            rm = make(map[string]map[string]string)
        }
        if name == "" {
            name = fmt.Sprintf("[%d]", len(rm))
        }
        rm[name] = rr
        r[root] = rm

    }

    return r,nil
}
