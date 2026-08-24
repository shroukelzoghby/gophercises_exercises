# Exercise #1: Quiz Game

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/quiz)

## Exercise Details

### Part 1

Create a program that will read in a quiz provided via a CSV file and will then give the quiz to a user, keeping track of how many questions they get right and how many they get incorrect. Regardless of whether the answer is correct or wrong, the next question should be asked immediately afterwards.

The CSV file should default to `problems.csv`, but the user should be able to customize the filename via a flag.

The CSV file should use the following format, where the first column is a question and the second column in the same row is the answer:

```csv
5+5,10
7+3,10
1+1,2
8+3,11
1+2,3
8+6,14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7
```

Quizzes can be assumed to be relatively short, with fewer than 100 questions, and to have single-word or number answers.

At the end of the quiz, the program should output the total number of questions answered correctly and how many questions there were in total. Questions given invalid answers are considered incorrect.

> CSV files may have questions with commas in them. For example, `"what 2+2, sir?",4` is a valid row in a CSV file. Use Go's [`encoding/csv`](https://pkg.go.dev/encoding/csv) package instead of writing a custom CSV parser.

### Part 2

Adapt the program from Part 1 to add a timer. The default time limit should be 30 seconds, but it should also be customizable via a flag.

The quiz should stop as soon as the time limit has been exceeded. It should not wait for the user to answer one final question and should stop even if it is currently waiting for an answer from the user.

Users should be asked to press Enter, or another key, before the timer starts. Questions should then be printed one at a time until the user provides an answer. Regardless of whether the answer is correct or wrong, the next question should be asked.

At the end of the quiz, the program should output the total number of questions answered correctly and how many questions there were in total. Questions with invalid answers or no answer are considered incorrect.

## Bonus

1. Add string trimming and cleanup so that correct answers with extra whitespace, capitalization, and similar differences are not considered incorrect. See Go's [`strings`](https://pkg.go.dev/strings) package.
2. Add an option, using a new flag, to shuffle the quiz order each time it is run.
