# Note format

Notes are plain Markdown. Hannote never writes analyzer output into them.

```markdown
## 새 단어

주로: usually
> 먼저 한국인들은 어디서 연인을 만나냐면, [주로] 회사나 학교 이런 데서 연인을 만나.

통하다: to understand / be on the same page
> 그냥 뭔가 feeling이 [통해]야 되는데...
```

- `word: gloss` — your own gloss. A lookup can insert this line for you.
- `> sentence` — where you saw the word.
- `[ ]` — marks the word in the sentence. It can cover part of a word
  (`[통해]야`); lookup always analyzes the whole space-separated word.
