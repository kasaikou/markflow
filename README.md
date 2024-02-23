<div align="center">

# docstak: Task Runner as a Document (TRaaD) tool<br>🐶🥞

[[English](./README.md)] [[日本語](./README.ja.md)]
</div>

## Getting Started

Write markdown document as configuration file.

````md
<!-- ./docstak.md -->

# Getting started

## hello_world

Echo "Hello World, docstak!"

```sh
echo "Hello World, docstak!"
```
````

And use `docstak` command.

```terminal
$ docstak hello_world
DOCSTAK INFO                task start (task: "hello_world")
STDOUT  hello_world         Hello World, docstak!
DOCSTAK INFO                task ended (task: "hello_world", exitCode: 0)
```

## Concepts

`docstak` is a task runner tool that reads dependencies between scripts and tasks from `.md` files and executes necessary scripts.

Traditionally, executing workflows through scripts involved executing tasks using script files like `.sh`, task runners like `make`, `task`, and other task runner tools.
However, as the number of build tools increases, or when dealing with large repositories such as monorepos, managing workflows becomes difficult due to the increasing number of workflows.

Typically, documentation is used within a team to share these workflows.
However, doesn't synchronizing changes between actual workflows and documentation become tiresome?
It's natural to feel this way because the aforementioned methods are solely focused on executing scripts and do not provide documentation functionality.

`docstak` reads Markdown, a means of documentation, and executes it.
While the aforementioned task runner tools primarily focus on executing scripts, `docstak` allows you to construct workflows along with documentation.
Please take a look at [`docstak.md`](./docstak.md).
You'll see that workflows are constructed using standard Markdown syntax, which can be rendered into HTML seamlessly using existing Markdown renderers.

## Contribute

We welcome contributions such as issues and pull requests following the [Code of Conduct](./CODE_OF_CONDUCT.md).

### Language

We understand that forcing the use of English, especially during the early stages of development, can slow down the process unnecessarily.
Therefore, we will follow the following guidelines for the time being.
Thank you for your understanding.

| | Language | Notes |
| :-: | :-: | :-- |
| Issue & Pull Request | English & Japanese | Please bear with us as English is not our strong suit. |
| Commit Message | English |  |
| Code Comment | English |  |
