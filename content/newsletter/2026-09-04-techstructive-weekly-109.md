---
title: "Techstructive Weekly #109"
date: 2026-09-04T18:45:14Z
slug: techstructive-weekly-109
draft: false
type: post
description: "Not much reading, but watching and learning a lot of things, working with AI to build and refactor systems, among the other things did in the week from 30th August to 5th September 2026"
tags: ["substack"]
---


<p></p>
<h2>Week #109</h2>
<p>A bit of back to work week. Enjoyed the process maybe. Not sure, still figuring out the agentic engineering environment. I am actually lucky to have peers and colleagues that respect each other and understand the current trend of agentic work. “You don’t have to know everything, but you have to own it” that’s a bit hard to digest pill, If you know what I mean. Generating code, but not knowing the details but knowing enough of it to steer the issues. I felt good, not comfortable but not stressed or disrespected the least. I understand not everyone might have the privilege to be in this ai-change-friendly environment.</p>
<p>While being excited for the week, I found a couple of use cases for playing with LLMs.</p>
<ul>
<li><p>There is no openrouter like app on android or ios right?</p></li>
<li><p>There is a scope for a python package or api that gives you free llm routing (like any free llm interface)</p></li>
<li>
<p>Something to make with the needle2 14MB model which is a function tool calling for android</p>
<p> </p>
</li>
</ul>
<p>Would be looking forward in the weekend to do some of these.</p>
<p></p>
<h3>Quote of the Week</h3>
<blockquote>
<p>“Presents are made for the pleasure of the one who gives them, not for the merit of those who receive them”</p>
<p><a href="https://www.goodreads.com/work/quotes/3209783https://www.goodreads.com/quotes/133882-presents-are-made-for-the-pleasure-of-who-gives-them">— Carlos Ruiz Zafón, The Shadow of the Wind</a></p>
</blockquote>
<p>I am reading the shadow of the wind, and father of the main character gives his son a present, the son is not happy. But the father is. Presents truly are for the giver and the by product is the receiver is happy too. What a beautiful little conversation it was, but it shows much about life.</p>
<p></p>
<h2>Read</h2>
<ol><li>
<p><a href="https://martiansoftware.com/articles/ai-written-code-is-still-yours">AI written code is still your code, are you ok with that?</a></p>
<ol>
<li><p>Maybe sometimes not always right? If we generate the code, its the responsibility of the developer to own it that’s the rule.</p></li>
<li><p>I find it a bit frustrating, quite less from last year. However some frustration is still there, since I haven’t found the reason for why each word in the code exists. To find that it takes a lot of effort, if its python, it can be easy with a REPL but other languages, not gonna take easy route, even with python in large code-bases and especially with LLM integration, it becomes harder and harder to make changes reliably and tested.</p></li>
</ol>
</li></ol>
<p></p>
<p>Didn’t read much, was a bit restless and read some book and worked the rest half of the week.</p>
<p></p>
<p></p>
<h2>Watched</h2>
<ul><li>
<p><a href="https://youtu.be/a3C1DMswClQ">Sriniously: HTTP for backend engineers</a></p>
<ul>
<li><p>It was a great explanation of the protocol and the actual details</p></li>
<li><p>Unlike the theoretical stuff, the speaker actually showed each concept with a live demo. The headers part was really great, the methods and request response schema was touched perfectly.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"a3C1DMswClQ","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/HqXr5uXOb5o">Dennis Hirsch Coding LS in C</a></p>
<ul>
<li><p>This was cool to see the file path. I didn’t we can get that low level with so little code in C. I want to do this in Go, but should I.</p></li>
<li><p>The feeling for doing it despite of AI is cool, but is it worth it? Should I be doing it? I am so overwhelmed with so many stuff.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"HqXr5uXOb5o","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<ul><li>
<p><a href="https://youtu.be/F7fe9pa8OeE">Omarchy Quattro</a></p>
<ul>
<li><p>Wow! This is so cool. It feels so lowkey, anybody who had used or customised his linux distro knows its a bit mechanical to do all of it. And omarchy specifically this release is the agentic way of doing that, not a big deal. But the vision to make someone onboard into linux is the big deal.</p></li>
<li><p>The effort to add functionality and the customisability is flawless and effortless.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"F7fe9pa8OeE","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<ul><li>
<p>PlanetScale: Postgres and MySQL Indexes are completely different</p>
<ul>
<li><p>This is cool to know. </p></li>
<li><p>Postgres creates a B Tree  for each index with the actual key and the tuple id , which makes look up of other data very trivial</p></li>
<li><p>MySQL  creates a B Tree with the actual key and the primary key to look up the data separately in the other B-Tree. I might be a little less efficient compared to Postgres, but I am not sure what the trade-off here is, a good question to answer for the weekend. I love databases.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"T3P9TLi5R08","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<p></p>
<h2>Learnt</h2>
<ul><li>
<p>Was reading the Learning SQL book and found some bits of information elsewhere to learn</p>
<ul>
<li><p>SQLite has 2k max columns limit per table</p></li>
<li><p>PostgreSQL has 1600 max columns</p></li>
<li><p>MySQL has 1017 columns</p></li>
<li><p>DuckDB has unlimited * (maybe a bit of very high number)</p></li>
<li>
<p>Oracle has 1000 columns </p>
<p></p>
</li>
</ul>
</li></ul>
<p></p>
<p></p>
<h2>Tech News</h2>
<ul>
<li><p><a href="https://blog.google/innovation-and-ai/models-and-research/gemini-models/3-8-flash-and-3-8-flash-cyber/">OpenAI releases GPT 6 Astra</a></p></li>
<li><p><a href="https://blog.google/innovation-and-ai/models-and-research/gemini-models/3-8-flash-and-3-8-flash-cyber/">Google drops Gemini 3.8 Flash</a></p></li>
<li><p><a href="https://www.anthropic.com/claude-fable-and-mythos-5-1">Anthropic drops Fable and Mythos 5.1</a></p></li>
</ul>
<p></p>
<div><hr></div>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/808">Hackernewsletter</a> (#808 edition), and for software development/coding articles, join daily.dev.</p>
<p></p>
<p>That’s it from the 109th Edition of techstructive weekly. I hope you found it helpful, and relaxing. If not please drop any suggestions, feedback or discussion about certain things you want to in the comments or drop me a message on my <a href="https://www.meetgor.com/contact">socials</a>. </p>
<p></p>
<p>Thank you for reading,</p>
<p>Until next week.</p>
<p>Happy Coding :)</p>
