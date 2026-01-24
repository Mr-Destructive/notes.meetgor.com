---
title: "CS50 SQL Problem set 1 packages"
date: 2026-01-24
draft: false
---

# CS50 SQL Problem set 1 packages

**Link:** https://cs50.harvard.edu/sql/psets/1/packages/

## Context

- You can’t do a unpacking in a sql query when using nested query for more than one column returned
    - Like example i have a query like this
      ```
      SELECT 
        s.id, 
        s.package_id,
        address_id, 
        a.address AS from_address, 
        (
          SELECT 
            address,
          FROM 
            addresses 
          WHERE 
            id = p.to_address_id
        ) AS to_address 
      FROM 
        scans AS s 
        INNER JOIN packages AS p ON p.id = s.package_id 
        INNER JOIN addresses AS a on a.id = p.from_address_id;

      ```
    - But let’s say for some reason I wanted to also get the to_address_type like a column from the addresses table
    - You might try to over-optimise the queries and try something like this
      ```
      SELECT 
        s.id, 
        s.package_id, 
        action, 
        contents, 
        address_id, 
        a.address AS from_address, 
        (
          SELECT 
            address, 
            type 
          FROM 
            addresses 
          WHERE 
            id = p.to_address_id
        ) AS (to_address, to_address_type) 
      FROM 
        scans AS s 
        INNER JOIN packages AS p ON p.id = s.package_id 
        INNER JOIN addresses AS a on a.id = p.from_address_id 
      WHERE 
        address = '900 Somerville Avenue' 
        AND s.action = 'Drop';

      ```
    - And ERROR, you can’t do that
    - This bit right here
      ```
      (SELECT address, type FROM addresses WHERE id=p.to_address_id)
      AS
      (to_address, to_address_type)
      ```
    - This is not feasible in SQL, you can’t unpack multiple columns from a subquery directly and alias them inline in a single SELECT clause.
    - Well I have to do it this way then, duhh

  ```
  SELECT 
    s.id, 
    s.package_id, 
    action, 
    contents, 
    address_id, 
    a.address AS from_address, 
    (
      SELECT 
        address 
      FROM 
        addresses 
      WHERE 
        id = p.to_address_id
    ) AS to_address, 
    (
      SELECT 
        type 
      FROM 
        addresses 
      WHERE 
        id = p.to_address_id
    ) AS to_address_type 
  FROM 
    scans AS s 
    INNER JOIN packages AS p ON p.id = s.package_id 
    INNER JOIN addresses AS a on a.id = p.from_address_id 

  ```
    - What a long query!
    - By the way, this is one of the questions in [CS50 SQL Problem set 1 packages](https://cs50.harvard.edu/sql/psets/1/packages/) section.
    - Any better way to do this? drop them in the comments or hit me up on my socials, will be completing more challenges this weekend.

## Tech News

**Source:** techstructive-weekly-53
