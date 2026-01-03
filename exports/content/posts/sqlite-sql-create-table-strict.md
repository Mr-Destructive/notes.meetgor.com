---
title: "SQLite SQL: Create Table STRICT "
date: 2025-12-12
slug: sqlite-sql-create-table-strict
draft: false
type: post
description: ""
tags: []
---

We have seen how to create `TABLE` loosely Not adhering to the types. Because if we create a table with column name of type text, and insert a integer, it will happily store it as text. It is very flexible as we saw. So, in such `CREATE TABLE` statement in SQLite, without the `STRICT` constraint, the types don't matter.

If you created a table with column type as `XYZ` it will accept it, because it really doesn't see that. It will see the data coming in from the insert statement and store it whatever it thinks is the best one for that piece of data. Look at the below example:

<pre><code>CREATE TABLE t1 (n xyz);
INSERT INTO t1 values(4);
INSERT INTO t1 values("gg");
SELECT rowid, n, typeof(n) FROM t1;</code></pre>
<pre><code>rowid	n	typeof(n)
1	    4	integer
2	   gg	text
</code></pre>See? The column type, it doesn't matter.
Unless it's strict or any constraints, or generated conditions have been added.
## The STRICT table optionLet's quote from the documentation what it means
<blockquote><ul><li>Every column definition must specify a datatype for that column.
</li><li>The freedom to specify a column without a datatype is removed.
</li><li>The datatype must be one of the following:
<ul><li>INT
</li><li>INTEGER
</li><li>REAL
</li><li>TEXT
</li><li>BLOB
</li><li>ANY
</li></ul></li><li>The [PRAGMA integrity_check](https://sqlite.org/pragma.html#pragma_integrity_check) and [PRAGMA quick_check](https://sqlite.org/pragma.html#pragma_quick_check) commands check the type of the content of all columns in STRICT tables and show errors if anything is amiss.
</li></ul></blockquote>There are other nuances of the STRICT table options and the kind of constraint that you put on the columns, but that requires studying very specific examples. We'll check those nuances later.
For now though, we need to understand how to create a strictly typed table, and what the strict option adds to the table.
<pre><code>CREATE TABLE users(
    name TEXT,
    age  INT,
    credits REAL,
    profile_pic BLOB
) STRICT;</code></pre>So, we have all the actual possible types we can use in a table column when defining a table. If you don't provide an column type, or provide any other type than `TEXT`, `INT` or `INTEGER`, `REAL`, `BLOB`, or `ANY` (don't put any, you lose the purpose of strict) it won't compile and execute the table creation. You need to provide a valid type among the 5 types.
However if you try to create a strict table with wrong column type or no column type.
<pre><code>CREATE TABLE t1 (t) STRICT;
-- Error: missing datatype for t1.t

CREATE TABLE t1 (t something) STRICT;
-- Error: unknown datatype for t1.t: "something"</code></pre>
Without STRICT it works as usual:
<pre><code>CREATE TABLE t1 (t  something);
INSERT INTO t1 values(123), ('abc'), (X''), (123.45);
SELECT t, typeof(t) FROM t1;</code></pre>
<pre><code>t	typeof(t)
123	integer
abc	text
    blob
123.45	real</code></pre>
Now back to the original example:
Insert a couple of rows:
<pre><code class="language-sql">-- All are NULL Values
INSERT INTO users DEFAULT VALUES;

INSERT INTO users (name, age, credits, profile_pic)
VALUES (
    'Alice',
    30,
    100.0,
    X'89504E470D0A1A0A'
);</code></pre>
This will insert two rows, the first one, all the columns will be `NULL` . If you look at the type of these statement. Those will be as per the table schema, consistent for all rows.
<pre><code>name	typeof(name)	age	typeof(age)	credits	typeof(credits)	profile_pic	typeof(profile_pic)
null		null		null		null
Alice	text	30	integer	100	real	137,80,78,71,13,10,26,10	blob
</code></pre>This has rightly added `NULL` type for the null values but when the data is in the row, it forces that type stated in the schema of the table.

Now, if we try to mess up the column data, it won't work
<pre><code class="language-sql">INSERT INTO users (name, age, credits, profile_pic)
VALUES (34, '4', 8, 123);
-- Error: cannot store INT value in BLOB column users.profile_pic

INSERT INTO users (name, age, credits, profile_pic)
VALUES (34, '4', 8, '');
-- Error: cannot store TEXT value in BLOB column users.profile_pic

INSERT INTO users (name, age, credits, profile_pic)
VALUES (34, 'abc', 8, X'');
-- Error: cannot store TEXT value in INT column users.age
</code></pre>
This will work, as type affinity and the conversion is possible within the column types here.
<pre><code>INSERT INTO users (name, age, credits, profile_pic)
VALUES (CAST(34 AS INT), '3', 8, X'');</code></pre>But if some data is not able to convert into that strict type, it will fail the constraint of strict column type.
For instance
<ul><li> `123` or `""` is not force convertible to BLOB which is binary large object. We need to parse it with X'' strings for some raw data to make it a BLOB like object in SQLite.
</li><li>`abc`  is not convertible/casteble to INTEGER or REAL Value.
</li></ul>So, the strict type is actually strict as we see the pattern repeating in SQLite.
<blockquote>It is flexible till you allow it to be, you can at anytime change the lever and make it strict
</blockquote>This is true for column-row level type checking with the STRICT table option while creating table. 
##