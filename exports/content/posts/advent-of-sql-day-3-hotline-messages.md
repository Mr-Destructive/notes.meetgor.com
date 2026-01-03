---
title: "Advent of SQL Day 3-Hotline Messages"
date: 2025-12-18
slug: advent-of-sql-day-3-hotline-messages
type: post
draft: false
description: ""
tags: []
---

This is day 3 from the Advent of SQL
## Grab the SQL StatementsLet's take the insert statements i.e. to create and populate tables and rows into the database. I am using SQLite.
It works without any special shenanigans, as it was intended to used for Postgres, but the table and use case looks very simple, so nothing specific to Postgres used yet! We are good!

Here is the SQL setup, if you want to play it in the playground here:
<pre><code class="language-sql">DROP TABLE IF EXISTS hotline_messages;

CREATE TABLE hotline_messages (
    id INT PRIMARY KEY,
    caller_name TEXT,
    transcript TEXT,
    tag TEXT,
    status TEXT
);

INSERT INTO hotline_messages (id, caller_name, transcript, tag, status) VALUES
    (1, 'Saanvi A.', 'I just found a refrigerator portal that leads to a disco party hosted by dancing llamas—please send help!', 'possible dragon', 'spam'),
    (2, 'Fatima Q.', 'Hi Santa, I would love a magical unicorn that lights up at night!', 'wish list', NULL),
    (3, 'Lillian Z.', 'Hi Santa, I would love the magical fairy garden set, please!', 'wish list', 'approved'),
    (4, 'Carter Y.', 'Thank you, Santa, for making Christmas so special with your wonderful spirit!', 'thank you', 'approved'),
    (5, 'Omar R.', 'Hi Santa, I would love a rainbow unicorn plushie that has a glittery horn!', 'wish list', 'approved'),
    (6, 'Diego Y.', 'Hi Santa, I would love a magical unicorn plushie that glows in the dark!', 'wish list', NULL),
    (7, 'Layla X.', 'Thank you, Santa, for spreading joy and magic every Christmas!', 'thank you', NULL),
    (8, 'Sophia K.', 'Santa, my cat said she wants to visit the candy cane forest next week.', NULL, NULL),
    (9, 'Eli H.', 'Hi Santa, I would love the magical fairy castle with twinkling lights!', 'wish list', 'approved'),
    (10, 'Logan F.', 'Santa, I think the reindeer are starting a band with the garden gnomes.', 'needs clarification', NULL),
    (11, 'Carlos P.', 'Thank you, Santa, for making Christmas so special every year!', 'thank you', NULL),
    (12, 'Zain G.', 'Thank you, Santa, for bringing joy to all the children around the world!', 'thank you', NULL),
    (13, 'Haruto R.', 'Thank you, Santa, for spreading so much joy and magic every Christmas!', 'thank you', 'approved'),
    (14, 'Oliver L.', 'Thank you, Santa, for spreading joy and making Christmas extra special!', 'thank you', NULL),
    (15, 'Luca M.', 'Hi Santa, could I please have the super cool glow-in-the-dark rocket ship?', 'wish list', NULL),
    (16, 'Samuel C.', 'sorry, Santa, my teddy bear said he wants to be a reindeer this year.', 'needs clarification', NULL);
</code></pre>
Here's the setup I did to check the data.
```
`SELECT * FROM hotline_messages;`
```
<pre><code>$ sqlite3
SQLite version 3.45.1 2024-01-30 16:01:20
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.

sqlite> .read day3-inserts.sql

sqlite> .schema
CREATE TABLE hotline_messages (
    id INT PRIMARY KEY,
    caller_name TEXT,
    transcript TEXT,
    tag TEXT,
    status TEXT
);

sqlite> .mode table
sqlite> SELECT * FROM hotline_messages LIMIT 10;
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| id | caller_name |                          transcript                          |         tag         |  status  |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 1  | Saanvi A.   | I just found a refrigerator portal that leads to a disco par | possible dragon     | spam     |
|    |             | ty hosted by dancing llamas—please send help!                |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 2  | Fatima Q.   | Hi Santa, I would love a magical unicorn that lights up at n | wish list           |          |
|    |             | ight!                                                        |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 3  | Lillian Z.  | Hi Santa, I would love the magical fairy garden set, please! | wish list           | approved |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 4  | Carter Y.   | Thank you, Santa, for making Christmas so special with your  | thank you           | approved |
|    |             | wonderful spirit!                                            |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 5  | Omar R.     | Hi Santa, I would love a rainbow unicorn plushie that has a  | wish list           | approved |
|    |             | glittery horn!                                               |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 6  | Diego Y.    | Hi Santa, I would love a magical unicorn plushie that glows  | wish list           |          |
|    |             | in the dark!                                                 |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 7  | Layla X.    | Thank you, Santa, for spreading joy and magic every Christma | thank you           |          |
|    |             | s!                                                           |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 8  | Sophia K.   | Santa, my cat said she wants to visit the candy cane forest  |                     |          |
|    |             | next week.                                                   |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 9  | Eli H.      | Hi Santa, I would love the magical fairy castle with twinkli | wish list           | approved |
|    |             | ng lights!                                                   |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
| 10 | Logan F.    | Santa, I think the reindeer are starting a band with the gar | needs clarification |          |
|    |             | den gnomes.                                                  |                     |          |
+----+-------------+--------------------------------------------------------------+---------------------+----------+
sqlite> </code></pre>
So, we just have one table, called `hotline_messages` and it has a few columns like
<ol><li>`caller_name`
</li><li>`transcript`
</li><li>`tag`
</li><li>`status`
</li></ol>What do we want to do with those?
Well! Let's get into the problem statement.

## ProblemHere goes the challenge for day 3
<blockquote>Using the `hotline_messages` table, update any record that has "sorry" (case insensitive) in the transcript and doesn't currently have a status assigned to have a status of "approved".
Then delete any records where the tag is "penguin prank", "time-loop advisory", "possible dragon", or "nonsense alert" or if the caller's name is "Test Caller".
After updating and deleting the records as described, write a final query that returns how many messages currently have a status of "approved" and how many still need to be reviewed (i.e., status is NULL).
</blockquote>
It's divided into 3 parts, so we need three queries? Maybe, I don't want to have a single long query for doing all of these. And after reading it, it seems it should not be a single query, it can be 2 queries, one is for updation and the other for selection after those updates.

So, we have to do three things.
<ol><li>Find the records which have `sorry` in the transcript text and mark their `status` as `approved` (What a lovely gesture)
</li><li>Find all records with the tags as either `penguin prank`, `time-loop advisory`, `possible dragon`, or `nonsense alert` and even if the `caller_name` is `Test Caller` then delete those records, yes take'em out of my way.
</li><li>After doing those 2 things, we have to count the number of records with `status` as `approved` and the number of records that are still not `approved` (they are in review or the status is `NULL` )
</li></ol>
So, let's do them step by step.
### Be generousLet's be generous like Santa says and mark the records with `status` as `approved` whose transcript have the word `sorry` in them. Let those children be gifted their reward of being generous and humble.
How do that in SQL, well let's first look at what we are updating.
```
`SELECT * FROM hotline_messages WHERE transcript LIKE '%sorry%';`
```So, will it be sufficient? I think so.
Because
<ul><li>`LIKE` is **case insensitive**, so it can catch `sorry`, `Sorry`, `SoRRY`, `sorrY`
</li><li>`%` before and after will catch the word `sorry` in middle of the sentence and not necessarily in the start.
</li></ul>I can see 104 rows selected with this condition. I always try to check before updation or deletion how many rows are affected. Because, sometimes we start `UPDATE hotline_messages SET status = 'approved'` and forget the where! This gets worse for delete believe me!
<pre><code class="language-sql">SELECT count(*) FROM hotline_messages WHERE transcript LIKE '%sorry%';
SELECT count(*) FROM hotline_messages;</code></pre>
<pre><code>sqlite> SELECT * FROM hotline_messages WHERE transcript LIKE '%sorry%';
+------+--------------+--------------------------------------------------------------+---------------------+----------+
|  id  | caller_name  |                          transcript                          |         tag         |  status  |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 16   | Samuel C.    | sorry, Santa, my teddy bear said he wants to be a reindeer t | needs clarification |          |
|      |              | his year.                                                    |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 33   | Jacob F.     | sorry, Hi Santa, I would love the magical unicorn plushie th | wish list           |          |
|      |              | at glows in the dark!                                        |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 46   | Jun Y.       | sorry, Hi Santa, I would love a magical unicorn stuffed anim | wish list           |          |
|      |              | al that glows in the dark!                                   |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+

sqlite> SELECT count(*) FROM hotline_messages WHERE transcript LIKE '%sorry%';
+----------+
| count(*) |
+----------+
| 104      |
+----------+
sqlite> SELECT count(*) FROM hotline_messages ;--WHERE transcript LIKE '%sorry%';
+----------+
| count(*) |
+----------+
| 1067     |
+----------+
sqlite> 
</code></pre>So, once I know `104` rows will be affected out of `1067` I can create the update statement.
We want to update the status and set it to `approved` for the rows which we selected just now (have `sorry` in the transcript text)

<pre><code class="language-sql">UPDATE hotline_messages
SET status = 'approved'
WHERE transcript LIKE '%sorry%';</code></pre>
Now, when we select again
```
`SELECT * FROM hotline_messages WHERE transcript LIKE '%sorry%'; `
```All approved!
<pre><code>sqlite> SELECT * FROM hotline_messages WHERE transcript LIKE '%sorry%';
+------+--------------+--------------------------------------------------------------+---------------------+----------+
|  id  | caller_name  |                          transcript                          |         tag         |  status  |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 16   | Samuel C.    | sorry, Santa, my teddy bear said he wants to be a reindeer t | needs clarification | approved |
|      |              | his year.                                                    |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 33   | Jacob F.     | sorry, Hi Santa, I would love the magical unicorn plushie th | wish list           | approved |
|      |              | at glows in the dark!                                        |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+
| 46   | Jun Y.       | sorry, Hi Santa, I would love a magical unicorn stuffed anim | wish list           | approved |
|      |              | al that glows in the dark!                                   |                     |          |
+------+--------------+--------------------------------------------------------------+---------------------+----------+</code></pre>To the next step then

### Remove SpamTo reiterate the second part of the challenge
<blockquote>Then delete any records where the tag is "penguin prank", "time-loop advisory", "possible dragon", or "nonsense alert" or if the caller's name is "Test Caller".
</blockquote>We basically need to 
Find all records with the tags as either `penguin prank`, `time-loop advisory`, `possible dragon`, or `nonsense alert` and even if the `caller_name` is `Test Caller` then delete those records, yes take'em out of my way.
So, again, select first update or delete later.
```
`SELECT * FROM hotline_messages WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert');`
```Here `IN` is a great helper as we can do equivalent of this into a compact statement.
<pre><code class="language-sql">SELECT * FROM hotline_messages 
WHERE 
    tag = 'penguin prank'
    OR tag = 'time-loop advisory'
    OR tag = 'possible dragon'
    OR tag = 'nonsense alert';</code></pre>That is a lot of `OR tag =` that is saved by `IN` a list of values. Handy little operator.
The count here is `68` 
```
`SELECT count(*) FROM hotline_messages WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert');`
```<pre><code>sqlite> SELECT count(*) FROM hotline_messages WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert');
+----------+
| count(*) |
+----------+
| 68       |
+----------+
sqlite> </code></pre>Also we need to check if the `caller_name` is `Test Caller`
It could be `OR` here
<pre><code class="language-sql">SELECT * FROM hotline_messages 
WHERE 
     tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert')
     OR caller_name = 'Test Caller';</code></pre>That's it let's count the number of rows we will be deleting soon.
<pre><code>SELECT count(*) FROM hotline_messages
WHERE
     tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert') 
     OR caller_name = 'Test Caller';</code></pre>So, we have `89` rows to deleted after considering the spamy tags and test callers out.
<pre><code>sqlite> SELECT count(*) FROM hotline_messages WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert') OR caller_name = 'Test Caller';
+----------+
| count(*) |
+----------+
| 89       |
+----------+
sqlite> 
</code></pre>
Let's get the spam outta here!
<pre><code class="language-sql">DELETE FROM hotline_messages
WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert') OR caller_name = 'Test Caller';</code></pre>Phew! Done 89 spammy records  removed! Santa might be relieved.
<pre><code>sqlite> DELETE FROM hotline_messages
WHERE tag IN ('penguin prank', 'time-loop advisory', 'possible dragon', 'nonsense alert') OR caller_name = 'Test Caller';

sqlite> SELECT changes();
+-----------+
| changes() |
+-----------+
| 89        |
+-----------+</code></pre>
The changes are done, now we simply have to select and count the things which are approved and in review.
### Count'em down<ul><li>After doing those 2 things, we have to count the number of records with `status` as `approved` and the number of records that are still not `approved` (they are in review or the status is `NULL` )
</li></ul>So, we need to get the count of
<ol><li>Records with `status` as `approved`
</li><li>Records with `status`  as `NULL` 
</li></ol>
#### Separate QueriesThis looks straight forward, you can write two separate queries for doing the things.
<pre><code class="language-sql">SELECT COUNT(*) as approved_count FROM hotline_messages WHERE status = 'approved';
SELECT COUNT(*) as in_review_count FROM hotline_messages WHERE status IS NULL;</code></pre>
<pre><code>sqlite> SELECT COUNT(*) as approved_count FROM hotline_messages WHERE status = 'approved';
SELECT COUNT(*) as in_review_count FROM hotline_messages WHERE status IS NULL;
+----------------+
| approved_count |
+----------------+
| 477            |
+----------------+
+-----------------+
| in_review_count |
+-----------------+
| 501             |
+-----------------+
sqlite> 
</code></pre>
#### Group by StatusBut however, can we do it in 1 query?
Think!
There are just 2 types of status right?
Let's check
```
`SELECT DISTINCT status FROM hotline_messages;`
```Hmm! 2? `NULL` and `approved` !
<pre><code class="language-sql">sqlite> SELECT DISTINCT status FROM hotline_messages;
+----------+
|  status  |
+----------+
|          |
| approved |
+----------+</code></pre>So we can simply do the same thing, but just group by `status` right? Like so:
<pre><code class="language-sql">SELECT status, COUNT(*) as count
FROM hotline_messages
GROUP BY status;</code></pre>And this should give us back the two rows with the count of `NULL` and `approved` .
<pre><code>sqlite> SELECT status, COUNT(*) as count
FROM hotline_messages
GROUP BY status; 
+----------+-------+
|  status  | count |
+----------+-------+
|          | 501   |
| approved | 477   |
+----------+-------+</code></pre>Is there a better way?
This looks a little wired! Like status is empty (NULL) and it kind of makes a little wired view for people to look at, can we do something different?
#### Cases when then elseThis is simple use case for a `CASE WHEN ... THEN ... ELSE ... END ` 
For each when we can check certain conditions and do certain things or do other thing.
In this case, if the status is `approved`, we can increment the count to 1 or we can keep 0, similarly the other when can be used for grouping the count of status being `NULL` .
<pre><code class="language-sql">SELECT 
    COUNT(CASE WHEN status = 'approved' THEN 1 END) AS approved_count,
    COUNT(CASE WHEN status IS NULL THEN 1 END) AS in_review_count
FROM 
    hotline_messages;</code></pre>What this will do is, for each row, we will count up the number of either `approved_count`or `in_review_count` depending on the value of the `status` cell. If that is `approved` we increment the `approved`  count else if that is `NULL` we increment the `in_review_count`.
Slick!
<pre><code class="language-sql">sqlite> SELECT 
    COUNT(CASE WHEN status = 'approved' THEN 1 END) AS approved_count,
    COUNT(CASE WHEN status IS NULL THEN 1 END) AS in_review_count
FROM 
    hotline_messages;
+----------------+-----------------+
| approved_count | in_review_count |
+----------------+-----------------+
| 477            | 501             |
+----------------+-----------------+
sqlite> </code></pre>
That's it from day 3 hopefully Santa is happy, and in sight of getting madder as the elves get dumber.