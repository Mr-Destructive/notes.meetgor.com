---
title: "Techstructive Weekly #97"
date: 2026-06-05T18:46:02Z
slug: techstructive-weekly-97
draft: false
type: post
description: "Read about OpenAPI specifications, watched a few livestreams on understanding API designing, among the other things learnt, explored from the week of 31st May to 6th June 2026"
tags: ["substack"]
---


<h2>Week #97</h2>
<p>It was a great week. A hope kindled out of random or maybe out of pity. I appreciate it with both hands, and don’t want to throw it away. Maybe I was unlucky for the past couple of months, but when it strikes, I have to be ready with both hands, when a opportunity comes, the preparation is what makes your luck work.</p>
<p>I will be on a roll on the weekend for exploring APIs and deeply understanding OpenAPI specifications, I know the fundamentals, have developed and maintained a few over the last couple of years, but need to make it to the next level. I think I have been ignorant of my interests and curiosities, not sure why I let the time get wasted. But anyhew, we need to move on, this week has surely inspired and bought back the spark that I have been missing since maybe a dozens of weeks. AI is over, I know what it can do and cannot. Enough of the hype, its time to get the steering wheels back.</p>
<p></p>
<h3>Quote of the Week</h3>
<blockquote>
<p>“In dreams begin responsibility”</p>
<p>— <a href="https://www.goodreads.com/quotes/97605-in-dreams-begin-responsibilities">William Butler Yeats</a></p>
<p>mentioned in “Kafka on the Shore” By haruki murakami</p>
</blockquote>
<p>I was reading this book “Kafka on the Shore” by Haruki Murakami, and this paragraph just sticks out. There are so many quotes from this book, so many thoughts I want to pen and paper to, but taking one bit at a time. This one is deep. Dreams do bring responsibilities, flip this around and you could say that where there’s no power to imagine, no responsibility can arise. Wow!</p>
<p>Isn’t it fascinating. I mulling over this quote which was mulled over by Haruki Murakami, through the lens of Kafka. Its just like a thread of thought is shared between the author and the people he is thinking of and me the reader. Just cannot put word to it, but its powerful thing to be able to do that. Reading does wonders to the mundane life, inspires and unclogs the brain from negative thoughts.</p>
<p></p>
<div><hr></div>
<h2>Read</h2>
<ol>
<li>
<p><a href="https://caio.ca/blog/a-love-letter-to-neovim">Caio Bianchi: A love Letter to Neovim</a></p>
<ul>
<li><p>I can feel the words here. It could be AI but I not 100% sure of that. It looks authentic in some angles, though some analogies do remind me of those sloppy comparisons.</p></li>
<li><p>But the gratitude and empathy for the editor is what I like about the article. It just appreciates how available and freeing is to use Neovim.</p></li>
<li><p>It just works. It doesn’t break with some updates (unless you run packer sync amidst mid-meeting). You are confident that it will work anytime you do certain things, type certain things exact same results. It feels smooth and doesn’t take eternity to load and boot up. I love it too.</p></li>
</ul>
</li>
<li>
<p><a href="https://www.seangoedecke.com/weird-projects-i-shipped-with-ai/">Sean Goedecke: Wired Projects I shipped with AI</a></p>
<ol>
<li><p>Man that is so cool. He literally came up with Endless Wiki (which was copied by xAI as Grokipidia, sad life). So he is AI maximalist. But in a good way I think.</p></li>
<li><p>He knows stuff, what would work and what won’t right. </p></li>
<li><p>He also clears the confusion that I constantly ask myself , “if LLMs are so good at writing code, where is the tsunami of new AI-generated apps, services and games?” he responds by saying “Writing code is only one of the bottlenecks involved in actually shipping a new product, after all”. There are other bottlenecks like decision making and the most bugging question for developers. “Why to build it”, what is a imagination and creativity problem and how , llms have taken over that.</p></li>
</ol>
</li>
<li>
<p><a href="https://lornajane.net/">Lorna Jane’s Blog:</a> This blog is a gold mine for API lovers. It also covers a lot of developer relations and dev-ex. Really a perfect combination for something I want to read and talk about.</p>
<ol>
<li>
<p><a href="https://lornajane.net/posts/2025/save-edits-to-openapi-as-an-overlay">Overlays in OpenAPI</a> and <a href="https://lornajane.net/posts/2024/openapi-overlays-to-avoid-api-oversharing">avoiding oversharing APIs</a></p>
<ol>
<li><p>This is a gem. We can trim off or add certain things to the spec in certain points of the pipeline of the API Spec. </p></li>
<li><p>This is so powerful for removing internal APIs and only exposing user facing ones, while still having the rich documentation for humans (now AI agents) to understand the actual api underneath.</p></li>
<li><p>Being able to do this programatically is amazing and makes a API pipeline so controllable and rich at the right place and time.</p></li>
</ol>
</li>
<li>
<p><a href="https://lornajane.net/posts/2019/using-tags-in-your-openapi-spec">Using Tags in OpenAPI Spec</a></p>
<ol>
<li><p>I was confused as to what the use of Tags was actually, might be for listing certain groups of related api. I didn’t think of grouping in the documentation page and elsewhere. I was just puzzled as to where it could probably be helpful. This article just made it clear.</p></li>
<li><p>“Quite a few of the documentation tools will render the endpoints grouped by tag which is very valuable”. That is what I wanted to hear. Nobody made it quite obvious.</p></li>
<li><p>Now I see the header like sections, I can now think, those are the tags, yes, so the APIs inside those tags might be related.</p></li>
</ol>
</li>
<li>
<p><a href="https://lornajane.net/posts/2020/whats-new-in-openapi-3-1">What’s new in OpenAPI 3.1 Spec</a></p>
<ol>
<li><p>The type having a array of possibilities is so cool. We can basically have a union like datatype in python here. Powerful.</p></li>
<li><p>The webhook thing, i must say it always confuses and frightens me. I always associate webhook to clients, I get scared. But this article clears that too</p></li>
<li><p>Suppose we are writing the Stripe’s API Spec. The customer of stripe, some fancy startup needs a event to send to the user when it happens, so the URL is defined by the user. So we don’t add the URL, we specify the event name (by convention not necessarily enforced) and the method. If there is a method then there must be request and response to provide those.</p></li>
<li><p>There are other changes as well like the override in ref and info block has a summary field now. </p></li>
<li><p>By the way, there is a 3.2.0 released as well https://github.com/OAI/OpenAPI-Specification/releases/tag/3.2.0 by the time of this post. It feels good that there is still innovation and releases happening at for granted taken things like APIs.</p></li>
</ol>
</li>
<li>
<p><a href="https://lornajane.net/posts/2023/tips-for-better-documentation-with-openapi">Tips for better documentation with OpenAPI</a></p>
<ol>
<li><p>Summary and description are really undervalued. I cannot love a API if it doesn’t have relatable and to  the point humanly descriptions and summary, even a simple one liner helps a ton.</p></li>
<li><p>Does than will do is really a good thing to keep in mind to build trust</p></li>
</ol>
</li>
</ol>
</li>
<li>
<p><a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods">MDN Docs for HTTP Methods</a>: I didn’t paid closer attention to these HEAD and OPTIONS method, but now all of them make sense. These are utility methods, nothing for data, but for communication meta information.</p>
<ol>
<li>
<p><a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods/OPTIONS">OPTIONS</a></p>
<ol>
<li><p>This is the request that asks for the available options in a API URI endpoint. It receives the response in response headers only, there is no response for a OPTIONS request. </p></li>
<li><p>It usually has a 204 No content response status and rightly so, since its purpose is to make the caller(client) aware which methods are available on the API.</p></li>
<li>
<p>It gets the headers like this `</p>
<pre><code><code>Access-Control-Allow-Methods: GET, POST, PUT, DELETE</code></code></pre>
</li>
</ol>
</li>
<li>
<p><a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods/HEAD">HEAD</a></p>
<ol>
<li><p>The header is powerful to get only the headers and not the actual response.</p></li>
<li><p>It is usually used to check if a resource exists, save bandwidth. It might still fetch the resource till the request leaves the server.</p></li>
<li><p>Pretty handy thing to know that this exists and can even use and implement one.</p></li>
</ol>
</li>
</ol>
</li>
<li>
<p><a href="https://github.com/github/rest-api-description">GitHub’s REST API Spec</a></p>
<ul><li><p>That came as a surprise a 3 MB file with references! What a rain of API Endpoints and Components. And the dereferenced one is 10 MB. WOW! I haven’t seen text specification that huge.</p></li></ul>
</li>
</ol>
<p></p>
<p></p>
<h2>Watched</h2>
<p></p>
<ul><li>
<p><a href="https://youtu.be/IvRmwiYdA_A">API Evangelist Conversation with Lorna and Kin Lane</a></p>
<ul>
<li><p>Overlay is a big thing that happened in the OpenAPI world. Rightly so. It creates a bridge between everyone by removing (mostly) the things that are not needed and keeps them for certain kinds of people (like developers those who work on integrations). </p></li>
<li><p>Developers also no longer need to clutter their core API description with language-specific metadata. Instead, these details can be applied as an overlay right before the SDK generation step that is neat little thing.</p></li>
<li><p>Kind of exciting that such a “boring” thing (I dont think its boring at all, just people have started to think in comparison to AI-tech) is still in progress and never stopping or getting discouraged by it. Gives me a lot of hope too.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"IvRmwiYdA_A","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/Il5btHG_D74">Hands on with Splotlights’ Spectral: API Linting for better Specs and Governance</a></p>
<ul>
<li><p>This was a great 1 hour video, really covered the most important thing to note about spectral, you can define your rules. That sometimes is a bit wired, right? Because what’s the point of the specification then, the point is that not everyone can accommodate the rules in every scenario, and more importantly the OpenAPI Spec is just the surface that brings the producer an the consumer on the same ground, so rules must be shared upon.</p></li>
<li><p>I immediately setup this tool with my Neovim config and it works like a charm.</p></li>
<li>
<p>Having the ability to write programmatic rules with js is so underrated. It opens a infinite possibilities.</p>
<p></p>
<div class="youtube-wrap" data-attrs='{"videoId":"Il5btHG_D74","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p> </p>
</li>
</ul>
</li></ul>
<p></p>
<ul><li>
<p><a href="https://youtu.be/-kQRBkBuwHg">Beyond the AI Hype</a></p>
<ul>
<li><p>This is a great video coming from the Indian background. It gives me hope that there are believers and people who truly understand the technology in and out and don’t play the hype.</p></li>
<li><p>The subject-matter-expert requirement is so crucial, the term “You create it, you only test it” is baseless as rightfully said. LLMs don’t have your conception of “right”. The professional who has spent years and decades learning it does.  Right now we are in the boom phase where the cost is lower and looks lucrative enough to make the bare minimum thing possible.</p></li>
</ul>
</li></ul>
<p></p>
<div class="youtube-wrap" data-attrs='{"videoId":"-kQRBkBuwHg","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://youtu.be/JpSHyEIZ_bo">Deep SWE Bench</a></p>
<ul>
<li><p>Ok, this is cool. I thought the same, if the benchmarks were open, won’t any lab just gamify it? Just like Meta did with LLama 4 I think.</p></li>
<li><p>Though the critique of making something closed sourced besides the open source test suite is important. </p></li>
<li><p>I don’t know if GPT models are that good, like I haven’t use Codex since a month, and not sure if I am missing on it. Claude code is good but the 5 hour window limit is so annoying.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"JpSHyEIZ_bo","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://www.youtube.com/live/QiAXxaLrt7E">JSON Schema Tips and Tricks</a></p>
<ul>
<li><p>This was a good live stream, from the maintainer of the JSON Schema org itself is so cool.</p></li>
<li><p>The ajv usage for building linting rules is like the right piece to build a tool. Like import a validator engine and write set of rules to make it your way, really nice.</p></li>
<li><p>It didn’t went into the JSON Schema specific details which I wanted to understand more in depth, but it was a good fundamental clearing guide on validation.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"QiAXxaLrt7E","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<ul><li>
<p><a href="https://www.youtube.com/live/JEBd78U9aBo?">OpenAPI 3.0 from idea to implementation</a></p>
<ul>
<li><p>Another great video from Postman team. They really love APIs don’t they.</p></li>
<li><p>It was a amazing refresher to OpenAPI, I haven’t looked into API specfications in a couple of months now, found it the right way to dab into the fundamentals and brush up the concepts. </p></li>
<li><p>AI Psychosis is real, it makes you forget certain things that you knew were interesting. I learn all of this first in 2024 and found I knew most of them since I have worked with them in the past couple of years. Be it writing the openapi spec from scratch for all the apis at work for a product company or updating and revamping an API and its specs too. It was a time which I loved to work with but Agents came and hijacked that urge to learn more. Here I am back on track.</p></li>
</ul>
</li></ul>
<div class="youtube-wrap" data-attrs='{"videoId":"JEBd78U9aBo","startTime":null,"endTime":null}' data-component-name="Youtube2ToDOM"><div class="youtube-inner"></div></div>
<p></p>
<h2>Learnt</h2>
<ul>
<li>
<p>HTTP Head and Options Methods</p>
<ul>
<li><p>I have written posts for the 5 methods  in Golang series(GET, POST, PATCH, PUT, DELETE). But never covered these two. I don’t think its covered much. I would like to write about it now. Weather people read it or  llm will read it or not, I will write it.</p></li>
<li><p>HEAD as explained above is for getting the header only part of a response without the server sending the actual response.</p></li>
<li><p>OPTIONS is for exploration of the API endpoint, it returns the header only, specifying what methods are available.  </p></li>
</ul>
</li>
<li>
<p>OpenAPI Specifications 3.1</p>
<ul>
<li><p>I didn’t knew it was still being developed and maintained with more additions. It felt really a breeze of fresh air, that we all need in the world of AI.</p></li>
<li><p>But the point being, that it added support for the API to convey that it will send webhook for a named event with a specific set of request and response for that given method. The name of the webhook is important because that is the only differentiator from the other webhooks if any (you can skim through the request and response body otherwise)</p></li>
<li><p>The tiny little additions is what makes me feel good, that there are people still trying to build and improve things.</p></li>
</ul>
</li>
<li>
<p>Spectral/YAML Lint Setup on Neovim</p>
<ul>
<li><p>I added vaccum and spectral to the LSP setup, and wow it made writing specs so easy and informative. </p></li>
<li><p>I am looking forward to doing livestreams by exploring these openapi specs and writing my own. I have a grand idea for openapi learning and mastery.</p></li>
</ul>
</li>
</ul>
<p></p>
<h2>Tech News</h2>
<p>I tend to take a break from this tech news, its a lot of paced things. Maybe this week was nothing noteworthy of catching my attention. But surely there is something brewing up and released. I am consciously aware that its moving too fast, and maybe I don’t want to keep up with it, I just cannot. I will read some articles if I find it interesting here. This week, nothing to speak and list about.</p>
<p>For more news, follow the <a href="https://buttondown.com/hacker-newsletter/archive/796">Hackernewsletter</a> (#796th edition), and for software development/coding articles, join daily.dev.</p>
<div><hr></div>
<p>That’s it from the week. Will be out on Sunday. Saturday will try to learn as much as I can about OpenAPI and API designs. I really had a itch from the beginning of my software career to devvelop APIs, but never really sat and spent months on it. Life is a bit tough, but focus and discipline is something that can get you what you want to.</p>
<p>Hope your life is good, and if not, I wish it would be soon. have hope and keep learning, keep the flame of curiosity alive, it will eventually play out.</p>
<p>Thanks for reading. See you next week!</p>
<p>Happy Coding :)</p>
<p class="button-wrapper" data-attrs='{"url":"https://techstructively.substack.com/p/techstructive-weekly-97/comments","text":"Leave a comment","action":null,"class":null}' data-component-name="ButtonCreateButton"><a class="button primary" href="https://techstructively.substack.com/p/techstructive-weekly-97/comments"><span>Leave a comment</span></a></p>
<div class="captioned-button-wrap" data-attrs='{"url":"https://techstructively.substack.com/p/techstructive-weekly-97?utm_source=substack&amp;utm_medium=email&amp;utm_content=share&amp;action=share","text":"Share"}' data-component-name="CaptionedButtonToDOM">
<div class="preamble"><p class="cta-caption">Thanks for reading Techstructive Weekly! This post is public so feel free to share it.</p></div>
<p class="button-wrapper" data-attrs='{"url":"https://techstructively.substack.com/p/techstructive-weekly-97?utm_source=substack&amp;utm_medium=email&amp;utm_content=share&amp;action=share","text":"Share"}' data-component-name="ButtonCreateButton"><a class="button primary" href="https://techstructively.substack.com/p/techstructive-weekly-97?utm_source=substack&amp;utm_medium=email&amp;utm_content=share&amp;action=share"><span>Share</span></a></p>
</div>
<div class="subscription-widget-wrap-editor" data-attrs='{"url":"https://techstructively.substack.com/subscribe?","text":"Subscribe","language":"en"}' data-component-name="SubscribeWidgetToDOM"><div class="subscription-widget show-subscribe">
<div class="preamble"><p class="cta-caption">Thanks for reading Techstructive Weekly! Subscribe for free to receive new posts and support my work.</p></div>
<div class="fake-input-wrapper">
<div class="fake-input"></div>
<div class="fake-button"></div>
</div>
</div></div>
<p></p>
<p></p>
<p></p>
