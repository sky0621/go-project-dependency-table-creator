# glideでパッケージ管理しているGoプロジェクト間の依存表({{.Datetime}} 時点)

#### ※ツール（ https://github.com/sky0621/go-project-dependency-table-creator ）による自動生成

{{range .Headers}}| [{{.Display}}]({{.URL}}) {{end}}|
{{range .Headers}}| :--- {{end}} |
{{range .Bodies}}{{range .}}| {{.}} {{end}}|
{{end}}
