---
type: til
title: Create a Non-Clustered Index in Django with Postgres as DB
description: Understanding how to add a non-clustered index in a postgres database
  in a django project.
status: published
slug: django-non-clustered-index-pg
date: 2022-11-10 22:10:00
tags:
- go
- python
- sql
- testing
---
## What is a non-clustered index?

A non-clustered index is a seperate structure than an actual table in the database, it stores the non-clustered index key(the column which we want to sort in the table), and a pointer to the actual values based on the index key. So, non-clustered indexes do not change the physical order of the table records, instead it holds a structure that can provide a easier and distinct way to fetch objects based on a particular column as the primary key in the structure.

## How to create a non-clustered index in django

In django, we can use the [db_index](https://docs.djangoproject.com/en/4.1/ref/models/indexes/) property on a field(s) to create a index on the table/model. 

### Add the property to the field in the model

Chose a field in which, you want to add a index. It can be a foreign key or any other normal field defined in your model.

We have used the typical blog model, so used in the some of my [TILS](https://www.meetgor.com/tils/) in django, it is just convenient to explain and understand as well. We have a django project named `core` and it has a app `blog` with a model defined below. The model `Article` has a few attributes like `title`, `description`, `content` and `status`.

```python
from django.db import models

ARTICLE_STATUS = [
    ("PUBLISHED", "Published"),
    ("DRAFT", "Draft"),
]

class Article(models.Model):
    title = models.CharField(max_length=128, db_index=True)
    description = models.CharField(max_length=512)
    content = models.TextField()
    status = models.CharField(max_length=16, choices=ARTICLE_STATUS, default="DRAFT")

    def __str__(self):
        return self.title
```

So, we have added a `db_index` to the title column in the model as a property. This will be equivalent to creating a index in `SQL` as follows:

```
$ python manage.py makemigrations

Migrations for 'blog':
  blog/migrations/0002_alter_article_title.py
    - Alter field title on article
```

```
$ python manage.py migrate

Operations to perform:
  Apply all migrations: admin, auth, blog, contenttypes, sessions
Running migrations:
  Applying blog.0002_alter_article_title... OK

```

Indexes are not standard as in SQL, but each vendor(sqlite, postgres, mysql) have their own flavour of syntax and naunces.

```sql
CREATE INDEX "blog_article_title_3c514952" ON "blog_article" ("title");

CREATE INDEX "blog_article_title_3c514952_like" ON "blog_article" ("title" varchar_pattern_ops);
```

The above index commands are specific to the field, as the title field is a varchar, it has two types of index, it can generate one with simple match and other for `LIKE` comparisons because of string comparison behaviour.

So, we just created a simple index and now if we query the db for a particular `title` which now has its own index for the table `blog_article`. This means, we will be able to fetch queries quickly if we are specifically filtering for `title`.

### Adding some data records

We can add a few data records to test the query from the databse, you can ignore this part as it would be just setting up a django project and adding a few records to the databse. This part won't make sense for people reading to get the actual stuff done, move to the next part please.

```
python manage.py createsuperuser
# Create a super user and run the server

python manage.py runserver
# Locate to http://127.0.0.1:8000/admin
# Create some records in the artilce model
```

So, after creating some records, you should have a simple database and a working django application.

```sql
SELECT * FROM blog_article;
```
```
blog_test=# SELECT * FROM blog_article;

 id |  title   | description |          content          |  status   
