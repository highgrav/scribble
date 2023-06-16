# Scribble
## A tiny DSL parser library

Scribble is a trivial framework for quickly defining languages, with a focus on lexing and state machine parsing,. It 
deliberately trades performance for simplicity, making it easy to add small DSLs to your application.

### Overview
Scribble consists of `Tokenizer`, which lets you define your lexical tokens using regular expressions; and `StateMachine`, 
which you use to define transitions between states as lexical tokens are encountered. You define your `StateMachine` and 
`Tokenizer` and compile them. Once you have text to parse, you `.Tokenize(...)` the text to get a sequence of `Token`s, 
which you can then feed to `StateMachine`'s `NextState(...)` function to get the next parse state.

Scribble uses `github.com/dlclark/regexp2` for its text matching, since Go's core `regexp` doesn't allow for backtracking 
and other useful features.

### Sample Usage
Check the `./test` directory for tokenizer and parser state machine usage examples.

### Caveats
Scribble is a port of a small C# library that has been previously used in production to define and parse SQL-like query 
languages and application control languages.