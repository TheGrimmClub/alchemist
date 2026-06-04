# The Alchemist's Grimoire

A spellbook for the Grimm Club. Each spell is a real git workflow underneath —
the incantation is just friendlier. The "beneath the spell" notes tell you what
git is actually doing, so you're learning the real craft, not just the costume.

> Golden rule of the order: **name the task before you cast.** Decide what
> you're making, *then* make it.

---

## The Brewing Cycle

These four are the heart of the craft. You'll cast them in order, over and over.

### `recipe` — lay the foundation

```
alchemist recipe python
```

Conjures a fresh project from a template and prepares it to record history.

*Beneath the spell:* creates the project files and runs `git init`.

### `start` — name your intent

```
alchemist start
```

Speaks the name of the task aloud before any code is written. That name becomes
the title of the commit you'll brew later. This is the most important spell in
the book — it's where you decide what you're doing.

*Beneath the spell:* saves your task to `.alchemist/current.json`. No git yet —
this is pure intention.

### `brew` — transmute work into a commit

```
alchemist brew
```

Gathers your changes, asks what you did, and prepares a detailed scroll (the
commit message). The scroll opens in your editor so you can perfect the wording
before it's sealed.

*Beneath the spell:* `git add -A`, then `git commit` with a message whose title
is your task name and whose body lists the files you changed.

### `bottle` — seal and send

```
alchemist bottle
```

Optionally stamps a version on your work, then sends it out into the world.
Clears the current task so you're ready to begin anew.

*Beneath the spell:* an optional `git tag -a`, then `git push --follow-tags`.

---

## Spells for Mishaps

Every alchemist spills something eventually. These three put it right. Notice
they do *different* things — choosing the right one is part of the skill.

### `discard` — undo a botched brew

```
alchemist discard
```

Wipes changes to files git already knows about, returning them to their last
sealed state. Asks for confirmation, because this **cannot be undone**.

*Beneath the spell:* `git restore`. Affects **tracked** files only.

### `stash` — pause the cauldron

```
alchemist stash            # set work aside
alchemist stash --resume   # bring it back
```

Tucks your half-finished work away so you can chase something urgent, then
pour it back when you return. Nothing is lost — it's a pause, not a deletion.

*Beneath the spell:* `git stash push` and `git stash pop`.

### `clean` — sweep the workbench

```
alchemist clean
```

Sweeps away clutter git was never tracking — build artifacts, stray scratch
files. Shows you the list first and asks before deleting.

*Beneath the spell:* `git clean -fd`. Affects **untracked** files only.

---

## Knowing which mishap-spell to cast

| Your situation                                   | Cast      |
|--------------------------------------------------|-----------|
| "I changed a tracked file and want it back"      | `discard` |
| "I need to step away mid-task without committing" | `stash`   |
| "My folder is full of build junk git ignores"    | `clean`   |

---

## The shortest possible quest

```
alchemist recipe python     # build the workshop
alchemist start             # name the task
# ...write your code...
alchemist brew              # commit it well
alchemist bottle            # tag and push
```

Cast in that order and you'll never write a mystery commit message again.
