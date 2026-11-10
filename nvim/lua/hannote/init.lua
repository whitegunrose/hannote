-- Hannote Neovim plugin (early stub).
-- Starts hannote-lsp for Markdown notes so `K` (hover) does dictionary lookup.
-- The word pane and input-source switching come later.

local M = {}

M.config = {
  cmd = { "hannote-lsp" },
  filetypes = { "markdown" },
}

function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})

  vim.api.nvim_create_autocmd("FileType", {
    pattern = M.config.filetypes,
    callback = function(args)
      if vim.fn.executable(M.config.cmd[1]) == 0 then
        return
      end
      vim.lsp.start({
        name = "hannote",
        cmd = M.config.cmd,
        root_dir = vim.fs.dirname(args.file),
      })
    end,
  })
end

return M
