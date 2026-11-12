# General instructions

Use English (American dialect) only.

## Info
- skill directory location: `.opencode/skills`
- ignore files and directories specified in the `.gitignore` file
- When the user mentions the root directory, treat it as the directory where this file is located.
- .work is a directory for a current working context
- docs/code repeats the code directory structure and is an additional context for generation, which complements
  the skills and agents with fine-tune items



## LLM Debug

`<llm-debug/>` is a debug tag used for debug actions
`<llm-debug>debug instructions</llm-debug>`

If the prompt contains `llm-debug=true` or `llm-debug:true`, follow the instructions in the tag;
otherwise, skip them.


## LLM Notes
Direct small notes for LLM usage, marked with `<llm-tip/>`.
