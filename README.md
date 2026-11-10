# Hannote (한노트)

A modal, Vim-style terminal note editor for learning Korean. Put the cursor on a
word, press `Shift+K`, and its dictionary entry appears in a side pane.

> Status: early scaffolding. Nothing works end to end yet.

## How it fits together

```
Neovim (first client) ─┐
                       ├─ LSP over stdio ─▶ hannote-lsp ─▶ analyzer (Kiwi)
Hannote TUI (later) ───┘                                ├─▶ krdict.db (dictionary, read-only)
                                                        └─▶ vocab.db  (statuses, notes, custom entries)
```

| Path | What it is |
| --- | --- |
| `cmd/hannote-import` | Builds `krdict.db` from the KRDICT XML download |
| `cmd/hannote-lsp` | Language server: hover = dictionary lookup |
| `cmd/hannote` | The custom terminal editor (later) |
| `internal/dict` | Dictionary types and `krdict.db` queries |
| `internal/analyzer` | `Analyzer` interface (Kiwi now, a hand-written one later) |
| `internal/vocab` | Word statuses, entry notes, custom entries |
| `internal/lsp` | LSP handlers |
| `internal/paths` | Where Hannote keeps its data |
| `nvim/` | Neovim plugin (word pane, input-source switching) |

## Dictionary data

Hannote does not ship dictionary data. Download the KRDICT dump yourself:

1. Go to <https://krdict.korean.go.kr/download/downloadPopup>.
2. Choose **XML**, with all items and English (영어) translations.
3. Run `go run ./cmd/hannote-import <xml-dir>` to build `krdict.db`.

Data from 한국어기초사전, © 국립국어원, licensed under
[CC BY-SA 2.0 KR](https://creativecommons.org/licenses/by-sa/2.0/kr/).

## Development

```sh
make build   # binaries in ./bin
make test
```

Design doc: kept in the Hannote project on claude.ai.

## License

Code: [MIT](LICENSE). Dictionary data you import is not part of this repo and
stays under its own license (CC BY-SA 2.0 KR).
