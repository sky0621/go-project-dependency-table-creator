# プロジェクト間の依存表({{.Datetime}} 時点)

#### ※ツール（ https://github.com/sky0621/go-project-dependency-table-creator ）による自動生成

#### 【前提】

##### ・Goのプロジェクトであること

##### ・glideでパッケージ管理している（＝glide.yamlを配置している）こと

{{range .Headers}}| {{.}} {{end}}|
{{range .Headers}}| :--- {{end}} |
{{range .Bodies}}{{range .}}| {{.}} {{end}}|
{{end}}
