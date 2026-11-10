// Package lsp implements Hannote's language server.
//
// Planned handlers:
//   - textDocument/hover: dictionary lookup for the word under the cursor (K)
//   - workspace/executeCommand:
//     hannote.setStatus, hannote.saveEntryNote,
//     hannote.appendGloss, hannote.addCustomEntry
package lsp
