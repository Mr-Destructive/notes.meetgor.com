---
type: post
title: 'Django Basics: Setup and Installation'
subtitle: Setting up environment and installing Django framework
date: 2021-11-20 16:30:00+05:30
slug: django-basics-setup
series:
- Django-Basics
image_url: https://res.cloudinary.com/dgpxbrwoz/image/upload/v1643290071/blogmedia/s8ahlep1e8lmgiboyjhz.png
tags:
- python
---
## Introduction

The crucial aspect of starting to learn any framework is the ease to set it up and Django by far is the easiest of the options out there. There is just a few lines of code to install django if you already have python installed in your system. In this article, we see how to setup a django project along with a virtual environment. 

If you already have python and pip installed, you can move on to the [virtual environment setup](#setting-up-virtual-environment-in-python).

## Installing Python and PIP

Django is a python based framework so that makes sense to have Python installed along with its package manager to use Django.  

To install Python, you can visit the official [Python](https://www.python.org/downloads/) website to download any relevant version for your system (recommended 3.7 and above). 

Mostly the Python installation comes with the option to install `pip`(python's package manager) but if you missed that, that's fine, you can install the [get-pip.py](https://bootstrap.pypa.io/get-pip.py) file into your system and run the below code:

```
python get-pip.py   
```

Make sure the include the relative path to the file if you are not in the same folder as the file.

So, that should be python setup in your local machine. To check that python was installed correctly, type in `python --version` and `pip --version` to check if they return any version number. IF they do, Congratulations !! You installed Python successfully and if not, don't worry there might be some simple issues that can be googled out and resolved easily. 
   
