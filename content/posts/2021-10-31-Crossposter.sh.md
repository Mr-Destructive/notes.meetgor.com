---
type: post
title: "Crossposting with a single script: Crossposter.sh"
subtitle: "Crosspost on dev.to, medium.com and hashnode.com with a single click"
date: 2021-10-31T15:45:00+05:30
tags: ["linux", "neovim", "testing"]
  - bash
slug: crossposter-shellscript
image_url: https://res.cloudinary.com/dgpxbrwoz/image/upload/v1643288989/blogmedia/trssl38erkdbcqlnjdvp.png
---

## Introduction

If you have been writing articles you know the pain to get some attention, if you have already been cross-posting your articles it usually takes some time to do that. This task can be automated with a shellscript. If you have been cross-posting articles on `medium.com`, `dev.to` and at `hashnode.com`, then I have a treat for you. 

Introducing **crossposter.sh**!!

## What is Crossposter.sh?

### Crosspost to dev.to/hahsnode/medium from the command line.

Crossposter.sh is a shellscript(BASH) to automate crossposting to platforms like dev.to, medium.com and hashnode.com. The script takes in markdown version of your post with a few inputs from you and posts it to those platforms. You would require a token/key for each of those platforms to post it from the command line. You can check out the official repository of [Crossposter](https://github.com/Mr-Destructive/crossposter).

The actual script is still not perfect (has a few bugs). Though it posts on `dev.to` and `medium.com` easily, the `hashnode.com` is buggy as it parses the raw markdown into the post and doesn't render as desired. So, **its a under-development script**, fell free to raise any issues or PRs on the official GitHub repo.   

Run the script on a bash interpreter with the command:

`bash crosspost.sh`

For posting the article you need to provide the following details:

## Front-Matter

### Meta data about the post

- Title of Post
- Subtitle of Post
- Publish status of post(`true` or `false`)
- Tags for the post (comma separated values)
- Canonical Url (original url of the post)
- Cover Image (URL of the post's image/thumbnail)

This information is a must for `dev.to` especially the `title`. This should be provide in the same order as given below:

```yaml

