# 文件批量改名工具

1. 遍历某文件夹里所有子目录，把文件全部提到第一层目录，并把文件夹全部删掉
2. 根据正则匹配找到匹配文件名并替换生成新的文件名，进行重命名

---

例：

    batrename path exgin extout

---

参数定义：

- `path` 要进行处理的路径
- `exgin` 匹配规则
- `exgout` 替换规则

---

to do:

文件名或文件夹名有空格，会导致恐慌