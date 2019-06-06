# quchong


家里有很多照片文档 ，怕弄丢了。复制来复制去的。备份了很多地方，都乱了。想找个去重的软件，windows下的看到个还是个收费的
感觉不难自己写一个。

目录路径问题，暂时支持windos，linux也可以，需要修改代码，未做测试

    usages：COMMAND -f DIR
            -f  指定要做去重的目录

    指定目录的父目录必须和目标目录在同一磁盘,widows,好像不存在这个问题。linux主要注意。
    重复文件会再父目录下创建相应层级的目录，并移动文件。因为是同一磁盘移动速度比较快。
    
    如果文件较多，切比较大的话计算md5值比较耗费时间
    测试计算单个8.32G的文件大概24s
    
    例如  C:\Users\yyw\go\src\github.com\yangyouwei\quchong\test
    会自动创建备份目录 [test]-samefile
    