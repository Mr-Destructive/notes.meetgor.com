---
title: "Techstructive Weekly #89"
date: 2026-04-10T18:45:43Z
slug: techstructive-weekly-89
draft: false
type: post
description: "Learning design patterns, reading about AI and people, building the intuition for chess board, among the other things read, created, watched, and learnt in the week from 5th to 11th April 2026"
tags: ["substack"]
---


<h2>Week #89</h2>
<p>I did a few things this week. I made a little app, wrote a half-baked draft for part 2 of the flight observatory, and learned design principles. </p>
<p>But something feels off.</p>
<p>I am feeling a bit rusty with programming lately, not sure where the joy in coding is lost. I hate managing agents, I hate reviewing 2k lines of code in 2 hours. It’s just not what I had signed up for. I am complaining, yes, but it’s just a weird thing.</p>
<p>Not feeling great to be honest, have ideas, but agentic-engineering seems to have completely overtaken. I am saying that I miss writing to code by hand, or I want to, but it feels like the layer of agents has replaced that muscle memory of thinking while coding. </p>
<p>Anyways, I still find the joy of building, not sure how long it will last. At least can be hopeful about it.</p>
<h3>Quote of the week</h3>
<blockquote>
<p>“We are often most in the dark when we are the most certain, and the most enlightened when we are the most confused.”</p>
<p>— <a href="https://www.goodreads.com/quotes/6503484-we-are-often-most-in-the-dark-when-we-are">M. Scott Peck</a></p>
</blockquote>
<p>Maybe, not sure. I feel I am up to something. Maybe the path is clearing up; the night is darkest just before the dawn. Keeping a positive attitude is the way forward. I think I have a lot of things to juggle and grasp, which is why I feel lost. Once I can manage the right things in the right direction at the right time, I can feel a bit confident. Cannot see the other way round. Escape is not an option, and I don’t like it either. There is no fun and learning in choosing the easier route. Ending with another quote, easy routes, hard life, hard routes, easy life.</p>
<p></p>
<div><hr></div>
<h2>Created</h2>
<ul>
<li>
<p><a href="https://chessight.meetgor.com/">Chessight</a> → https://chessight.meetgor.com/</p>
<ul>
<li><p>I built it as my experiment to learn agentic-driven development, of course, why not! I used Codex, and it just one-shotted the entire UI, I deliberately asked it to use Vue, and it did it phenomenally well.</p></li>
<li><p>Ok, what is the project? It's just a visual training application that can help you understand the chessboard and its algebraic notation.</p></li>
<li><p>The idea is simple, not to teach chess, but to train how you see the board, like g5 shouldn’t be something you calculate, it should just click when looking at the board.</p></li>
<li><p>I was watching FIDE Candidates 2026 (Prag vs Sindarov), and I was intrigued by Sagar Shah and all the other chess players and commentators who can just name the moves. I was still figuring out the move, and they named the exact step. I wanted to build that muscle memory. I looked for chess.com and Duolingo, but neither of them provides this. They don’t start from the fundamentals.</p></li>
<li><p>After playing it for a while, it feels like one of those small skills that unlock a completely different level of thinking in the game. I am not a chess expert, I don’t play chess regularly, but I like it, just for fun and flexing the -100 IQ brain of mine.</p></li>
</ul>
</li>
<li><p>Draft for Flight Observatory technical report and devlog.</p></li>
</ul>
<p></p>
<h2>Read</h2>
<ol>
<li>
<p><a href="https://medium.com/@lukas_kosinski/ai-and-remote-work-is-a-disaster-for-junior-software-engineers-a377b1d8ed20">AI and Remote work is a disaster for software engineers</a></p>
<ul>
<li>
<p>Spicy take. I’ve been working remotely for ~2 years as a junior-ish backend engineer, and I use AI daily. If AI + remote sabotages juniors were broadly true, I should be a pretty weak engineer by now. That hasn’t been my experience at all.</p>
<blockquote><p>remote work sabotages careers of youngsters</p></blockquote>
</li>
<li>
<p>I think remote doesn’t sabotage. It removes forced structure. If you don’t build your own, you stagnate and just feel lost. There’s no overhearing seniors, no accidental learning, no pressure to “look busy”. If you just do assigned tickets and log off, yeah, sure, you’ll stagnate. But that’s not a remote problem; that’s an ownership problem. You can sit in an office and do the same thing.</p>
<blockquote><p>AI — too little brain stimulation</p></blockquote>
</li>
<li><p>Same with AI. If you use it to skip thinking, you’ll get faster at producing things you don’t understand. But if you use it to explore, ask “why does this break?”, generate edge cases, compare approaches, it’s like having a patient senior who will walk through things with you endlessly. The difference is in how you engage with it, the intention rather than the environment.</p></li>
<li><p>I do agree that juniors need tight feedback loops and exposure to better engineers. That’s harder to get remotely, and most companies don’t compensate for it well. But the answer isn’t “go back to office”, it’s “design better learning environments”, more deliberate mentorship, better code reviews, more context sharing.</p></li>
<li><p>Also, the market point is real: AI is eating the bottom layer of trivial work. But that just raises the bar; it doesn’t remove the path. Juniors now need to show they can reason about systems, not just implement tickets.</p></li>
<li><p>The point is, if an individual is naturally curious, remote or on-site doesn’t matter, he’ll succeed in whichever environment.</p></li>
</ul>
</li>
<li>
<p><a href="https://shkspr.mobi/blog/2026/04/someone-at-browserstack-is-leaking-users-email-address/">Someone at BrowserStack is leaking mail addresses of customers</a></p>
<ul>
<li><p>This is all happening. The security of software is crazy right now. We just had a week of supply chain attacks on Python and JavaScript ecosystems. And now its getting into the weeds of the application. Not too far from AGI, right?</p></li>
<li><p>I don’t reckon that the company might be leaking it intentionally, it might be some vendor they aren’t fully aware of, or have really vibe-shipped something. Not sure.</p></li>
<li><p>Excited for the next post.</p></li>
</ul>
</li>
<li>
<p><a href="https://ben.page/microservices">Does coding with LLMs mean more microservices</a></p>
<ul>
<li><p>Yep, this is a valid and good observation. I have mostly written my side projects in cloud functions (Netlify, Vercel, or Cloudflare). I don’t mind the complexity of managing a bunch of environment variables and setup. It’s a one-time thing, and I know looking for them is a bit of chaos, but I have AI-agents to clean up as and when necessary.</p></li>
<li><p>Maintaining a monolith is no joke with LLMs, it used to change hundreds of out-of-the-blue things, nonsensically, but it has reduced almost to 0. But the cost of changing one thing is still high in those environments,  since a change without review can cause a catastrophic production failure. (Just like GitHub is having right now) </p></li>
</ul>
</li>
<li>
<p><a href="https://www.jmduke.com/posts/software-never-had-a-soul.html">Software never had a soul</a></p>
<ul>
<li>
<p>Beautifully put.</p>
<blockquote><p>“You do not even need better or faster tools. You just need to really mean it.”</p></blockquote>
</li>
<li><p>You need to mean it, you need to care it, this is the crux of being a developer. We should not be caring about the code but caring for the user’s problem, the product. </p></li>
</ul>
</li>
<li>
<p><a href="https://www.anildash.com/2026/04/06/people-love-to-work-hard/">People love to work hard</a></p>
<ul>
<li><p>Great point. This is about not what you do, but how people empathize with your work. If they don’t they say “you don’t want to work”, if they do you can see the title of this post and read it.</p></li>
<li><p>A very human post from this author. Kind of heart-touching. There is nothing more to say, like it’s something you feel and can’t describe in words. </p></li>
</ul>
</li>
<li>
<p><a href="https://ultrathink.art/blog/sqlite-in-production-lessons">Lessons from using SQLite in production</a></p>
<ul>
<li><p>Useful bit here. The WAL mode is absolutely a clutch, most people don’t know of, and just lament about it being a single-writer constrained DB. It’s a super-powerful and very versatile lightweight database.</p></li>
<li><p>Folks at Turso are making it even better with read-write replicas and syncing, having a daemon for SQLite. The future might be apps full of SQLite.</p></li>
</ul>
</li>
</ol>
<p></p>
<h2>Watched</h2>
<ul><li>
<p><a href="https://youtu.be/wc8FBhQtdsA">Simon Wilison on Lenny’s Podcast</a></p>
<ul>
<li><p>This was a relieving video. I thought I was the only a minority amongst the developers who got exhausted after a few hours of AI prompting. It is a natural thing. We are doing the creme layer, the ephiony of the task, this won’t be sustainable.</p></li>
<li><p>The other thing is that he had more than 2 decades of experience.I am not giving excuses, but still, that is a lot more than just a 2-3 year me trying to wrestle the concept of agentic engineering.</p></li>
<li><p>It was refreshing to see this, his enthusiasm is truly contagious. I want to build more now. But mindfully. His tools section is really wild. </p></li>
<li><p>I fit in the middle, I not a junior or a fresher anymore, I am not more than 2 years into the industry. Alas! I am in the middle. But still I think, its not doom and gloom, I am not going to be passive and let the AIs take over my brain.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"wc8FBhQtdsA","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/vvhC64hQZMk">WhatsApp System Design By GKCS</a></p>
<ul><li><p>Well explained. I thought it might be more complex or elegant. But it was neither. Its a chat app nonetheless that scales, that’s it.</p></li></ul>
</li></ul>
<p></p>
<div class="youtube-wrap" data-attrs='{"videoId":"vvhC64hQZMk","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/FcbsppIX0bg">Building Zepto LLD, System Design</a></p>
<ul>
<li><p>This was a good design. The system feels a bit chaotic, but was understandable.</p></li>
<li><p>The inventory, delivery, product, and everything could be a system on its own, but in the real world, systems are created by combining multiple systems.</p></li>
<li><p>Basically, we have products, orders, and order items. The inventory manager finds the product in a dark store(large storage or inventory) and then assigns the nearest dark store with those products or multiple stores, and maps the available delivery agents to those shipping those products as orders.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"FcbsppIX0bg","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/fqySz1Me2pI">How UPI Payments work</a></p>
<ul>
<li><p>Oh, so every bank, small or national or not, has to go through a nationalised bank(or one of them), damm!</p></li>
<li><p>I wonder how and why that is so fast? Is it because the banks already have mapped out the route to each bank amongst them, yes, that’s probably the reason for the unified name. </p></li>
</ul>
</li></ul>
<p></p>
<div class="youtube-wrap" data-attrs='{"videoId":"fqySz1Me2pI","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/v9ejT8FO-7I">Strategy</a> and <a href="https://youtu.be/_BpmfnqjgzQ">Observer</a> Design Pattern</p>
<ul>
<li><p>Banger of a series. It was a great explanation, and it felt really enthusiastic.</p></li>
<li><p>Strategy patterns are like polymorphism, but for a specific component or functionality. Like payment can be a strategy pattern, where we can implement UPI, card, net banking, etc.</p></li>
<li><p>The strategy pattern is like implementing a family of algorithms and encapsulating the algorithms in their own, and we can make them interchangeable.</p></li>
<li><p>Observer pattern is like a push strategy rather than the client polling it. It makes it efficient and clean. When the observable changes, it notifies (pushes) the changes to the observers, and then it can fetch the changes since it can decide what it wants.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"_BpmfnqjgzQ","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<h2>Learnt</h2>
<ul>
<li>
<p>System design is a really broad field</p>
<ul>
<li><p>I thought it was a bunch of arrows and squares. But its a little more subtle than that. It’s about abstract design, and schema too, without going into the code details. That’s a bit tricky part, because creating an abstraction not low-level as code and as high as a class or an entity is a bit tricky balance to maintain. You’ll either go all in or barely touch the core problem.</p></li>
<li><p>There are design patterns and principles to use, diagrams for them, UML, ER, HLD, LLD, and whatnot.</p></li>
</ul>
</li>
<li>
<p>Design Patterns are wired, but important</p>
<ul>
<li><p>This came as a bummer in an interview. I knew patterns, but couldn’t recollect their names, and was dead on the lines. I wonder if people actually care about theory in the age of LLMs. It has made learning on the fly so easy that I am starting to not learn anything deeply. That’s a hard realisation to come out of.</p></li>
<li><p>This came to me, and I immediately started learning and reading about, I crunched almost halfway through “<a href="https://theobjectorientedway.com/">The Object Oriented Way</a>” by Christopher Okhravi, and referencing the “Head First Design Patterns”, the mammoth of a book by O’Reilly.</p></li>
<li><p>I am finding it rewarding now to go through them. Really, LLMs might make us the dumbest creature on earth one day. The greatest devolution in the history of Earth.</p></li>
</ul>
</li>
<li>
<p>We are in a software boom, and not doom</p>
<ul>
<li><p>I am realizing that there are more software engineer roles open, more need for software, but not quite one-to-one. It’s like people can make software that they once couldn’t, but now can do a bit easily. But the only thing that keeps them holding back is that they or anyone doesn’t know what the code actually does. And I think as software developers, this edge will be vital more than ever.</p></li>
<li><p>Till people cope with it, it’s time to go deep and explore the depths of software.</p></li>
</ul>
</li>
<li>
<p><a href="https://github.com/arman-bd/guppylm">You can make your own small LLM</a></p>
<ul><li><p>This is very cool. I like it and found it almost at the right time. I wanted to build an LLM with some specific data, but I was not sure what that could be. This guy just created a bunch of fish and food sentences and called it a day. Legendary stuff. </p></li></ul>
<p></p>
</li>
</ul>
<p></p>
<h2>Tech News</h2>
<ul><li>
<p>Anthropic with Claude Mythos, Project Glasswing, eh?</p>
<ul><li><p>Maybe it was a great mode, maybe it was a meh! model. I don’t know since I don’t have access. Period.</p></li></ul>
<p></p>
</li></ul>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/789">Hackernewsletter</a> (#789th edition), and for software development/coding articles, join <a href="http://daily.dev/">daily.dev</a>.</p>
<div><hr></div>
<p>It was a slow week. To be honest, disappointing in some sense. Not much movement or progress. But the next week might just be good or great. That’s the hope as usual.</p>
<p>That’s it for this week.</p>
<p></p>
