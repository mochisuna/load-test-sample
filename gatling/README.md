# gatling用シナリオ
sbtを利用したテストプロジェクトです。

# Scenarios
- computerdatabase.LoadTestSimulation: 標準シナリオテスト
- computerdatabase.LoadTestSimulationRampUp10Users: 5秒かけてユーザー数を増加させる加負荷テスト
- computerdatabase.LoadTestSimulationCompoundTest: 複合テスト（並列実行）
- computerdatabase.LoadTestSimulationWithPreprocessing: 事前処理を入れた標準シナリオテスト

# USAGE
1. gatlingディレクトリに移動
1. 検証対象のURIを `TARGET_URI` として環境変数に定義
1. sbtコマンドでターゲットを指定して実行

```bash
$ cd <YOUR_PATH>/load-test-sample/gatling
$ set -x TARGET_URI "http://localhost:<PORT>/v1"
$ sbt "gatling:testOnly computerdatabase.LoadTestSimulation"
```
