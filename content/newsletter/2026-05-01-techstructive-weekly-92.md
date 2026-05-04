---
title: "Techstructive Weekly #92"
date: 2026-05-01T18:45:22Z
slug: techstructive-weekly-92
draft: false
type: post
description: "Back to work week, reading ai discourse, watching the end of subsidized token era, among the other things read, watched and learnt in the week from 26th April to 2nd May 2026"
tags: ["substack"]
---


<h2>Week #92</h2>
<p>I enjoyed working this week, I am surprised I am saying that. The same thing a month back was dreading me, and now its making me wake up each morning. AI is kind of a progression and also a setback in some sense. I don’t know, but the usage will steer where the world moves with it. There will be a divide, of course, and if I am correct, it will create the biggest divide the world has ever seen.</p>
<p>But till that happens, just sip in the beauty of the nature, and read some handwritten articles, human voice, and some hot takes on the internet about tech.</p>
<h3>Quote of the week</h3>
<blockquote>
<p>“To know that we know what we know, and that we do not know what we do not know, that is true knowledge.”</p>
<p>— Confucius</p>
<p>― <a href="https://www.goodreads.com/quotes/383876-confucius-said-to-know-that-we-know-what-we-know">Henry David Thoreau, Walden or, Life in the Woods</a></p>
</blockquote>
<p>How true is that? If you know something, you are eager to express or impress, but when you don’t, you shy away or fear away. The later part is what makes us humans. The latter, the complete part, is what distinguishes us from the AI bots. They will hallucinate and pretend the opposite when they don’t know something, they are unpredictable even in the first known part.</p>
<p>Knowledge is not learning something, its about intuition and learning to adapt. Being honest is the best thing to do to yourself.</p>
<p></p>
<div><hr></div>
<h2>Read</h2>
<ol>
<li>
<p><a href="http://ozark.hendrix.edu/~yorgey/forest/00FD/index.xml">To my students</a></p>
<ul><li><p>Read this article, if you are a student or not. You will find answers and sentences that matter to your soul and not your ego.</p></li></ul>
<blockquote><p>Above all, be motivated by love instead of fear.</p></blockquote>
<ul><li><p>Everyone is in fear, unsure of the future, don’t be. Recollect the feeling when you started to do the thing in the first place, I bet it was not money, if it was, and still is, then I don’t know, if it was love, keep on steering.</p></li></ul>
<blockquote><p>Have the courage to go slowly, especially when everyone else is telling you that you need to go fast and cut corners</p></blockquote>
<ul><li><p>Vibe coding, running agents, and shipping without reviewing, these are hype and will fade with time. In professional development, at least, this won’t be the standard as far as I know.</p></li></ul>
</li>
<li>
<p><a href="https://mitchellh.com/writing/ghostty-leaving-github">Ghostty is leaving GitHub</a></p>
<ul>
<li><p>This is bad. Oh my god, he has spent almost two decades on GitHub. I might not have lived consciously for that many years.</p></li>
<li><p>Man, this person has breathed and lived through GitHub. I can resonate with the line “When I went through tough breakups, I lost myself in open source... on GitHub”. I have done that to some extent, too. Not at his scale, but when I didn’t get any internships, got rejected in interviews, I found SQLC, Turbot, and a handful of organisations to fill my hunger for learning and programming. GitHub was the center.</p></li>
<li><p>I don’t rely on GitHub since I don’t maintain or do any hustles anymore, just because life is tough and so is the AI hype. I can relate to the outages happening these days. One of my projects <a href="https://github.com/Mr-Destructive/flight-observatory">flight-observatory</a>, had a few failing actions notifications. I thought it might be 403 or some upstream external API issues,  but GitHub was the reason behind it.</p></li>
</ul>
</li>
<li>
<p><a href="https://lemire.me/blog/2026/04/27/you-can-beat-the-binary-search/">You can beat binary search</a>: Quad Search</p>
<ul>
<li><p>This is creative. Instead of just halving the list, we do a buckets of 16, quads, and compute which of the buckets the element can be. Since we have the max of each bucket. We can find the bucket and parallelly then compare each element in it.</p></li>
<li><p>Surprisingly, this is better for large arrays, and with multiple cores this will outperform binary search easily. Because we are not diving the array into halves  anymore, we are picking the most possible region where the element could be and parallel searching it in those.</p></li>
</ul>
</li>
<li>
<p><a href="https://peterdohertys.website/blog-posts/full-text-search-w-duckdb.html">Full Text search with DuckDB</a></p>
<ul><li><p>This is quite similar to what we do in SQLite. Pretty simple. I like how DuckDB just keeps on making transitioning to modern databases, but still keeping the backwards compatibility and keeping it simple and dead easy to use just like SQLite.</p></li></ul>
</li>
<li>
<p><a href="https://gregraiz.com/blog/local-vibe/">I got sick of remembering port numbers</a></p>
<ul>
<li><p>He got sick of remembering port numbers, so he just reinvented ngrok for local?</p></li>
<li><p>Yes, it is different than assigning names to services, but I feel that it’s doing that similarly. It’s a clever trick.</p></li>
<li><p>This is the phase where people are building anything they want; the idea of personal software and developer tooling is just skyrocketing.</p></li>
</ul>
</li>
<li>
<p><a href="https://www.pootlepress.com/2026/04/ai-tokens-and-the-gathering-storm/">AI, Tokens and the gathering storm</a></p>
<ul>
<li><p>This is quite evident from the recent events and rug pull of Anthropics subsidized tokens, and GitHub’s move from tier based to usage based billing. Its on the wall, its the end of the subsidized token era.</p></li>
<li><p>The next few months might actually tell us how expensive are these tokens.</p></li>
<li><p>Local LLMs be ready for shinny demo, the dark days might be about to end.</p></li>
</ul>
</li>
<li>
<p><a href="https://viggy28.dev/article/local-llm-seven-wrong-answers/">I asked Local LLM to add 23 numbers, I got 7 different wrong answers.</a></p>
<ul>
<li><p>My last sentence was jinxed. LLMs can’t do math yet. Sigh!</p></li>
<li><p>I did explore this last year and they failed miserably, till date this hasn’t improved. I think we can create a skill for this shall we, I don’t want it to use a python interpretter for such a trivial thing to add. Can we make them add? It would consume a bit more tokens but let’s see. I have some ideas here.</p></li>
<li><p>But the point being, Local LLMs are not quite the hype that the propreitary models live up to.</p></li>
</ul>
</li>
<li>
<p><a href="https://www.0xsid.com/blog/agentic-coding-fatigue">Agentic Coding Fatigue</a></p>
<ul>
<li><p>True. This is me reading to myself. Writing code gives us the clarity, becuase we were in the weeds, we know why each of the if else was written, we knew the edge cases, but now? We just run the slot machine to “fix it” prompt and cross our fingers so that it works.</p></li>
<li><p>The balance of understanding and the need to understanding, is wide and is growing faster than ever. I don’t know which one to lean onto. The former sounds like there won’t be any harm. But the software is such a term a thing, that nobody really cares how well you understand it. Only the output matters, not the usage of design patterns or Python or Golang.</p></li>
</ul>
</li>
<li>
<p><a href="https://ky.fyi/posts/ai-burnout">Do I belong in Tech Anymore?</a></p>
<ul>
<li><p>This one is really rough to read through. Not that it’s bad writing or thoughts, its is dead real and truthful. I myself am not able to express the true things that I have to go through, maybe I am not as privileged or have made my situation like this.</p></li>
<li>
<p>The principles that the author listed are gold, make sure we point it here too</p>
<ul>
<li><p>Things that are worth doing are worth doing well.</p></li>
<li><p>Things that are done well require time and effort.</p></li>
<li><p>You make meaning through the doing.</p></li>
<li><p>Ideas are common; effort is not.</p></li>
<li><p>There are no shortcuts.</p></li>
</ul>
</li>
<li><p>Frame it in gold, or memorise them, it is better to remember this than to be a slave to an agent or a corporation that is.</p></li>
</ul>
</li>
<li>
<p><a href="https://migrainebrain.bearblog.dev/people-who-dont-use-ai-will-be-left-behind/">People who don’t use AI will get left behind</a></p>
<ul>
<li><p>I will frame it as “People who use only AI will definitely get left behind”</p></li>
<li><p>People would stop learning and thinking if they just learnt to use AI and not anything else around it.</p></li>
</ul>
</li>
<li>
<p><a href="https://idiallo.com/blog/have-you-seen-the-new-xl-ai-parody">Have you seen the new Excel?</a></p>
<ul>
<li><p>Learn Excel or get left behind, learn AI or get left behind. Everyone is on a roll with hypes. I do agree, excel has just revolutionized the way people think of software or even any form of work.</p></li>
<li><p>It has abstracted the code in such a way that people rely on the results and not the code. Superb product.</p></li>
<li><p>I worry AI is doing the same, but for text generation, for code generation. It might be the de facto thing to produce custom software on the fly.</p></li>
<li><p>Until that happens, keep learning Excel ;)</p></li>
</ul>
</li>
<li>
<p><a href="https://openai.com/index/where-the-goblins-came-from/">Where the goblins came from</a></p>
<ul>
<li><p>This was an interesting nerdy read. I liked they are open about it. Though they don’t back away from shoe-horning codex glazing. I don’t doubt they use codex, but its quite a bit of token-hungry and slow model.</p></li>
<li><p>The problem of the goblin is really interesting. The goblin behavior wasn’t a bug, according to them, it came from reinforcement learning rewards that favored playful, creature-based metaphors in the “Nerdy” personality.</p></li>
<li><p>Those rewarded patterns spread and amplified across training loops, even outside the original context. It resulted in a small stylistic quirk becoming a generalized model habit, showing how local incentives can unintentionally shape global behavior.</p></li>
</ul>
</li>
</ol>
<p></p>
<p></p>
<h2>Watched</h2>
<ul><li>
<p><a href="https://youtu.be/HgNKa9UlRF8">Making Sense of the AI hype</a> - Wading through AI -Episode 2</p>
<ul>
<li><p>The topic of engineer writing the blog and the manager or the sales person making the video is a reality and the only needle making the AI labs floating.</p></li>
<li><p>Its all gimmick play. Models aren’t quite capable of creating a novel thing like a compiler from scratch, that is, way far currently.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"HgNKa9UlRF8","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<ul><li>
<p><a href="https://youtu.be/NZa5lApeFic">AI isn’t taking jobs its worse</a></p>
<ul><li><p>Anthropic is eating the money, by creating fear. Like organisations need things quick so they burn tokens, for that they need to cut down on cost, which surprise is got by firing employees. So, the money eventually flows to Anthropic by cutting the job, which is hype honestly speaking. Not sustainable, they themselves have proven it.</p></li></ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"NZa5lApeFic","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<h2>Learnt</h2>
<ul><li>
<p>String Pool in Java</p>
<ul>
<li><p>String Pool is a special memory area inside the heap where unique string literals are stored. The idea is instead of creating a new object every time, it reuses existing strings from this pool when possible.</p></li>
<li><p>A string in a string pool is a deduplicated storage of string literals inside the heap.</p></li>
</ul>
</li></ul>
<ul><li>
<p>Cursor agent &gt; Claude code</p>
<ul>
<li><p>I like cursor agent cli, it is intuitive and doesn’t block the user, Claude code is like I(the agent) am the owner. I don’t like that.</p></li>
<li><p>Cursor has the ability to view the output of the tool calls with Ctrl+o and let the screen continue, whereas the Claude just blocks the view.</p></li>
<li><p>I hate anthropic that is one reason. Maybe, but honestly speaking, claude code is blaoted.</p></li>
</ul>
</li></ul>
<p></p>
<h3>Random Tidbits</h3>
<ul><li>
<p>https://www.numberempire.com/</p>
<ul><li><p>A cool webpage to see different patterns for a given number</p></li></ul>
</li></ul>
<p></p>
<h2>Tech News</h2>
<ul>
<li>
<p><a href="https://github.blog/news-insights/company-news/github-copilot-is-moving-to-usage-based-billing/">Copilot is moving to usage-based billing: The end of a subsidized token era</a></p>
<ul><li><p>This is a sign guys, if you still think, AI is not hype. The hype was at its peak, and now its wading. The subsidized token era is over. The real cost is on the labs and the providers that are now realising the importance of sustenance.</p></li></ul>
</li>
<li>
<p><a href="https://www.warp.dev/blog/warp-is-now-open-source">Warp is now open source: Another attempt to win the race?</a></p>
<ul>
<li><p>This looks like a desparate attempt to get into the lead of the race, I am not against it by any chance, its a solid move.</p></li>
<li><p>More organisations should do this now, espcially the ones that created this right? Yes I am talking about A…. Never mind.</p></li>
</ul>
</li>
<li>
<p><a href="https://github.com/anthropics/claude-code/issues/53262">HERMES.md in git commit messages causes requests to route to extra usage billing instead of plan quota</a></p>
<ul>
<li><p>Again. Never mind. Nobody would be using their coding agent.</p></li>
<li>
<p>They don’t leave a chance to not to hate them.</p>
<p></p>
</li>
</ul>
</li>
</ul>
<div><hr></div>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/792">Hackernewsletter</a> (#792nd edition), and for software development/coding articles, join daily.dev.</p>
<div><hr></div>
<p>That’s it from this week. It was a bit of a work-heavy routine. Finally feeling back, maybe not peacefully good but warmly good. Looking up to slowing down with manual coding and building some impactful things in the era of slop.</p>
<p>Till then,</p>
<p>Happy Coding :)</p>
