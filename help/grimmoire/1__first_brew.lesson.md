# Lesson 1: Your First Brew

**Grimm Club · ~30 minutes · no git experience needed**

By the end of this lesson you'll have made a real project, written a tiny bit of
code, and saved it to history with a commit message you can actually be proud of.

You'll need: Alchemist installed (`alchemist --help` should work) and `git`.

---

## The big idea

Most beginners write code first and figure out what to say about it afterwards —
which is why so many commits just say "stuff" or "fix." Alchemist flips that
around. **You name the task first.** That one habit makes everything downstream
easier.

So our motto today: *name it, make it, brew it.*

---

## Step 1 — Build your workshop (3 min)

Pick an empty folder to work in, then:

```
alchemist recipe python
```

Choose a project name when asked (try `greeter`). Then move into it:

```
cd greeter
```

You now have a tiny Python project that's already tracking history.

## Step 2 — Name the task (2 min)

Before touching any code, decide what you're going to do and say so:

```
alchemist start
```

- **Task name:** `Add a friendly greeting`
- **Description:** `Print a hello message with the user's name`

That name is now reserved as your future commit title. Notice you haven't
written any code yet — that's the point.

## Step 3 — Make it (10 min)

Open `main.py` and change it to ask for a name and greet it:

```python
def main():
    name = input("What's your name? ")
    print(f"Greetings, {name}! Welcome to the Grimm Club.")


if __name__ == "__main__":
    main()
```

Run it to make sure it works:

```
python main.py
```

## Step 4 — Brew it (5 min)

Now turn that work into a commit:

```
alchemist brew
```

Watch what happens:

1. It lists the files you changed and asks to stage them — say yes.
2. It asks, *"In one line, what did you do?"* — answer something like
   `Ask for a name and greet the user`.
3. A commit message opens in your editor, already half-written for you. The
   first line is your task name; the body has your summary and the file list.
   Tidy it up if you like, then save and close.

You just wrote a clear, detailed commit — without staring at a blank prompt.

## Step 5 — Reflect (5 min · together)

Talk it through as a club:

- How was naming the task *first* different from how you'd normally start?
- Look at your commit message. Could a teammate understand it in six months?
- What's the difference between `discard`, `stash`, and `clean`? (Peek at the
  [Grimoire](../GRIMMOIRE.md) if you're unsure — knowing which to reach for is a
  real skill.)

---

## If something goes wrong

- **Changed a file and want it back the way it was?** `alchemist discard`
- **Need to stop mid-task but not lose your work?** `alchemist stash`, and later
  `alchemist stash --resume`
- **"no task in progress" when you brew?** You skipped Step 2 — run
  `alchemist start` first.

---

## Homework

Start a second task with `alchemist start` (try *"Greet by time of day"*), add a
morning/afternoon/evening message, then `brew` it. Bring your commit message to
the next meeting and we'll read a few aloud.

*Next lesson: bottling a release — version tags and pushing to GitHub.*
