# mylog

##1、支持往不同的地方输出日志
##2、日志分级别
    * Debug  
    * Trace  
    * Warning  
    * Info  
    * Error  
    * Fatal  
##3、日志要支持开关
##4、完整的日志记录要包含有时间、行号、文件名、日志级别、日志信息
##5、日志文件要切割
    * 按文件大小切割  
    * 每次记录日志之前判断一下当前写的这个文件的文件大小  
##2、按日期切割
    * 在日志结构体中设置一个字段记录上一次切割的小时数
    * 在写日志之前检查一下当前时间的小时数和之前保存的是否一致，不一致要切割