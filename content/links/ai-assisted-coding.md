---
title: "AI Assisted Coding"
date: 2026-01-24
draft: false
---

# AI Assisted Coding

**Link:** https://youtu.be/5vp9ypOUgMw

## Context

1. [How to stop feeling lost in tech](https://www.yacinemahdid.com/p/how-to-stop-feeling-lost-in-tech))
    - This is a premium quality post. So much valuable advice for juniors and people who are all over the place and feeling overwhelmed or burned out.
    - What do you want to be in like 5 years, 3 years, 1 year? That hit me like a truck. I never think about the future, but still there needs to be specific I enjoy and would want to keep enjoying and learning. I know that's a too generic AI like answer but that is where the real friction will arise and I think we will emerge out of it with a fresh perspective. After having asked a lot of questions and having a better idea about our likes and dislikes, thinking without any pressure just gives us the opportunity to truly see beneath us.
    - I don't think the steps mentioned there are necessary, like I know this might be done for relaxing the user (taking leave in the middle of a week, going a entire day at the waffle-house or cafe). But the process is to the point.
    - Listing down the goals and what I want to be like a tree like structure, no cycles. Breaking down one goal from 5 year to 3, then 1 and then month and then weeks. That makes it really clear and forces us to double down on our true "calling" or the muse if you think in that way.
2. [What makes a good software engineer](https://sanitarium.se/blog/2025/08/21/what-makes-a-good-software-engineer)
    - Curiosity to build the right thing, and being honest about it.
    - Trust is super valuable, I have experienced it working remotely at my 1+ year stay at my current company.
    - I often just do it for doing it, then I hit it with I need to know why this needs to be done and done correctly. Then things take shape, curiosity out performs every other emotion.
3. [Go is still not good](https://blog.habets.se/2025/07/Go-is-still-not-good.html)
    - I can agree to only one points here especially the standard library swallowing exceptions part. That sometimes is little on the border of Go's philosophy vs the actual issue. Just read the f-ing manual would be the argument, but it kind of becomes too verbose and critical if mishandled.
4. [Exploring EXIF](https://hturan.com/writing/exploring-exif)
    - A great post, diving into a specific tool for a lot of things.
    - I think the media metadata is messed up, the software is not consistently adhering to the standards. Look at the blog post and the Apple photos have a lot more metadata then Google or any other photo taking software.
    - I have to say the different of fields are in orders of magnitude more than the rest. How much additional info they cram into a single photo is bewildering.
    - I really liked the post, it had dumping of data into sqlite database and writing queries to get certain things, how cute and interesting that idea is. My brain is running in all directions at this now.
    - LLMs to write queries to get photos from albums with natural language and metadata without actually processing the image, is one I can hit straight off my head.
    - I confirmed that it Apple who is pouring love into those photos. It probably has to do a lot of things in the background to get the metadata. Such a irony of tech today. Lagging behind the trend but way ahead of everyone in the core.
5. [How AI writing supercharged an addictive pattern](https://bowendwelle.substack.com/p/ai-writing-addictive)
    - We'll have been there and done that. The initial WOW! to the glorious vibe crap that comes after 12 hours of prumpting.
    - I know there is a lot of value in AI, I am not a AI doomer, but the things people are assuming like 10x dev, replacing a junior dev, an intern, I don't buy that both emotionally and mentally. It just hurts the morale of humans to say that. If that is the case, then just flip the lever and let the software world be in chaos.
6. [5 Docker networking concepts that everyone should know](https://medium.com/@vishnugoswami6000/5-docker-networking-basics-every-beginner-should-actually-understand-7a558a6c5c0a)
    - This was a fantastic walk-through of one of the un-explained parts of docker. Surely there are exhaustive lists of commands to perform things in the docker networking, but very few dive into the why part. This article specifically did that and all the commands made sense.
    - The bridge is the safe default, setting custom host is for advanced usage, avoid if possible, only dive if necessary.
    - Custom networks makes communication clear between containers
      user names, not IPs.
7. [Left to Right Programming](https://graic.net/p/left-to-right-programming)
    - If python looks odd, how about SQL? That is left to right but evaluated mostly from the right to left right? mostly I asking loosely here, though aggregates might be exceptions and other things I am not aware of.
      > Programs should be valid as they are typed.
    - Good points, yes python’s list comprehension is an outlier at it as it needs the last bit first to understand the whole context.
8. [No AI is not making Engineers 10x productive](https://colton.dev/blog/curing-your-ai-10x-engineer-imposter-syndrome/)
    - See the copium? Too much bearish on AI now. Suddenly people realise they need to put breaks to the hype and focus on improving the craft before cracks start forming into dents and large holes.
      > Making all your engineers feel constantly anxious about their performance is *bad for your company*. It will make your engineers not want to work for you. This is a recipe for short term thinking that will encourage engineers to max out bad metrics, like lines of code. Code review will get neglected, tech debt will compound, and in the long term the whole company will be footing the bill of those errors
    - What a relief, to hear or sorry read it. Every word of this article is healing me. I will write a separate thought on this article.
9. [LLMs makes us dumber in the longer run](https://desunit.com/blog/in-the-long-run-llms-make-us-dumber)
  > The comfort we get when offloading our cognitive load to LLMs is bad for us. Cognitive load should exist, and if we reduce it too much – if we stop thinking – we can actually unlearn how to think.
    - Well said. Consistent reliance on AI is and will distinguish some human from an honest caring human.
10. [What's the point of vibe coding if I still have to pay a dev to fix](https://old.reddit.com/r/vibecoding/comments/1mu6t8z/whats_the_point_of_vibe_coding_if_i_still_have_to)it
    - People are coping seriously on AI now. Not that AI is bad, it's not yet there to build software blindly and let laymen control the steering. Devin, you tried, unfortunately, we won't be moving forward with your application.
11. [Why building my blog is more fun than filling it](https://clojurecivitas.github.io/civitas/why/i_should_write_more.html)
    - Its a rabbit hole. I have done it and seen people do it too. Waylon Walker's Markata (its a python plugin based ssg, so easy to extend and setup) and my static site generator (I am calling it burrow, inspired as a home for golang gophers). I used jekyll > markata and felt the calling for writing my own ssg and here I am.
    - It feels rewarding to show the world what you have built for yourself and brag about it when someone sees it. It's like a garden, what is writing equivalent in gardening? Planting trees, that's how the blog gets populated, and adding different things to the garden is what is addictive and feels like a gardener.
12. [Do blogs need to be so lonely](https://thehistoryoftheweb.com/do-blogs-need-to-be-so-lonely/)
    - Collaborative blogging, I think hashnode tried to do it with community blogs and co-authoring. I like it but I feel that it's quite rare to collaborate on shared interesting topics, a podcast or a video interview might be a better fit to show the communication. Blogs don't really fit that quite well in my opinion, but yes distinct parts of the blog could be collaborated.
13. [Stupid things that work](https://ryanglover.net/blog/stupid-things-that-work)
    - I mean why not? looping one billion times in javascript makes your computer warm, fine. I would rather install and run Android studio than touch javascript. These days, running a local model around 3-5 Billion parameters (on 8GB ram) could easily burn your computer not just warm it.

**Source:** techstructive-weekly-56
