# AI Writing Pattern Catalog

The full reference behind the pattern index in SKILL.md. Read
this before any substantial rewrite. Each entry has the tell,
why it reads as generated, and a before/after pair.

"Before" specimens intentionally contain the defects under
discussion. They are exhibits, not style.

A single match is not a defect. Clusters are. Judge paragraphs,
not words.

## Table of Contents

- [Content Patterns](#content-patterns) (1-7)
- [Language and Grammar Patterns](#language-and-grammar-patterns) (8-13)
- [Style Patterns](#style-patterns) (14-17)
- [Communication Patterns](#communication-patterns) (18-20)
- [Filler and Hedging](#filler-and-hedging) (21-24)
- [Structure and Drama](#structure-and-drama) (25-28)
- [Adding Voice Without Inventing It](#adding-voice-without-inventing-it)
- [Full Worked Example](#full-worked-example)
- [Attribution](#attribution)

## Content Patterns

### 1. Inflated Significance

Watch for: "stands as", "serves as", "is a testament to",
"pivotal moment", "underscores the importance of", "reflects
broader", "enduring legacy", "setting the stage for", "marking
a shift", "shaping the future", "deeply rooted", "load-bearing",
"focal point".

LLM prose inflates ordinary facts into grand statements about
history, legacy, or trends.

Before:

> The Statistical Institute of Catalonia was officially
> established in 1989, marking a pivotal moment in the evolution
> of regional statistics in Spain. This initiative was part of a
> broader movement to decentralize administrative functions.

After:

> The Statistical Institute of Catalonia was established in 1989
> to collect and publish regional statistics independently from
> Spain's national statistics office.

### 2. Forced Notability

Watch for: "independent coverage", "featured in major
publications", "active social media presence", "written by a
leading expert".

The prose tries to prove importance by listing attention instead
of saying what happened.

Before:

> Her views have been cited in The New York Times, BBC, and The
> Hindu. She maintains an active social media presence with over
> 500,000 followers.

After:

> In a 2024 interview, she argued that AI regulation should
> focus on outcomes rather than methods.

### 3. Superficial -ing Analysis

Watch for: "highlighting", "underscoring", "emphasizing",
"ensuring", "reflecting", "symbolizing", "fostering",
"showcasing", "cultivating", "contributing to".

Present-participle phrases tacked onto sentences simulate depth
without adding any.

Before:

> The temple's color palette of blue, green, and gold resonates
> with the region's natural beauty, symbolizing Texas bluebonnets
> and the Gulf of Mexico, reflecting the community's deep
> connection to the land.

After:

> The temple uses blue, green, and gold. The architect said the
> colors reference local bluebonnets and the Gulf coast.

### 4. Brochure Language

Watch for: "boasts", "vibrant", "profound", "nestled", "in the
heart of", "renowned", "breathtaking", "stunning", "seamless",
"intuitive", "powerful", "groundbreaking" (figurative), "rich"
(figurative), "commitment to", "natural beauty".

The prose drifts into advertisement copy, especially for places,
products, and project summaries.

Before:

> Nestled within the breathtaking region of Gonder, Alamata Raya
> Kobo stands as a vibrant town with a rich cultural heritage and
> stunning natural beauty.

After:

> Alamata Raya Kobo is a town in the Gonder region of Ethiopia,
> known for its weekly market and 18th-century church.

### 5. Vague Attribution

Watch for: "experts argue", "observers have cited", "industry
reports", "some critics argue", "many believe", "it is widely
regarded", "publications have noted".

Claims get attributed to authorities nobody can check.

Before:

> Due to its unique characteristics, the Haolai River is of
> interest to researchers. Experts believe it plays a crucial
> role in the regional ecosystem.

After:

> The Haolai River supports several endemic fish species,
> according to a 2019 survey by the Chinese Academy of Sciences.

If no source exists, say less. Do not dress a guess as
consensus. And removing a fake attribution must not silently
strengthen the claim: asserting flatly what the source only
attributed to "observers" changes the certainty. Prefer saying
less, or a hedged agentless form ("tools like it are starting
to be treated as essential"), over a flat assertion.

### 6. Challenges-and-Outlook Filler

Watch for: "faces several challenges", "despite these
challenges", "future outlook", "looking ahead", "the road
ahead", "the future looks bright".

A generic closing section that could be pasted under any topic.

Before:

> Despite its industrial prosperity, Korattur faces challenges
> typical of urban areas. Despite these challenges, with its
> strategic location, Korattur continues to thrive as an integral
> part of Chennai's growth.

After:

> Traffic congestion increased after 2015 when three new IT
> parks opened. The municipal corporation began a stormwater
> drainage project in 2022 to address recurring floods.

### 7. Forced Triplets

Watch for: three-part slogans, three parallel fragments, three
abstract nouns, three examples where one would do.

> Not invention, not interaction, but execution.

> The product improves speed, quality, and collaboration.

Humans use threes too, so do not ban them. Remove them when they
feel decorative rather than load-bearing.

Before:

> The event features keynote sessions, panel discussions, and
> networking opportunities. Attendees can expect innovation,
> inspiration, and industry insights.

After:

> The event includes talks and panels, with time for informal
> networking between sessions.

## Language and Grammar Patterns

### 8. Overused AI Vocabulary

Watch for clustering of: "additionally", "align with",
"crucial", "delve", "enhance", "fostering", "garner",
"highlight" (verb), "interplay", "intricate", "key" (generic
adjective), "landscape" (abstract), "pivotal", "showcase",
"tapestry" (abstract), "testament", "underscore" (verb),
"valuable", "vibrant", "enduring", "emphasizing".

None of these words is banned. The tell is density: several in
one paragraph means the paragraph needs simplification.

Before:

> Additionally, a distinctive feature of Somali cuisine is the
> incorporation of camel meat. An enduring testament to Italian
> colonial influence is the widespread adoption of pasta in the
> local culinary landscape.

After:

> Somali cuisine also includes camel meat, which is considered a
> delicacy. Pasta dishes, introduced during Italian colonization,
> remain common, especially in the south.

### 9. Copula Avoidance

Watch for: "serves as", "stands as", "marks", "represents",
"boasts", "features", "offers".

The model dodges plain "is", "are", and "has".

Before:

> Gallery 825 serves as LAAA's exhibition space for contemporary
> art. The gallery features four separate spaces and boasts over
> 3,000 square feet.

After:

> Gallery 825 is LAAA's exhibition space for contemporary art.
> The gallery has four rooms totaling 3,000 square feet.

### 10. Negative Parallelisms and Tailing Negations

"Not only X but Y", "not just X, it is Y", and clipped negations
bolted onto sentence ends.

Before:

> It is not just about the beat riding under the vocals; it is
> part of the aggression and atmosphere. It is not merely a song,
> it is a statement.

After:

> The heavy beat adds to the aggressive tone.

Before:

> The options come from the selected item, no guessing.

After:

> The options come from the selected item without forcing the
> user to guess.

### 11. Elegant Variation

Synonym-cycling to avoid repetition when repetition would be
clearer.

Before:

> The protagonist faces many challenges. The main character must
> overcome obstacles. The central figure eventually triumphs. The
> hero returns home.

After:

> The protagonist faces many challenges but eventually triumphs
> and returns home.

### 12. False Ranges

"From X to Y" where X and Y sit on no meaningful scale.

Before:

> Our journey has taken us from the singularity of the Big Bang
> to the grand cosmic web, from the birth and death of stars to
> the enigmatic dance of dark matter.

After:

> The book covers the Big Bang, star formation, and current
> theories about dark matter.

### 13. Reflexive Passive Voice

Passive voice is not automatically wrong. Keep it when the actor
is unknown, irrelevant, or deliberately omitted ("The files are
encrypted at rest"). Rewrite when the actor matters.

Before:

> The incident was resolved after the service was restarted.

After:

> The on-call engineer restarted the service and resolved the
> incident.

## Style Patterns

### 14. Boldface Spray

Mechanically bolded phrases that direct emphasis nowhere.

Before:

> It blends **OKRs (Objectives and Key Results)**, **KPIs (Key
> Performance Indicators)**, and visual strategy tools such as
> the **Business Model Canvas (BMC)**.

After:

> It blends OKRs, KPIs, and visual strategy tools like the
> Business Model Canvas.

### 15. Inline-Header Bullets

Bullets that open with a bold label and colon when prose would
read better.

Before:

> * **User Experience:** The interface is significantly improved.
> * **Performance:** Performance is enhanced through optimized algorithms.
> * **Security:** Security is strengthened with end-to-end encryption.

After:

> The update improves the interface, speeds up load times through
> optimized algorithms, and adds end-to-end encryption.

Keep lists when the structure genuinely helps the reader. Do not
collapse a useful checklist into prose just to dodge a pattern.

### 16. Mechanical Headings

Too many headings, generic heading names, or a heading followed
by a one-line warm-up that restates the heading.

Before:

> ## Performance
>
> Speed matters.
>
> When users hit a slow page, they leave.

After:

> ## Performance
>
> When users hit a slow page, they leave.

Preserve the document's existing heading convention. Do not
impose one universal style.

### 17. Emoji Decoration

Emojis as ornaments on headings and bullets. Strip them from
prose that is not chat, and fold the content into sentences.

## Communication Patterns

### 18. Chatbot Residue

Watch for: "I hope this helps", "Great question!", "Certainly!",
"You're absolutely right", "Would you like me to", "let me
know", "let's dive in", "let's explore", "here's what you need
to know", "without further ado".

Chat-turn framing pasted into content.

Before:

> Great question! Here is an overview of the French Revolution.
> I hope this helps! Let me know if you would like me to expand
> on any section.

After:

> The French Revolution began in 1789 when financial crisis and
> food shortages led to widespread unrest.

### 19. Cutoff Disclaimers and Speculative Gap-Filling

Watch for: "as of my last update", "based on available
information", "while specific details are limited", "not
publicly available", "maintains a low profile", "likely grew
up", "it is believed that".

Two related tells: capability disclaimers left in the prose, and
plausible filler written around missing information.

Before:

> While specific details about the company's founding are not
> extensively documented in readily available sources, it appears
> to have been established sometime in the 1990s.

After:

> The company was founded in 1994, according to its registration
> documents.

When the fact is genuinely unavailable, state the absence in one
sentence or omit the section. Never backfill with "likely".

### 20. Sycophancy

Flattering the reader instead of answering.

Before:

> Great question! You're absolutely right that this is a complex
> topic. That's an excellent point about the economic factors.

After:

> The economic factors you mentioned are relevant here.

## Filler and Hedging

### 21. Filler Phrases

Replace with the direct form:

| Filler | Direct |
|--------|--------|
| in order to achieve this goal | to achieve this |
| due to the fact that | because |
| at this point in time | now |
| in the event that | if |
| has the ability to | can |
| it is important to note that | (delete it) |

### 22. Excessive Hedging

Before:

> It could potentially possibly be argued that the policy might
> have some effect on outcomes.

After:

> The policy may affect outcomes.

Keep honest uncertainty. Tighten it instead of deleting it.

### 23. Generic Positive Conclusions

Before:

> The future looks bright for the company. Exciting times lie
> ahead as they continue their journey toward excellence.

After:

> The company plans to open two more locations next year.

### 24. Authority Tropes

Watch for: "the real question is", "at its core", "in reality",
"what really matters", "fundamentally", "the heart of the
matter".

These promise insight; the next sentence usually restates an
ordinary point.

Before:

> The real question is whether teams can adapt. At its core, what
> really matters is organizational readiness.

After:

> The question is whether teams can adapt. That depends mostly on
> whether the organization is ready to change its habits.

## Structure and Drama

### 25. Diff-Anchored Writing

Documentation that narrates a recent change instead of
describing the current system. Unless the document is a
changelog, release note, or migration guide, write current
state.

Before:

> This function was added to replace the previous approach of
> iterating through all items, which caused O(n²) performance.

After:

> This function uses a hash map for O(1) lookups, avoiding the
> O(n²) cost of naive iteration.

### 26. Staccato Drama

One short sentence for emphasis is fine. A chain of clipped
fragments feels engineered.

Before:

> Then AlphaEvolve arrived. It had no preference for symmetry. No
> aesthetic prior. No nostalgia for human taste. The old rules
> were gone.

After:

> AlphaEvolve changed the search because it did not favor
> symmetry or human-looking designs. That made some older
> assumptions less useful.

### 27. Aphorism Formulas

Watch for: "X is the Y of Z", "X becomes a trap", "the language
of", "the currency of", "the architecture of".

Ordinary claims dressed as quotable profundity.

Before:

> Symmetry is the language of trust. Efficiency becomes a trap
> when teams forget the human layer.

After:

> Symmetric layouts often feel more predictable to users. Teams
> can over-optimize workflows and miss how people actually use
> them.

### 28. Theatrical Openers

Watch for standalone: "Honestly?", "Look,", "Here's the thing",
"Let's be honest", "Real talk".

Fine mid-sentence in casual registers. The tell is the
theatrical standalone opener before an ordinary point.

Before:

> Is it worth the price? Honestly? It depends on how often you
> will use it.

After:

> Whether it is worth the price depends on how often you will
> use it.

## Adding Voice Without Inventing It

Removing tells is half the job; sterile prose can still read as
synthetic. Signs of soulless writing: every sentence the same
length and shape, no opinion where one would be natural, no
honest uncertainty, no specific detail, conclusions that sound
like press releases.

Voice repair is allowed only within the invariants:

- Vary rhythm. Short sentences are fine. So are long ones that
  take their time.
- Swap inflated abstractions for the specific claims they gesture
  at.
- Let useful asides and unresolved tension survive.
- First-person reactions, jokes, and editorial attitude may be
  *kept* or *restored* when the author's material contains them.
  They may not be *manufactured*. If the source is a neutral
  report, the humanized version is a cleaner neutral report, not
  a personal essay. When the piece would clearly benefit from a
  stance the author never stated, flag it as a suggestion instead
  of writing it in.

Never inject personality into reference docs, legal text,
security reports, specs, or compliance material unless the
source already carries that voice.

## Full Worked Example

Before:

> Great question! Here is an essay on this topic. I hope this
> helps!
>
> AI-assisted coding serves as an enduring testament to the
> transformative potential of large language models, marking a
> pivotal moment in the evolution of software development. In
> today's rapidly evolving technological landscape, these
> groundbreaking tools are reshaping how engineers ideate,
> iterate, and deliver, underscoring their vital role in modern
> workflows.
>
> Industry observers have noted that adoption has accelerated
> from hobbyist experiments to enterprise-wide rollouts, from
> solo developers to cross-functional teams. Additionally, the
> ability to generate documentation, tests, and refactors
> showcases the intricate interplay between automation and human
> judgment.
>
> * Speed: Code generation is significantly faster, reducing friction and empowering developers.
> * Quality: Output quality has been enhanced through improved training.
> * Adoption: Usage continues to grow, reflecting broader industry trends.
>
> While specific details are limited based on available
> information, it could potentially be argued that these tools
> might have some positive effect. In order to fully realize this
> potential, teams must align with best practices.
>
> In conclusion, the future looks bright. Exciting times lie
> ahead as we continue this journey toward excellence. Let me
> know if you would like me to expand on any section!

First draft of the rewrite:

> AI coding assistants can make some programming tasks faster.
> They are useful for boilerplate, test scaffolding, config
> files, and repetitive refactors. They are less useful when the
> problem requires product judgment, architecture, or domain
> context.
>
> The main risk is not that the tool is always wrong. The risk is
> that it can sound right while being wrong. A suggestion can
> compile, pass lint, and still miss the point.
>
> Productivity metrics are slippery. Accepted suggestions are not
> the same thing as correct suggestions, and correct suggestions
> are not always valuable.

Then ask: what still sounds generated? The paragraphs land with
tidy essay symmetry, and "the main risk is not X, the risk is Y"
is a negative parallelism. One more pass:

> AI coding assistants make the boring parts faster: boilerplate,
> test scaffolding, config files, repetitive refactors. They do
> not help much with architecture or product judgment, and they
> are good at sounding right while being wrong. A suggestion can
> compile, pass lint, and still miss the point.
>
> Treat the assistant as autocomplete for chores, not a
> substitute for review. The output still needs tests and a human
> who knows what the code is supposed to do.
>
> Be careful with the productivity numbers. Accepted suggestions
> are not correctness, and correctness is not always value.

Note what the second pass did not do: it did not add first-person
anecdotes, invented statistics, or opinions absent from the
source. It tightened claims the source already made.

## Attribution

The pattern taxonomy draws on Wikipedia's "Signs of AI writing"
page, maintained by WikiProject AI Cleanup. Use external
references as guidance rather than text to paste, and follow the
source's license when quoting substantially.

The underlying principle: LLM writing gravitates toward broadly
applicable, statistically likely phrasing. The cure is not "make
it messy". The cure is to make it specific, intentional, and
appropriate for the reader.
