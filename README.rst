GoSQL(SQLite)
==================


.. contents::

简化SQL查询，查询返回结果为 `Map` 结构。


.. image:: https://img.vim-cn.com/c7/c310fc26e8a6be6a1ad55b8692acea27363a85.jpg

使用
-------

程序改为了 命令行调用模式。

.. code:: bash
    # 编译
    go build main.go
    
    # 初始化数据库
    ./main init
    # 插入数据
    ./main query "INSERT INTO auth (uid, name, nickname) VALUES(1, '名字', '昵称') "
    # 查询数据
    ./main query "SELECT * FROM auth"
    # 修改数据
    ./main query "UPDATE auth SET name = '名字很长' WHERE uid=1 "


参考资料：

*   https://golang.org/doc/code.html
*   https://segmentfault.com/a/1190000002623278#articleHeader21
*   https://golang.org/pkg/database/sql/#DB
*   https://github.com/go-sql-driver/mysql/wiki/Examples
*   https://golang.org/pkg/encoding/json/
*   https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
