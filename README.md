<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

<br />
<div align="center">

  <!-- Replace with a logo/banner once available -->
  <!-- <img src="docs/assets/images/hannote_banner.png" alt="Hannote Logo" width="400"> -->
  <h1 align="center">Hannote (한노트)</h1>

  <p align="center">
    A modal, Vim-style terminal note editor for learning Korean — press <code>Shift+K</code> on any word to see its dictionary entry.
    <br />
    <a href="https://github.com/whitegunrose/hannote/wiki"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/whitegunrose/hannote/issues/new?labels=bug&template=bug-report.md">Report Bug</a>
    ·
    <a href="https://github.com/whitegunrose/hannote/issues/new?labels=enhancement&template=feature-request.md">Request Feature</a>
  </p>
</div>

---

## Table of Contents

1. [About The Project](#about-the-project)
   - [Built With](#built-with)
2. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
   - [Dictionary data](#dictionary-data)
3. [Usage](#usage)
4. [Roadmap](#roadmap)
5. [Contributing](#contributing)
6. [License](#license)
7. [Contact](#contact)
8. [Acknowledgments](#acknowledgments)

---
<!-- Replace with a demo GIF once available -->
<!-- ![](https://github.com/whitegunrose/hannote/blob/main/docs/assets/demos/demo.gif) -->

## About The Project

> **Status:** early development. The repo is scaffolded; the dictionary importer and language server are being built now.

Hannote is a keyboard-driven terminal note editor for Korean learners, written in Go. It replaces the paper vocabulary notebook: write down a new word, its gloss, and the sentence you found it in — then put your cursor on any word and press `Shift+K` to look it up, just like `K` looks up a man page in Vim.

Korean words change shape as endings stack up (`일어났어요` is `일어나다` + past tense + polite ending), so Hannote analyzes the word first and then shows **every** dictionary entry it could be.

**Why Hannote?**

- 📖 `Shift+K` lookup on any word — conjugated forms included
- 🔌 Fully offline — the dictionary lives on your machine
- 📝 Your notes stay plain Markdown, readable anywhere
- 🏷️ Track each word as Unknown, Seen, Known, or Hidden, with your own notes per word
- ⌨️ Vim motions with Normal, Insert, and Visual modes — and English input outside Insert mode

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

[![Go][Go-badge]][Go-url]
[![SQLite][SQLite-badge]][SQLite-url]
[![Neovim][Neovim-badge]][Neovim-url]

> Hannote is a Go language server that any editor can talk to. Neovim is the first client; a custom Hannote TUI comes next. Word analysis is powered by [Kiwi](https://github.com/bab2min/Kiwi), and definitions come from the [Basic Korean Dictionary (한국어기초사전)](https://krdict.korean.go.kr) by the National Institute of Korean Language.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Getting Started

Follow the steps below to get Hannote running on your local machine.

### Prerequisites

- **Go 1.24+** — [Install Go](https://go.dev/dl/)

  ```sh
  go version   # verify installation
  ```

- **Neovim 0.10+** — the first editor Hannote supports
- A terminal with good Unicode support (developed on macOS with [Ghostty](https://ghostty.org))

### Installation

1. Clone the repository

   ```sh
   git clone https://github.com/whitegunrose/hannote.git
   cd hannote
   ```

2. Build the binaries into `bin/`

   ```sh
   make build
   ```

3. *(Optional)* Put them on your `PATH`

   ```sh
   export PATH="$PWD/bin:$PATH"
   ```

### Dictionary data

Hannote does not ship dictionary data — you download it once and import it.

1. Go to the [한국어기초사전 download page](https://krdict.korean.go.kr/download/downloadPopup).
2. Choose **XML**, with all items and **English (영어)** translations selected.
3. Put the files in `data/krdict/` (git ignores `data/`) and import them:

   ```sh
   hannote-import data/krdict
   ```

Hannote keeps its databases in one folder:

| File          | What it holds                                  |
|---------------|------------------------------------------------|
| `krdict.db`   | The imported dictionary (read-only)            |
| `vocab.db`    | Your word statuses, word notes, custom entries |

The folder is `~/.local/share/hannote/` by default, or `$XDG_DATA_HOME/hannote/`. Set `HANNOTE_DATA_DIR` to use a different one.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Usage

> Usage below describes the planned v1 behavior.

Notes are plain Markdown. Write the word, your gloss, and the sentence you saw it in, with the word in `[brackets]`:

```markdown
## 새 단어

주로: usually
> 먼저 한국인들은 어디서 연인을 만나냐면, [주로] 회사나 학교 이런 데서 연인을 만나.
```

| Key             | Action                                              |
|-----------------|-----------------------------------------------------|
| `K`             | Look up the word under the cursor                   |
| `K` `K`         | Jump into the word pane                             |
| `1` – `4`       | *(Word pane)* Set status: Unknown / Seen / Known / Hidden |
| `]c` / `[c`     | *(Word pane)* Next / previous candidate             |

When a word has no dictionary entry, Hannote offers to add a custom one. See [`docs/note-format.md`](docs/note-format.md) for the full note format.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Roadmap

- [x] Repo scaffold and design
- [ ] Dictionary importer (`hannote-import`)
- [ ] Lookup with Kiwi word analysis
- [ ] Language server with `K` lookup in Neovim
- [ ] Word statuses, word notes, and the word pane
- [ ] Custom Hannote TUI editor
- [ ] Grammar breakdown (`일어나다` + `-았-` + `-어요`)
- [ ] Hand-written Korean analyzer

<!-- See the [open issues](https://github.com/whitegunrose/hannote/issues) for the full list of proposed features and known bugs. -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Contributing

Contributions are welcome! If you have an idea that would make Hannote better, open an issue tagged `enhancement` or send a pull request. Bug reports are just as valuable — use the [bug report template](https://github.com/whitegunrose/hannote/issues/new?labels=bug&template=bug-report.md).

### Workflow

1. Fork the project and clone your fork
2. Create a feature branch

   ```sh
   git checkout -b feature/AmazingFeature
   ```

3. Make your changes, then check them

   ```sh
   gofmt -l .     # should print nothing
   make vet
   make test
   make build
   ```

4. Commit and push

   ```sh
   git commit -m "Add AmazingFeature"
   git push origin feature/AmazingFeature
   ```

5. Open a Pull Request against `main` describing what changed and why

### Project layout

```
cmd/
├── hannote-import/         # builds krdict.db from the KRDICT XML
├── hannote-lsp/            # language server: hover = dictionary lookup
└── hannote/                # custom TUI editor (later)
internal/
├── dict/                   # dictionary types and krdict.db queries
├── analyzer/               # Analyzer interface (Kiwi now, hand-written later)
├── vocab/                  # word statuses, word notes, custom entries
├── lsp/                    # LSP handlers
└── paths/                  # where Hannote keeps its data
nvim/                       # Neovim plugin
docs/                       # note format and other docs
```

### Things to keep in mind

- **Never commit dictionary data.** The KRDICT files and `*.db` are gitignored; small hand-made samples go in `testdata/`.
- **Notes stay plain text.** Hannote never writes analyzer output into a user's note files.
- **New analyzers** implement the `Analyzer` interface in `internal/analyzer` and must return every candidate, not just the most likely one.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## License

Distributed under the MIT License. See `LICENSE` for more information.

Dictionary data from [한국어기초사전](https://krdict.korean.go.kr), © National Institute of Korean Language (국립국어원), is not part of this repository and is licensed separately under [CC BY-SA 2.0 KR](https://creativecommons.org/licenses/by-sa/2.0/kr/).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Contact

**whitegunrose** — [@whitegunrose](https://github.com/whitegunrose)

<!-- Project Link: [https://github.com/whitegunrose/hannote](https://github.com/whitegunrose/hannote) -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Acknowledgments

Resources and tools that made Hannote possible:

- [한국어기초사전 (Basic Korean Dictionary)](https://krdict.korean.go.kr) — National Institute of Korean Language
- [Kiwi](https://github.com/bab2min/Kiwi) — Korean morphological analyzer
- [Kimchi Reader](https://kimchi-reader.app) and [Yomitan](https://github.com/yomidevs/yomitan) — popup-dictionary inspiration
- [Neovim](https://neovim.io) — the `K` lookup this whole project started from
- [Go Documentation](https://go.dev/doc/)
- [Best-README-Template](https://github.com/othneildrew/Best-README-Template) — README structure inspiration
- [Shields.io](https://shields.io) — Badge generation
- [Choose an Open Source License](https://choosealicense.com)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- MARKDOWN LINKS & BADGES -->
[contributors-shield]: https://img.shields.io/github/contributors/whitegunrose/hannote.svg?style=for-the-badge
[contributors-url]: https://github.com/whitegunrose/hannote/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/whitegunrose/hannote.svg?style=for-the-badge
[forks-url]: https://github.com/whitegunrose/hannote/network/members
[stars-shield]: https://img.shields.io/github/stars/whitegunrose/hannote.svg?style=for-the-badge
[stars-url]: https://github.com/whitegunrose/hannote/stargazers
[issues-shield]: https://img.shields.io/github/issues/whitegunrose/hannote.svg?style=for-the-badge
[issues-url]: https://github.com/whitegunrose/hannote/issues
[license-shield]: https://img.shields.io/github/license/whitegunrose/hannote.svg?style=for-the-badge
[license-url]: https://github.com/whitegunrose/hannote/blob/main/LICENSE
[Go-badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[SQLite-badge]: https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white
[SQLite-url]: https://www.sqlite.org/
[Neovim-badge]: https://img.shields.io/badge/Neovim-57A143?style=for-the-badge&logo=neovim&logoColor=white
[Neovim-url]: https://neovim.io/
