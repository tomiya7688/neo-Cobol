# 開発コンテキストの入口

Neo COBOL は Go で実装するコンパイラです。開発時は、選択した Issue に必要な情報を絞って読みます。
この入口は参照先と探索手順を示します。仕様や実装の内容は原典で確認してください。

## 正式な参照先

- 言語の設計・仕様: [SPECIFICATION.md](docs/language/SPECIFICATION.md) と、その項目から参照する仕様書。
- コンパイラの構成・実装範囲: [ARCHITECTURE.md](docs/compiler/ARCHITECTURE.md)。
- 実装の現状: 対象のソースとテスト。仕様の草案を実装済みの保証として扱わない。
- 今回の目的と受け入れ条件: 選択した Issue と関連する PR。要約は原典への案内として使う。
- CLI の基本操作: [README.md](README.md)。

## 作業手順と探索の停止条件

1. Issue から Goal（目的）、Required（必要な変更・制約）、Acceptance（受け入れ条件）を短く整理する。
2. checkout、remote、変更中のファイル、対象ブランチ・PR の現在状態を確認する。リモートは概要と変更ファイルから確認し、必要な差分だけを読む。
3. 下表から該当する行を選ぶ。複数の仕様が並ぶ行でも、今回の変更に関係する仕様だけを開く。
4. `rg --files` で候補を絞り、`rg -n '<symbol>' <対象ディレクトリ>` で位置を探して、必要な行と直接の依存先を読む。
5. 目的、制約、受け入れ条件、変更箇所、検証方法が分かったら探索を止める。未解決の疑問があるときだけ、その疑問に必要な参照を追加する。
6. 変更と検証結果を Issue に対応させて報告する。実行していない検証は「未検証」と明記する。

ファイル全体や長いログが必要なのは、部分的な参照では判断できない場合です。
過去の要約を使う場合は対象のブランチ・リビジョンを確認し、変更されている原典を読み直します。

## 作業別の参照先

| 作業 | 必要に応じて読む仕様 | 最初に検索する実装 |
| --- | --- | --- |
| 字句・ソース形式 | docs/language/LEXICAL.md、SOURCE_FORMAT.md | internal/lexer |
| 文法・文 | docs/language/GRAMMAR.md、STATEMENTS.md | internal/parser、internal/ast |
| 型・PIC・変数・データ・変換 | docs/language/TYPES.md、PIC.md、VARIABLES.md、DATA_MODEL.md、CONVERSIONS.md | internal/sema。構文変更時は internal/parser |
| NIR・C 出力・処理パイプライン | docs/compiler/ARCHITECTURE.md | internal/nir、internal/backend/c、internal/compiler |
| CLI | README.md | cmd/neoc、internal/compiler |
| クラス・関数・引数の仕様 | docs/language/OOP.md、FUNCTIONS.md、PARAMETERS.md | 仕様を先に確認し、対象の実装を検索する |

依存範囲は対象コードの import・呼び出し元・共有データ構造から広げます。
表にあるディレクトリ全体を一括でコンテキストへ投入する必要はありません。

## 検証の選び方

- 文書のみの変更: 差分、参照先、既存仕様との整合を確認する。
- Go の変更: `rg --files <対象ディレクトリ> -g '*_test.go'` で関連テストを探し、該当パッケージを検証する。
- 共有 AST・NIR・処理パイプラインなど影響範囲が広い変更: `go test ./...` も検討する。
- C 出力・実行の変更: C コンパイラの利用可否を確認し、関連サンプルの生成・ビルド・実行を検証する。
- 結果はコマンド、成否、必要なエラー抜粋を短く記録する。静的確認を実行時の保証として報告しない。

## 通常は読み込まないもの

ビルド出力、キャッシュ、生成物、長い実行ログ、全 Issue の本文、全 PR の差分、無関係な履歴。
受け入れ条件の判断に必要なものは、該当箇所に絞って読みます。

## 品質基準

- レビュー時は [レビュー基準.md](レビュー基準.md) を使います。
- 実装・文書作成時は [コーディング規約.md](コーディング規約.md) を使います。
- 対象の仕様・実装に関係する項目だけを参照し、無関係な箇所を一括で読み込まないでください。

## 導入範囲

[ai-context-reducer](https://github.com/tomiya7688/ai-context-reducer) の最小コアを基に、
小さな入口、作業別の参照、検索してから読む手順、探索の停止条件を採用しています。
専用ツールや大きな依存グラフの導入は、具体的な必要性が出た段階で判断します。
削減量は未測定です。
