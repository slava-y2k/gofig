---
name: grammar-check-code
description: "Identify and fix typos, grammar, logical, and flow errors in text, 
    code comments, and other documents. Use this skill only when the user explicitly requests a grammar check."
---

# Typos and Grammar Check

## Purpose
Analyze text for typos and for grammar, logical, and flow errors. Provide specific, focused suggestions on how to fix each issue. Focus on clarity, correctness, and readability.

## Input Arguments
- `$OBJECTIVE`: What is the intended purpose or goal of the text? (e.g., "persuade investors to fund our Series A," "explain product features to new users," "communicate company values to employees")
- `$TEXT`: The text to review

## Process

### Step 1: Scan for errors and categorize them
Read through the text once, identifying:
- **Typo**: Mistyped, misplaced, or doubled characters; unnecessary double spaces between words; doubled words; incorrect names
- **Grammar**: Spelling, punctuation, subject-verb agreement, tense consistency, modifier placement
- **Logical**: Contradictions, unsupported claims, unclear cause-and-effect, incomplete thoughts
- **Flow**: Choppy transitions, unclear organization, redundancy, passive voice overuse, vague pronouns, awkward phrasing

Rule of thumb: a typo is a slip of the fingers (transposed or doubled characters, doubled spaces); a grammar error is the wrong word or the wrong form.

### Step 2: Fix obvious errors
Obvious errors are typos and grammar errors.
- **Typo**: fix and specify the location (file and line, or paragraph/sentence)
- **Grammar**: fix and specify the rule that was applied

If the text lives in a file, apply each fix directly in the file. If it does not, provide the corrected form in the report. Either way, report every fix individually in the output format below — do not output the entire corrected text.

### Step 3: Create Fix Suggestions
For logical and flow errors provide:
- **Location**: Where in the text (e.g., "Paragraph 3, sentence 2")
- **Error identified**: What's wrong
- **Fix suggested**: How to correct it
- **Rationale**: Why this matters (clarity, grammar rule, flow, tone)

### Step 4: Prioritize
Flag highest-impact issues first:
- Critical: Grammar or logic errors that confuse readers
- Important: Flow issues that hurt readability or persuasiveness
- Minor: Stylistic suggestions or polish

---

## Error Categories and Examples

### Typo Errors

**Transposed or Doubled Characters**
- Example error: "buisness" instead of "business"
- Fix: Correct to "business"

**Doubled Words or Spaces**
- Example error: "the the" or an unnecessary double space between words
- Fix: Remove the duplicate

**Incorrect Names**
- Example error: A file, function, or product name misspelled in the text
- Fix: Match the official spelling of the name

### Grammar Errors

**Spelling**
- Example error: "seperate" instead of "separate"
- Fix: Correct spelling to "separate"

**Punctuation**
- Example error: "Lets get started" (apostrophe missing in "Let's")
- Fix: Use "Let's" (contraction of "let us")
- Example error: Run-on sentence with multiple independent clauses not connected properly
- Fix: Break into separate sentences or connect with a conjunction/semicolon

**Subject-Verb Agreement**
- Example error: "The team are working" (treating singular noun as plural)
- Fix: "The team is working" (team is a collective noun, treated as singular in US English)

**Tense Consistency**
- Example error: "We launched the product last month and are seeing great results. Users report high satisfaction and prefer our solution." (mix of past and present)
- Fix: Keep tense consistent based on timeframe

**Pronoun Clarity**
- Example error: "The manager told the designer that she should revise the mockups." (Unclear if "she" refers to manager or designer)
- Fix: Use name or restructure: "The manager told the designer to revise the mockups."

**Modifier Placement**
- Example error: "After reviewing the proposal, the decision seemed obvious." (Who reviewed? Unclear.)
- Fix: "After reviewing the proposal, we saw the decision was obvious."

---

### Logical Errors

**Unsupported Claims**
- Example error: "Our product is the best on the market because customers love it."
- Fix: Provide evidence: "Our product has a 4.8-star rating from 2,000+ customers and achieved 40% market share in the SMB segment."

**Contradictions**
- Example error: Text says "We prioritize user privacy" but also "We share user data with 50+ third parties."
- Fix: Clarify or reconcile the statements with detail

**Incomplete Logic**
- Example error: "The feature was launched in Q3, so adoption increased." (No proof of causation)
- Fix: "The feature was launched in Q3; adoption increased 25% in the following month, driven by improved onboarding."

**Vague Claims**
- Example error: "Our solution saves time and money."
- Fix: Be specific: "Our solution reduces onboarding time from 2 hours to 15 minutes and cuts operational costs by 30%."

---

### Flow Errors

**Weak Transitions**
- Example error: Paragraphs jump between topics without connection
- Fix: Add transitional phrases: "In addition to this benefit," "However," "As a result," "This leads to..."

**Choppy Sentences**
- Example error: "We launched the product. We got great feedback. We iterated quickly. We improved the feature."
- Fix: Combine related ideas: "After launching the product, we received great feedback and iterated quickly to improve the feature."

**Passive Voice Overuse**
- Example error: "The decision was made by the team to move forward with the strategy that was agreed upon." (Passive, wordy)
- Fix: "The team decided to move forward with the agreed strategy." (Active, clearer)

**Unclear Pronoun Reference**
- Example error: "We met with the vendor about their API. It was complicated, so we decided against it." (What is "it"? The API? The vendor? The meeting?)
- Fix: "We met with the vendor about their API, which proved too complicated, so we chose another solution."

**Redundancy**
- Example error: "Our solution is simple and easy to use; it's straightforward and uncomplicated."
- Fix: "Our solution is simple and easy to use." (Remove redundant synonyms)

**Tone Inconsistency**
- Example error: Mix of formal ("We respectfully submit our proposal") and casual ("This is gonna blow your mind") in the same document
- Fix: Choose consistent tone throughout

---

## Output Format

Do NOT include the corrected text in full. Instead, provide:

**[ERROR SUMMARY]**
Total number of errors found, organized by category:
- X typo errors
- X grammar errors
- X logical errors
- X flow errors

**[FIXES BY CATEGORY]**
List all errors with fixes as bullet points. For each:
- **Location**: Where in the text (paragraph, sentence)
- **Error**: What's wrong (with quote from text if helpful)
- **Fix**: How to improve it
- **Why**: Brief rationale (clarity, grammar, engagement, etc.)

**[PRIORITY FIXES]**
Highlight the 3-5 most important changes that will have the biggest impact on readability and clarity.

**[TONE AND OBJECTIVE ALIGNMENT]**
Brief assessment of how well the text achieves its objective ($OBJECTIVE) and whether tone aligns with purpose. State whether tone adjustments are needed.

---

## Important Guidelines

- **Tone**: Use straightforward, professional language. Be encouraging about the writing.
- **Focus on clarity**: Grammar matters, but clarity is paramount. A sentence can be grammatically correct but still confusing.
- **Use primary-school language**: Explain fixes in simple terms. Don't assume the reader knows grammar terminology.
- **Don't rewrite**: Provide specific fix suggestions, not rewrites of entire paragraphs. Let the author maintain their voice.
- **Include rationale**: Explain why each fix matters. This helps the author understand the principle, not just the rule.
- **Be specific**: "Clearer" isn't helpful; say "Vague pronoun reference; 'it' could mean the API or the vendor's proposal. Change to: 'The vendor's API proved too complex.'"
- **Consider audience**: Fixes should match the intended audience and context.

---

## Checklist for Review

Use this checklist to ensure thorough review:

- [ ] Check for typos (transposed or doubled characters, doubled words, double spaces)
- [ ] Check for spelling errors (use spell-check, manual review)
- [ ] Check for punctuation issues (missing commas, apostrophes, periods)
- [ ] Verify subject-verb agreement throughout
- [ ] Check tense consistency (past, present, future should align)
- [ ] Identify vague pronouns that could be clearer
- [ ] Look for sentences that could be combined or split for better flow
- [ ] Identify passive voice; flag if overused
- [ ] Check for unsupported claims; ask "Is this proven?" or "Do we have evidence?"
- [ ] Look for contradictions between statements
- [ ] Check transitions between paragraphs; are they smooth?
- [ ] Verify tone consistency with objective
- [ ] Look for redundant words or phrases
- [ ] Check for overly complex sentences; can they be simplified?
- [ ] Verify that claims support the stated objective

---

## Examples of Effective Feedback

**Poor feedback**: "This sentence is unclear."
**Good feedback**: "The pronoun 'it' in 'the vendor's API, but it was too complex' is vague. Change to 'the vendor's API was too complex' for clarity."

**Poor feedback**: "Fix the grammar here."
**Good feedback**: "Subject-verb agreement: 'The data show' vs 'The data shows'. In formal English, 'data' is plural ('The data show'); in everyday usage, it is often treated as singular ('The data shows'). Pick one form and use it consistently."

**Poor feedback**: "This doesn't flow well."
**Good feedback**: "Choppy transitions between paragraphs. Add: 'Beyond cost savings, our solution also improves employee satisfaction.' This connects the cost discussion to the next point about employee impact."

---

## When to Suggest No Change

Not every phrase needs fixing. Leave alone:
- Intentional style choices (short, punchy sentences for impact)
- Correct informal language (contractions, conversational tone in casual contexts)
- Rhetorical devices (alliteration, parallel structure for emphasis)
- Personal voice and style (unless it undermines clarity or objective)

Focus on clarity and correctness, not perfection or style uniformity.
