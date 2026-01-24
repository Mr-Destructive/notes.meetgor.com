---
type: post
title: 'Comment/Uncomment Code: Vim for Programmers'
subtitle: A short guide to comment/uncomment chunks of code effectively in Vim
date: 2021-10-07 16:45:42+05:30
slug: vim-un-comment-p1
image_url: https://res.cloudinary.com/dgpxbrwoz/image/upload/v1643288865/blogmedia/jbs2hcaiplfvwsgrcfgl.png
tags:
- frontend
- go
- linux
- neovim
- python
- rust
- testing
- vim
---
## Introduction

We as programmers always fiddle with commenting out code for code testing, documenting the function of code, and most importantly debugging. So you can't wait to comment on a large chunk of code manually, as it is quite a tedious thing to do. Let's do it effectively in Vim.

In this part of the series, I'll cover how to comment/uncomment chunks/blocks of code effectively in Vim. We will see and use some commands, keybindings for doing so, and also we would add certain components to our vimrc file as well to design some custom key mappings.  Let's get faster with Vim.
  
## How to comment multiple lines effectively

To comment on multiple lines of code, we can use the Visual Block mode to select the lines, and then after entering into insert mode, we can comment a single line and it would be reflected on all the selected lines.

1. Press `CTRL+V` and Select the line using j and k

2. After Selecting the lines, Press `Escape`

3. Press `Shift + I`, to enter insert mode

4. Enter the comment code (`//`, `#`, or other)


![vimcoment.gif](https://cdn.hashnode.com/res/hashnode/image/upload/v1633518136135/06dfBTq2T.gif)

So, using just simple steps you can comment out large chunks of code quite easily and effectively. If you are using some other language that has multiple characters for commenting like `//`, `- -`, etc, you can type in any number of characters while being in insert mode after selecting the lines.

 
![vimcppcom.gif](https://cdn.hashnode.com/res/hashnode/image/upload/v1633520509953/0q-k2ZHC7.gif)

This might look a bit wired on the first try but just try it every day, It is a life-saving and very satisfying experience once applied in a real-world scenario.


## How to uncomment multiple lines effectively

Now, as we have seen to comment out a large chunk of code, we can even uncomment the code very easily. It's even simpler than commenting the code.

1. Press `CTRL + V` to enter Visual Block mode

2. Select the commented characters

3. Press `d` to delete the comments

4. Press `Escape`

![vimuncoment.gif](https://cdn.hashnode.com/res/hashnode/image/upload/v1633518156818/GJzRPTI3I.gif)

We can simply use the CTRL + V to select the comment, and then press d to delete all the comment characters. 

**We are using the Visual Block mode as we only want the comment to be selected and not the entire code associated with the lines.**

## Using Multiline Comments for Programming languages

Now you might say, why use multiple single-line comments when we can use multiline comments in almost all programming languages. Well, Of course, you can do that, it's easier for reading the code if syntax highlighting is accurate and greys out the commented part. We can simply add those characters to the start of the block and at the end of the block.  

But in Vim, we can customize that too, just imagine when you just select the chunk/block of code that you need to comment out and then simply press a few keystrokes (just 2) and the multiline comments are automatically (programmatically) added as per the programming language extension of the file.

Isn't that cool? Well, you just need to copy-paste the below code to your Vimrc file and source it and you are good to go. 

```vim
function! Comment()
    let ext = tolower(expand('%:e'))
    if ext == 'py' 
        let cmt1 = "'''"
	    let cmt2 = "'''"   
    elseif ext == 'cpp' || ext =='java' || ext == 'css' || ext == 'js' || ext == 'c' || ext =='cs' || ext == 'rs' || ext == 'go'
	    let cmt1 = '/*'
	    let cmt2 = '*/'
    elseif ext == 'sh'
	    let cmt1 = ":'"
	    let cmt2 = "'"
    elseif ext == 'html'
	    let cmt1 = "<!--"
	    let cmt2 = "-->"
    elseif ext == 'hs'
	    let cmt1 = "{-"
	    let cmt2 = "-}"
    elseif ext == "rb"
	    let cmt1 = "=begin"
	    let cmt2 = "=end"
    endif
    exe line("'<")."normal O". cmt1 | exe line("'>")."normal o". cmt2 
endfunction

function! UnComment()
    exe line("'<")."normal dd" | exe line("'>")."normal dd"   
endfunction


vnoremap ,m :<c-w><c-w><c-w><c-w><c-w>call Comment()<CR>
vnoremap m, :<c-w><c-w><c-w><c-w><c-w>call UnComment()<CR>

```
The below screencast is an example of `HTML` snippet in a file that is getting commented using mapping with the keys `,m` you can put any other keybinding you like. 
![htmcm.gif](https://cdn.hashnode.com/res/hashnode/image/upload/v1633595891674/hbhrbtRHd.gif)

