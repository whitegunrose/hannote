# Hannote for Neovim

Load it from your local clone while developing (lazy.nvim):

```lua
{
  dir = "~/code/hannote/nvim",
  config = function()
    require("hannote").setup()
  end,
}
```

`hannote-lsp` must be on your `PATH` (`make build`, then add `bin/`).
