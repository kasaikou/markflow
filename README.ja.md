<div align="center">

# docstak: Task Runner as a Document (TRaaD) tool<br>🐶🥞

[[English](./README.md)] [[日本語](./README.ja.md)]

[![codecov](https://codecov.io/gh/kasaikou/docstak/graph/badge.svg?token=QZTCJ1A852)](https://codecov.io/gh/kasaikou/docstak)
</div>

## Getting Started

Markdown ファイルを設定ファイルとして記述します。

````md
<!-- ./docstak.md -->

# Getting started

## hello_world

Echo "Hello World, docstak!"

```sh
echo "Hello World, docstak!"
```
````

そして `docstak` コマンドを実行します.

```terminal
$ docstak hello_world
DOCSTAK INFO                task start (task: "hello_world")
STDOUT  hello_world         Hello World, docstak!
DOCSTAK INFO                task ended (task: "hello_world", exitCode: 0)
```

## Concepts

`docstak` は .md からスクリプトやタスク間の依存関係を読み取り、必要なスクリプトを実行していくタスクランナーツールです。

従来からスクリプトによるワークフローの実行には、 `.sh` などのスクリプトファイルによるタスクの実行、タスクランナーとしての `make` 、 `task` などのタスクランナーツールなどの手段があります。
少量であればいいものの、多くのビルドツールを使い始めたり、モノリポなどの大規模なリポジトリになるとワークフローの数も増えてきて管理することが困難になっていきます。

ところで、一般的にはこれらワークフローをチーム内で共有するためにドキュメンテーションを行っていくものであると思うのですが、実際のワークフローとドキュメンテーションを同期しながら変更を加えていく事は疲れませんか？
当然ではありますが、前述した手段はいずれもスクリプトを実行することにのみ特化しており、ドキュメンテーションとしての機能を提供していません。

`docstak` はドキュメンテーションの手段である Markdown を読み取って実行します。
前述したタスクランナーツールが単にスクリプトの実行を主眼に置いているのに対して、 `docstak` はドキュメントとセットでワークフローを構成することができるのです。
ぜひ [`docstak.md`](./docstak.md) をご覧ください。
既存の Markdown レンダラーを使うことで問題なく HTML にレンダリング可能な普通の Markdown 記法でワークフローを構成していることがわかります。

## Contribute

[Code of Conduct](./CODE_OF_CONDUCT.md) に則った Issue や Pull Request などのコントリビューションを歓迎します。

### Language

開発初期段階で苦手な英語を無理に使って開発速度を落とすいわれはないので、しばらくは以下の運用で行っていきます。
よろしくお願いします。

| | 言語 | 補足
| :-: | :-: | :--
| Issue & Pull Request | English & Japanese | 英語が苦手なので勘弁してください
| Commit Message | English |
| Code comment | English |
