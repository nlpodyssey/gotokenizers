# GoTokenizers

GoTokenizers is a Go module implementing some of the most used
state-of-the-art natural language tokenizers.

This is a **work-in-progress** port from the original Hugging Face's
[tokenizers](https://github.com/huggingface/tokenizers), written in Rust.

The project's goal is the creation of a similar tokenization library
entirely written in Go - that is, without requiring foreign language
integrations. Without having any competitive performance ambition, we simply
aim to provide a good alternative for the Go community.

**Warning: the development is currently in its very early stage**. We are
progressively porting to Go the original Rust code, which in turn is under
active development.
We want to quickly reach a stable complete port of the library, so we are
adopting a mostly quantitative approach. That means that we are consciously
sacrificing code elegance and performance, keeping instead a stricter
parallelism between the original Rust code and the newly written Go code.

Only on a second phase, once we have a first complete library port, we expect
to focus on cleaning, refactoring, testing, and maybe even re-engineering
some components. Eventually, we will evaluate further expansions, and a
possible integration or combination with the
[spago](https://github.com/saientist/spago) library. However, as you can guess,
it's too soon to tell!

## License

GoTokenizers is licensed under the
[BSD 2-Clause "Simplified" License](https://github.com/saientist/gotokenizers/blob/master/LICENSE).

This project is strongly based on
[Hugging Face Tokenizers](https://github.com/huggingface/tokenizers),
which is licensed under the
[Apache License 2.0](https://github.com/huggingface/tokenizers/blob/master/LICENSE).
