# General instructions

Use English (American dialect) only

## Info
- skill directory location: `.opencode/skills`
- ignore files and directories specified in the `.gitignore` file
- When the user mentions the root directory, treat it as the directory where this file is located.
- .work is a directory for a current working context
- docs/code repeats the code directory structure and it is an additional context for generation which complements
skills/agents and a fine tune items 



## LLM Debug

`<llm-debug/>` is a debug tag which used for debug action
`<llm-debug>debug instructions</llm-debug>`

If prompt contains `llm-debug=true` or `llm-debug:true` then follow instructions in the tag,
otherwise skip them. 

