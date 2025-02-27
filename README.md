# Programming Languages, Spring 2025

- Instructor: Prof. Y C Cheng
- Class meeting time: Mon 2-3-4
- Class meeting place:
  - Second Academic Building-207 (in person)
  - [Microsoft Teams](https://teams.microsoft.com/l/team/19%3ACxxkswKurkJFiUujdCtAfg7t9IO6ZwPGO2EiY4AyicA1%40thread.tacv2/conversations?groupId=763aba6a-09fe-4b95-9837-b7cf8d08d47d&tenantId=dfb5e216-2b8a-4b32-b1cb-e786a1095218)
- Office hours: Mon 5-6-7, Tue 2-3-4
- Course Repo: Course repository: http://140.124.181.100/yccheng/pl2025s
- TA: 
  - Benny Wang \<benny870704@gmail.com\>, Regina Yen\<ycycchre@gmail.com\> (Hong-Yue technology Research Building, Room 1324)
  - Office hours: 10-12 am, Wed
- TA's homework repository: [pl2025_TA](http://140.124.181.100/course/pl2025s_ta)
- Please hand in the homework on [Microsoft Teams](https://teams.microsoft.com/l/team/19%3ACxxkswKurkJFiUujdCtAfg7t9IO6ZwPGO2EiY4AyicA1%40thread.tacv2/conversations?groupId=763aba6a-09fe-4b95-9837-b7cf8d08d47d&tenantId=dfb5e216-2b8a-4b32-b1cb-e786a1095218).

- HW1 is announced, deadline: 3/3(Mon.) 23:59.

## Synopsis

We will take a practice-first approach to learn the most frequently seen language features topic by topic.

Each topic consists of three parts given a problem:
1. A program we will write in **a style satisfying the constraints** in solving the problem. We will discuss the consequences of the style
2. The principles in language design that implements the constraints
3. The same program in Go written either in class or as an exercise to illustrate how the constraints are satisfied.

Some problems we will solve:
- [Term frequencies](https://github.com/crista/exercises-in-programming-style)

We will repeat the three-parts for a number of styles under the following programming paradigms:
- structured programming
    - Cookbook
    - Monolithic
    - Pipeline
- functional programming
    - Infinite Mirror
    - Kick Forward
    - The One
- object-oriented programming
    - Things
    - Letter Box
    - Abstract Things
    - Hollywood
    - Bulletin Board
- concurrency
    - Map Reduce

## Reference Books
- Donovan, Alan AA, and Brian W. Kernighan. The Go programming language. Addison-Wesley Professional, 2015.
- Cristina Videira Lopes, Exercises in Programming Style, First Edition, Chapman and Hall/CRC, 2014.
- Michael L. Scott, Programming Language Pragmatics, Fourth Edition, Morgan Kaufmann (Elsevier), 2016.

**Week 1**

- Program: Cookbook
    - run the program
    - unit-tests it
- Principle: Names, Bindings, and scope
- Program in Go and constraints resolution
    - TDD the Go program
    - Reexamine the Principle

**Week 2**
- pipeline style and mathematical function
- why single input and output?
- Currying
- idempotence
- runtime environment - stack needed (Pragmatics).

== HW1, deadline: 3/3(Mon.) 23:59

1. (30%) Do exercise 5.2 of **Exercises in Programming Style** (see below.) Hint: apply currying to remove_stop_words. Write unit tests for this and change the main program to use this function. 

Example program: https://github.com/crista/exercises-in-programming-style/blob/master/06-pipeline/tf-06.py

In the example program, the name of the file containing the list of stop words is hard-coded (line #36). Modify the program so that the name of the stop words file is given as the second argument in the command line. You must observe the following additional stylistic constraints: (1) no function can have more than 1 argument, and (2) the only function that can be changed is remove_stop_words; it’s ok to change the call chain in the main block in order to reflect the changes in remove_stop_words.

2. (30%) Show with unit testing for the functions in the Cookbook style that are **not idempotent** and describe the reason in the comment.

3. (40%) Rewrite the Pipeline style in Go, including unit tests.