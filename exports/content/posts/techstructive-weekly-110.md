---
title: "Techstructive Weekly #110"
date: 2026-09-11T18:45:09Z
slug: techstructive-weekly-110
draft: false
type: post
description: "Created free llm router, annoyed with gemini models, among the other things read, watched and learnt in the week from 6th to 12th September 2026"
tags: ["substack"]
---


<p></p>
<h2>Week #110</h2>
<p>It was a hammering week. 3k lines shipped to prod, no issues. Agentic coding on the move. I am not sure, I should be proud or not, but not bragging about it, it was a win for me. The code was obviously reviewed from me :| and no one else even spent an second of eye on it, shipped, nothing broke. It was a model updation + some tracing and standardization of model routing. There were 2 major things at once, if you are updating models, that could lead to unknowns that you even don’t know, and that’s what happened only 3 instances among 3k requests in the week. The other being tracing and model standardization library addition. It worked in the end, there are improvements to be done, there always are. But a fair learning after recovering from sick 2 week leave. I call it a comeback motivation from nature. Thanking and grateful for it.</p>
<p>On to the week work, it was good work, I now like working with agents, there is a learning involved, if you are curious and like exploring codebases, its a great tool. I mostly use claude and occasional cursor, I am designing apis or building deeper integrations. I would like to explore more things in the weekend, but it gets a bit tiring. So might read some books and chill out.</p>
<p></p>
<h3>Quote of the Week</h3>
<blockquote>
<p>“I thought there were no more ghosts there than those of absence and loss” <br><br>“And for a moment I thought there were no more ghosts there than those of absence and loss, and that the light that smiled on me was borrowed light, only real as long as I could hold it in my eyes, second by second.”</p>
<p>― <a href="https://www.goodreads.com/quotes/6648408-and-for-a-moment-i-thought-there-were-no-more">Carlos Ruiz Zafón </a> in The Shadow of the Wind</p>
</blockquote>
<p>What a beautiful quote that hits timely. There are no ghosts, after a age in childhood, you don’t fear darkness anymore, as grow from child to adult, you start fearing the loneliness and absence of someone. For me there is a fear of losing my mother and father, everyone has, but its more scarier than my own death, because these two are the only one that can truly understand me, without me saying anything. Light is something they have provided me constantly and its absence is what brings me fear.</p>
<p></p>
<h2>Created</h2>
<ul><li>
<p><a href="https://github.com/mr-Destructive/free-llm-router">Free LLM Router</a></p>
<ul>
<li><p>I discussed it last week and here is the repo today. I did manage to get it working after gpt not willing to implement the api, opencode to the rescue and it created it correctly with the freely available llms on the web</p></li>
<li><p> It uses llms like gemini, kilo, zen, ovh and llm7 providers. I know they will get blocked in some time, but its worth keeping a repo of a quick check on how somethings can work, just a free call without having the auth or a paywall.</p></li>
<li><p>Yes of course share private info in these chats, it might be used for training, but these days people are rarely concerned about security and privacy right?</p></li>
<li><p>This week I’ll try to deploy it on cloudflare or vercel cloud functions. Need to first verify if those environments will block the request, my suspicions is that they will.</p></li>
</ul>
</li></ul>
<p></p>
<p></p>
<h2>Read</h2>
<ol><li>
<p><a href="https://shopify.engineering/back-to-native">Native is the future of mobile at Shopify</a></p>
<ol>
<li><p>This is cool, using ai to review the screen to make incremental changes. Helix is a nice idea to port one system to other, in this case React Native to Kotlin or native mobile application port.</p></li>
<li><p>The premise of software development has changed and the technology is also changing, because it does still mater on which technology is used to make it a pleasant experience to users.</p></li>
<li><p>The code can be generated at fly, so generating code is not the issue, maintaining two repos is also not a big deal but porting or migrating is. </p></li>
</ol>
</li></ol>
<p></p>
<p></p>
<p></p>
<h2>Watched</h2>
<ul><li>
<p><a href="https://youtu.be/4B4R2T4w7Kg">AGI is here, and its bad right?</a></p>
<ul>
<li><p>This is really bad. The person working at Anthropic is saying it, that they are working irresponsibly. </p></li>
<li><p>If they know its wrong direction, why even move, you know I think its about the “power” comes with “ego” and ego dissolves the ability to think. Maybe they have become too arrogant to address the elephant in the room.</p></li>
<li><p>They can make AGI but what the price humans will pay for that? Losing the soul? Devastating the whole earth? Who knows, we can be optimistic but for what reason?  </p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"4B4R2T4w7Kg","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<ul><li>
<p><a href="https://youtu.be/5KvY8CnBB3w">You should stop pretending you understand the codebase</a></p>
<ul>
<li><p>You cannot know everything about the codebase. You can know something but not everything. If you think you know everything then you are probably wrong, or you are not shipping enough, or thinking and dreaming big enough.</p></li>
<li><p>Agents make it really easy to make the change but the context to make that is really critical. </p></li>
<li><p>The other thing that caught me was that you cannot re-create codebases from agents unless you have the context and the mental model of the existing one, truly a great piece of thought to keep in mind.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"5KvY8CnBB3w","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/8faAKQeFkgk">You are the project </a></p>
<ul>
<li><p>The project is not the point, the process of transforming yourself into someone you love to be is.</p></li>
<li><p>banger of a video, shows that AI cannot replace the satisfaction and the development of skill that the friction or the process of going through the actual act of doing it.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"8faAKQeFkgk","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<p></p>
<h2>Learnt</h2>
<ul><li>
<p>Gemini models have moved from thinking budget to thinking tiers and its annoying as hell</p>
<ul>
<li><p>This was a shock for me, I had kept 5k for some thing and expected it to work for over a year, and now some model is being deprecated and it sucks to switch from token count to levels like low,medium and high where there is no constraint and adherence to the level.</p></li>
<li><p>This is a bad strategy from Google, I don’t like it, there is no restriction on what low or high means, and don’t even select medium that is the worse of all.</p></li>
</ul>
<p></p>
</li></ul>
<p></p>
<p></p>
<h2>Tech News</h2>
<ul>
<li><p><a href="https://www.deepseek.com/en/news/deepseek-v4-1-flash/">DeepSeek releases V4.1 Flash</a></p></li>
<li><p><a href="https://tailwindcss.com/blog/tailwind-is-joining-shopify">Tailwind labs is joining Shopify</a></p></li>
<li>
<p><a href="https://ai.meta.com/muse/">Meta releases Muse agent on the cloud</a> (only in US though)</p>
<p></p>
</li>
</ul>
<p></p>
<div><hr></div>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/809">Hackernewsletter</a> (#809th edition), and for software development/coding articles, join daily.dev.</p>
<p></p>
<p>That’s it from the 110th Edition of techstructive weekly. I hope you found it helpful, and relaxing. If not please drop any suggestions, feedback or discussion about certain things you want to in the comments or drop me a message on my <a href="https://www.meetgor.com/contact">socials</a>. </p>
<p></p>
<p>Thank you for reading,</p>
<p>Until next week.</p>
<p>Happy Coding :)</p>
