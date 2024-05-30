package golang

func (this_ *Generator) GenConf() (err error) {
	dir := this_.Dir + "conf/"
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "application.yml"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	builder.AppendCode(this_.GetApp().Text)

	builder.NewLine()

	builder.AppendTabLine("log:")
	builder.AppendTabLine("  console: false # 输出到控制台")
	builder.AppendTabLine("  filename: ./logs/app.log")
	builder.AppendTabLine("  maxSize: 100   # 文件大小单位M")
	builder.AppendTabLine("  maxAge: 7      # 保留多少天")
	builder.AppendTabLine("  maxBackups: 10 # 最多几个文件")
	builder.AppendTabLine("  level: debug   # 级别，debug，info，warn，error")

	builder.NewLine()
	return
}
