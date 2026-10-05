---
title: "Techstructive Weekly #113"
date: 2026-10-02T18:45:33Z
slug: techstructive-weekly-113
draft: false
type: post
description: "Reading AI slop for learning systems, watching some useful tricks, among the other things observed and learnt in the week from 27th September to 3rd October 2026"
tags: ["substack"]
---


<p></p>
<h2>Week #113</h2>
<p>A week with a lot of maturing on as a developer. Learnt a lot of stuff, pondered over things.</p>
<ul>
<li><p>We are reading less human content, because most of the times we are talking to the agents. If you are not, then you are not learning enough? I think over the last few months, my reading articles have dropped considerably because there is no time of that leisure to read an 1k&gt; word technical article. I have changed habbits of reading ebooks but that was even there in the last year, this year the agents have taken over the mental cognition of reading the slop, but in a good way maybe.</p></li>
<li><p>With LLMs, we now have the ability to understand codebases and flows of a system at our fingertips. Now its not about is the model smart enough, now its about can you frame the right question out, can you get from it what you want. And yes prompting plays a role, but its not 100% prompting, it also about thinking with first principles and thinking out loud to make the intent clear. Communication, written or verbal has never been more important.</p></li>
</ul>
<p></p>
<p>These are some of my observations over the week, while I was solving some issues in codebases where I had no clue about. Having a LLM can help, but I struggled to ask the questions and didn’t do quite well. But now with that learning, I am approaching it more open-mindedly and framing the right questions.</p>
<p></p>
<h3>Quote of the Week</h3>
<blockquote>
<p>“The search for truth is more precious than its possession.”</p>
<p>― <strong><a href="https://www.goodreads.com/quotes/45649-the-search-for-truth-is-more-precious-than-its-possession">Gotthold Ephraim Lessing</a></strong></p>
</blockquote>
<p>This is true. Now with LLMs, we don’t have the friction of getting an answer, to satisy our curiosity is the possession. The process was the most important part, which now is no more a step. Its like going on a picnic but skipping the travel. What? Its hard to digest but that’s what is the call for adoption or survival maybe? Who know, maybe darwin was true or he will be proven wrong in due course of the next 5 years right? AGI are you listening?</p>
<p></p>
<h2>Read</h2>
<ol>
<li>
<p><a href="https://calnewport.com/its-time-to-investigate-the-ai-labs/">It’s time to investigate the frontier labs</a></p>
<ol>
<li><p>What a time we live in. The people who say, we need to be safe against AI are the only ones creating unsafe technologies. Like, once the incidents have happened, they are still researching on and conducting such experiments. Like what? Does it even make sense, the thing that they think is going to make “humanity vanish” are the one still pursuing the goal of it.</p></li>
<li><p>Really, we need to investigate them. What is the purpose of their experiments if it can go and attack other systems?</p></li>
</ol>
</li>
<li>
<p><a href="https://bookofrevenue.com/blog/6ab81e9a97a13f0001f7e4e1/postgres-at-time-zone-u-does-not-do-what-you-think-it-does">Postgres timestamp vs timestamptz</a></p>
<ol><li><p>Oh timezones and developers. Love the ever love-hate relationship. There is no love actually. Dealing with time zones as a developer especially in databases, is a classic brain fader.</p></li></ol>
</li>
<li>
<p><a href="https://simonw.substack.com/p/2026-in-llms-so-far">2026 in LLM so far</a></p>
<ol>
<li><p>WOW, this is so nostalgic now. It seems like eternity, the claw moment. The overnight rise and fall of the moltbook and what not.</p></li>
<li>
<p>The progression of the pelican riding bicycle has become one of the classic test for any LLM.</p>
<p></p>
</li>
</ol>
</li>
</ol>
<p></p>
<p></p>
<p></p>
<h2>Watched</h2>
<ul><li>
<p><a href="https://youtu.be/4TE1xErXwGc">Scaling 7M Tables in Postgres, talk from Kailash Nadh CTO at Zerodha</a></p>
<ul><li><p>This is cool. I didn’t knew postgres can handle that much tables. But it makes sense, why would a battle-tested database hold off in the number of tables we can create.</p></li></ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"4TE1xErXwGc","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<p></p>
<ul><li>
<p><a href="https://youtu.be/Ag3nn9BgWP0">Casey Muratori rating the Digging Vibe Games: The Standup </a></p>
<ul><li><p>This was entertaining. We knew how far the agents can go. THe taste is the moat now for developers.</p></li></ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"Ag3nn9BgWP0","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<p></p>
<h2>Learnt</h2>
<ul>
<li>
<p>We can hack DNS to query text only records: check out <a href="https://www.dns.toys/">dns.toys</a></p>
<ul><li><p>I watched a talk from Kailash Nadh the CTO of Zerodha who had created this and its so intruiging to make such tools</p></li></ul>
</li>
<li>
<p>We can make Postgres store over million tables</p>
<ul><li><p>That is a clever hack</p></li></ul>
<p></p>
</li>
</ul>
<p></p>
<p></p>
<h2>Tech News</h2>
<ul>
<li>
<p><a href="https://www.worldlabs.ai/blog/amd-announcement">World Labs is joining AMD</a></p>
<ul><li><p>This is good. At least it's not Nvidia.</p></li></ul>
</li>
<li><p><a href="https://openai.com/index/introducing-gpt-6-1-sol/">OpenAI drops GPT 6.1 Sol</a></p></li>
<li><p><a href="https://openai.com/index/introducing-dots/">OpenAI releases dots, always-on agents</a></p></li>
<li><p><a href="https://blog.google/innovation-and-ai/models-and-research/gemini-models/gemini-4-argon/">Google drops Gemini 4 Argon: who cares?</a></p></li>
<li>
<p><a href="https://www.anthropic.com/claude-sonnet-5-5">Anthropic releases Sonnet 5.5</a></p>
<p></p>
</li>
</ul>
<p></p>
<div><hr></div>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/812">Hackernewsletter</a> (#812 edition), and for software development/coding articles, join daily.dev.</p>
<p></p>
<p>That’s it from the 113th Edition of techstructive weekly. I hope you found it helpful, and relaxing. If not please drop any suggestions, feedback or discussion about certain things you want to in the comments or drop me a message on my <a href="https://www.meetgor.com/contact">socials</a>. </p>
<p></p>
<p>Thank you for reading,</p>
<p>Until next week.</p>
<p>Happy Coding :)</p>
